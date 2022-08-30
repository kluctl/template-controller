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
	"github.com/kluctl/kluctl/v2/pkg/jinja2"
	"github.com/kluctl/kluctl/v2/pkg/types/k8s"
	"github.com/kluctl/kluctl/v2/pkg/utils"
	"github.com/kluctl/kluctl/v2/pkg/utils/uo"
	"k8s.io/apimachinery/pkg/api/errors"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/json"
	templatesv1alpha1 "kluctl/template-controller/api/v1alpha1"
	"kluctl/template-controller/controllers/generators"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ResourceTemplateReconciler reconciles a ResourceTemplate object
type ResourceTemplateReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=templates.kluctl.io,resources=resourcetemplates,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=templates.kluctl.io,resources=resourcetemplates/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=templates.kluctl.io,resources=resourcetemplates/finalizers,verbs=update

// Reconcile a resource
func (r *ResourceTemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = ctrl.LoggerFrom(ctx)

	var rt templatesv1alpha1.ResourceTemplate
	err := r.Get(ctx, req.NamespacedName, &rt)
	if err != nil {
		return ctrl.Result{}, err
	}

	patch := client.MergeFrom(rt.DeepCopy())
	err = r.doReconcile(ctx, &rt)
	if err != nil {
		c := metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionFalse,
			ObservedGeneration: rt.GetGeneration(),
			Reason:             "Error",
			Message:            err.Error(),
		}
		apimeta.SetStatusCondition(&rt.Status.Conditions, c)
	} else {
		c := metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionTrue,
			ObservedGeneration: rt.GetGeneration(),
			Reason:             "Success",
			Message:            "Success",
		}
		apimeta.SetStatusCondition(&rt.Status.Conditions, c)
	}
	err = r.Status().Patch(ctx, &rt, patch)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{
		RequeueAfter: rt.Spec.Interval.Duration,
	}, nil
}

func (r *ResourceTemplateReconciler) doReconcile(ctx context.Context, rt *templatesv1alpha1.ResourceTemplate) error {
	baseVars, err := r.buildBaseVars(ctx, rt)
	if err != nil {
		return err
	}

	j2, err := jinja2.NewJinja2()
	if err != nil {
		return err
	}
	defer j2.Close()

	var allResources []*uo.UnstructuredObject

	for _, g := range rt.Spec.Generators {
		g, err := r.buildGenerator(ctx, rt.GetNamespace(), g)
		if err != nil {
			return err
		}
		contexts, err := g.BuildContexts()
		if err != nil {
			return err
		}

		for _, c := range contexts {
			vars := baseVars.MergeCopy(c.Vars)

			resources, err := r.renderTemplates(ctx, j2, rt, vars)
			if err != nil {
				return err
			}
			allResources = append(allResources, resources...)
		}
	}

	toDelete := make(map[k8s.ObjectRef]k8s.ObjectRef)
	for _, n := range rt.Status.AppliedResources {
		ref := k8s.NewObjectRef(n.Group, n.Version, n.Kind, n.Name, n.Namespace)
		ref2 := ref
		ref2.GVK.Version = ""
		toDelete[ref2] = ref
	}

	rt.Status.AppliedResources = nil

	var errs []error
	for _, resource := range allResources {
		ref := resource.GetK8sRef()

		ari := templatesv1alpha1.AppliedResourceInfo{
			Group:     ref.GVK.Group,
			Version:   ref.GVK.Version,
			Kind:      ref.GVK.Kind,
			Namespace: ref.Namespace,
			Name:      ref.Name,
			Success:   true,
		}

		ref.GVK.Version = ""
		delete(toDelete, ref)

		err = r.applyTemplate(ctx, rt, resource)
		if err != nil {
			ari.Success = false
			ari.Error = err.Error()
			errs = append(errs, err)
		}

		rt.Status.AppliedResources = append(rt.Status.AppliedResources, ari)
	}

	for _, ref := range toDelete {
		m := metav1.PartialObjectMetadata{}
		m.SetGroupVersionKind(ref.GVK)
		m.SetNamespace(ref.Namespace)
		m.SetName(ref.Name)

		err = r.Delete(ctx, &m)
		if err != nil {
			if !errors.IsNotFound(err) {
				errs = append(errs, err)
			}
		}
	}

	return utils.NewErrorListOrNil(errs)
}

func (r *ResourceTemplateReconciler) applyTemplate(ctx context.Context, rt *templatesv1alpha1.ResourceTemplate, rendered *uo.UnstructuredObject) error {
	log := ctrl.LoggerFrom(ctx)

	n := rendered.Clone()
	mres, err := controllerutil.CreateOrUpdate(ctx, r.Client, n.ToUnstructured(), func() error {
		if err := controllerutil.SetControllerReference(rt, n.ToUnstructured(), r.Scheme); err != nil {
			return err
		}
		for k, v := range rendered.GetK8sAnnotations() {
			n.SetK8sAnnotation(k, v)
		}
		for k, v := range rendered.GetK8sLabels() {
			n.SetK8sLabel(k, v)
		}
		for k := range rendered.Object {
			if k != "metadata" && k != "status" {
				n.Object[k] = rendered.Object[k]
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	if mres != controllerutil.OperationResultNone {
		log.Info(fmt.Sprintf("CreateOrUpdate returned %v", mres), "ref", rendered.GetK8sRef())
	}
	return nil
}

func (r *ResourceTemplateReconciler) renderTemplates(ctx context.Context, j2 *jinja2.Jinja2, rt *templatesv1alpha1.ResourceTemplate, vars *uo.UnstructuredObject) ([]*uo.UnstructuredObject, error) {
	var ret []*uo.UnstructuredObject
	for _, t := range rt.Spec.Templates {
		b, err := t.MarshalJSON()
		if err != nil {
			return nil, err
		}

		r, err := j2.RenderString(string(b), nil, vars)
		if err != nil {
			return nil, err
		}

		u, err := uo.FromString(r)
		if err != nil {
			return nil, err
		}
		ret = append(ret, u)
	}
	return ret, nil
}

func (r *ResourceTemplateReconciler) buildBaseVars(ctx context.Context, rt *templatesv1alpha1.ResourceTemplate) (*uo.UnstructuredObject, error) {
	vars := uo.New()

	b, err := json.Marshal(rt)
	if err != nil {
		return nil, err
	}
	u, err := uo.FromString(string(b))
	if err != nil {
		return nil, err
	}

	_ = vars.SetNestedField(u.Object, "resourceTemplate")

	return vars, nil
}

func (r *ResourceTemplateReconciler) buildGenerator(ctx context.Context, namespace string, spec templatesv1alpha1.Generator) (generators.Generator, error) {
	if spec.PullRequest != nil {
		return generators.BuildPullRequestGenerator(ctx, r.Client, namespace, *spec.PullRequest)
	} else {
		return nil, fmt.Errorf("no generator specified")
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *ResourceTemplateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&templatesv1alpha1.ResourceTemplate{}).
		Complete(r)
}
