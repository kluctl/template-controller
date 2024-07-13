/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"github.com/kluctl/go-jinja2"
	templatesv1alpha1 "github.com/kluctl/template-controller/api/v1alpha1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// TextTemplateReconciler reconciles a TextTemplate object
type TextTemplateReconciler struct {
	BaseTemplateReconciler
}

//+kubebuilder:rbac:groups=templates.kluctl.io,resources=texttemplates,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=templates.kluctl.io,resources=texttemplates/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=templates.kluctl.io,resources=texttemplates/finalizers,verbs=update

// Reconcile a resource
func (r *TextTemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	logger := log.FromContext(ctx)

	logger.V(1).Info("Starting reconcile")
	defer logger.V(1).Info("Finished reconcile", "err", err)

	var tt templatesv1alpha1.TextTemplate
	err = r.Get(ctx, req.NamespacedName, &tt)
	if err != nil {
		err = client.IgnoreNotFound(err)
		if err != nil {
			logger.Error(err, "Get failed")
		}
		return
	}

	// Add our finalizer if it does not exist
	if !controllerutil.ContainsFinalizer(&tt, templatesv1alpha1.TextTemplateFinalizer) {
		patch := client.MergeFrom(tt.DeepCopy())
		controllerutil.AddFinalizer(&tt, templatesv1alpha1.TextTemplateFinalizer)
		if err := r.Patch(ctx, &tt, patch, client.FieldOwner(r.FieldManager)); err != nil {
			logger.Error(err, "unable to register finalizer")
			return ctrl.Result{}, err
		}
	}

	// Examine if the object is under deletion
	if !tt.GetDeletionTimestamp().IsZero() {
		return r.finalize(ctx, &tt)
	}

	// Return early if the object is suspended.
	if tt.Spec.Suspend {
		logger.Info("Reconciliation is suspended for this object")
		return ctrl.Result{}, nil
	}

	patch := client.MergeFrom(tt.DeepCopy())
	err = r.doReconcile(ctx, &tt)
	if err != nil {
		c := metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionFalse,
			ObservedGeneration: tt.GetGeneration(),
			Reason:             "Error",
			Message:            err.Error(),
		}
		apimeta.SetStatusCondition(&tt.Status.Conditions, c)
	} else {
		c := metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionTrue,
			ObservedGeneration: tt.GetGeneration(),
			Reason:             "Success",
			Message:            "Success",
		}
		apimeta.SetStatusCondition(&tt.Status.Conditions, c)
	}
	err = r.Status().Patch(ctx, &tt, patch, SubResourceFieldOwner(r.FieldManager))
	return
}

