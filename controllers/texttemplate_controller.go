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
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	templatesv1alpha1 "github.com/kluctl/template-controller/api/v1alpha1"
)

const forTemplateRefKey = "spec.templateRef"
const forInputsObjectKey = "spec.inputs.object.ref"

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
		logger.Error(err, "Get failed")
		err = client.IgnoreNotFound(err)
		return
	}

	// Return early if the object is suspended.
	if tt.Spec.Suspend {
		logger.Info("Reconciliation is suspended for this object")
		return ctrl.Result{}, nil
	}

	err = r.addWatchForKind(ctx, schema.GroupVersionKind{
		Version: "v1",
		Kind:    "ConfigMap",
	}, forTemplateRefKey, r.buildWatchEventHandler(forTemplateRefKey))
	if err != nil {
		return ctrl.Result{}, nil
	}
	for _, me := range tt.Spec.Inputs {
		if me.Object != nil {
			gvk, err2 := me.Object.Ref.GroupVersionKind()
			if err2 != nil {
				err = err2
				return
			}
			err = r.addWatchForKind(ctx, gvk, forInputsObjectKey, r.buildWatchEventHandler(forInputsObjectKey))
			if err != nil {
				return
			}
		}
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

// SetupWithManager sets up the controller with the Manager.
func (r *TextTemplateReconciler) SetupWithManager(mgr ctrl.Manager, concurrent int) error {
	r.watchedKinds = map[schema.GroupVersionKind]bool{}

	// Index the TextTemplate by the objects they are for.
	if err := mgr.GetCache().IndexField(context.TODO(), &templatesv1alpha1.TextTemplate{}, forTemplateRefKey,
		func(object client.Object) []string {
			o := object.(*templatesv1alpha1.TextTemplate)
			ref := r.buildTemplateRef(o)
			if ref == nil {
				return nil
			}

			return []string{BuildRefIndexValue(*ref, "")}
		}); err != nil {
		return fmt.Errorf("failed setting index fields: %w", err)
	}
	if err := mgr.GetCache().IndexField(context.TODO(), &templatesv1alpha1.TextTemplate{}, forInputsObjectKey,
		func(object client.Object) []string {
			o := object.(*templatesv1alpha1.TextTemplate)
			var ret []string
			for _, input := range o.Spec.Inputs {
				if input.Object != nil {
					ret = append(ret, BuildRefIndexValue(input.Object.Ref, o.GetNamespace()))
				}
			}
			return ret
		}); err != nil {
		return fmt.Errorf("failed setting index fields: %w", err)
	}

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

	return nil
}

func (r *TextTemplateReconciler) buildWatchEventHandler(indexField string) handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(func(object client.Object) []reconcile.Request {
		var list templatesv1alpha1.TextTemplateList

		err := r.List(context.Background(), &list, client.MatchingFields{
			indexField: BuildObjectIndexValue(object),
		})
		if err != nil {
			return nil
		}
		var reqs []reconcile.Request
		for _, x := range list.Items {
			reqs = append(reqs, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: x.GetNamespace(),
					Name:      x.GetName(),
				},
			})
		}
		return reqs
	})
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