func (r *TextTemplateReconciler) doReconcile(ctx context.Context, tt *templatesv1alpha1.TextTemplate) error {
	j2, err := NewJinja2()
	if err != nil {
		return err
	}
	defer j2.Close()

	objClient, err := r.getClientForObjects(tt.Spec.ServiceAccountName, tt.GetNamespace())
	if err != nil {
		return err
	}

	wt := r.watchesUtil.getWatchesForTemplate(client.ObjectKeyFromObject(tt))
	wt.setClient(objClient, tt.Spec.ServiceAccountName)
	newObjects := map[templatesv1alpha1.ObjectRef]struct{}{}
	if tt.Spec.TemplateRef != nil && tt.Spec.TemplateRef.ConfigMap != nil {
		ns := tt.Spec.TemplateRef.ConfigMap.Namespace
		if ns == "" {
			ns = tt.Namespace
		}
		objRef := templatesv1alpha1.ObjectRef{
			APIVersion: "v1",
			Kind:       "ConfigMap",
			Namespace:  ns,
			Name:       tt.Spec.TemplateRef.ConfigMap.Name,
		}
		err = wt.addWatchForObject(ctx, objRef)
		if err != nil {
			return err
		}
		newObjects[objRef] = struct{}{}
	}
	for _, me := range tt.Spec.Inputs {
		if me.Object != nil {
			err = wt.addWatchForObject(ctx, me.Object.Ref)
			if err != nil {
				return err
			}
			newObjects[me.Object.Ref] = struct{}{}
		}
	}
	wt.removeDeletedWatches(newObjects)

	var templateStr string
	if tt.Spec.Template != nil {
		templateStr = *tt.Spec.Template
	} else if tt.Spec.TemplateRef != nil {
		ref := r.buildTemplateRef(tt)
		if ref == nil {
			return fmt.Errorf("no template ref specified")
		}
		if tt.Spec.TemplateRef.ConfigMap != nil {
			jp := fmt.Sprintf("data[\"%s\"]", tt.Spec.TemplateRef.ConfigMap.Key)
			elems, err := r.buildObjectInput(ctx, r.Client, tt.GetNamespace(), *ref, &jp, false, true)
			if err != nil {
				return fmt.Errorf("failed to template from %s: %w", ref, err)
			}
			x, ok := elems[0].(string)
			if !ok {
				return fmt.Errorf("unexpected error. Element is not a string")
			}
			templateStr = x
		} else {
			return fmt.Errorf("no template ref specified")
		}
	} else {
		return fmt.Errorf("no template specified")
	}

	statusBackup := tt.Status
	tt.Status = templatesv1alpha1.TextTemplateStatus{}

	vars, err := r.buildBaseVars(tt, "textTemplate")
	tt.Status = statusBackup
	if err != nil {
		return err
	}

	for _, input := range tt.Spec.Inputs {
		if input.Object != nil {
			elems, err := r.buildObjectInput(ctx, objClient, tt.GetNamespace(), input.Object.Ref, input.Object.JsonPath, false, true)
			if err != nil {
				return fmt.Errorf("failed to get object %s: %w", input.Object.Ref.String(), err)
			}
			MergeMap(vars, map[string]interface{}{
				"inputs": map[string]any{
					input.Name: elems[0],
				},
			})
		} else {
			return fmt.Errorf("missing input")
		}
	}

	rendered, err := j2.RenderString(templateStr, jinja2.WithGlobals(vars))
	if err != nil {
		return err
	}

	tt.Status.Result = rendered

	return nil
}

func (r *TextTemplateReconciler) finalize(ctx context.Context, obj *templatesv1alpha1.TextTemplate) (ctrl.Result, error) {
	r.watchesUtil.removeWatchesForTemplate(client.ObjectKeyFromObject(obj))

	// Remove our finalizer from the list and update it
	controllerutil.RemoveFinalizer(obj, templatesv1alpha1.TextTemplateFinalizer)
	if err := r.Update(ctx, obj, client.FieldOwner(r.FieldManager)); err != nil {
		return ctrl.Result{}, err
	}

	// Stop reconciliation as the object is being deleted
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TextTemplateReconciler) SetupWithManager(mgr ctrl.Manager, concurrent int) error {
	r.Manager = mgr

	c, err := ctrl.NewControllerManagedBy(mgr).
		For(&templatesv1alpha1.TextTemplate{}, builder.WithPredicates(
			predicate.Or(predicate.GenerationChangedPredicate{}),
		)).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: concurrent,
		}).
		Build(r)
	if err != nil {
		return err
	}
	r.controller = c

	err = r.watchesUtil.init(r.RawWatchContext, r.controller)
	if err != nil {
		return err
	}

	return nil
}

func (r *TextTemplateReconciler) buildTemplateRef(tt *templatesv1alpha1.TextTemplate) *templatesv1alpha1.ObjectRef {
	if tt.Spec.TemplateRef == nil {
		return nil
	}
	if tt.Spec.TemplateRef.ConfigMap != nil {
		ns := tt.Spec.TemplateRef.ConfigMap.Namespace
		if ns == "" {
			ns = tt.GetNamespace()
		}
		return &templatesv1alpha1.ObjectRef{
			APIVersion: "v1",
			Kind:       "ConfigMap",
			Namespace:  ns,
			Name:       tt.Spec.TemplateRef.ConfigMap.Name,
		}
	} else {
		return nil
	}
}
