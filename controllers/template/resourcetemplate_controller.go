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

package template

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/kluctl/go-jinja2"
	"github.com/kluctl/template-controller/controllers"
	generators2 "github.com/kluctl/template-controller/controllers/template/generators"
	"io"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
	"strings"

	templatesv1alpha1 "github.com/kluctl/template-controller/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
//+kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=get;list;watch

// Reconcile a resource
func (r *ResourceTemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
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

	j2, err := controllers.NewJinja2()
	if err != nil {
		return err
	}
	defer j2.Close()

	var allResources []*unstructured.Unstructured

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
			vars := runtime.DeepCopyJSON(baseVars)
			controllers.MergeMap(vars, c.Vars)

			resources, err := r.renderTemplates(ctx, j2, rt, vars)
			if err != nil {
				return err
			}
			allResources = append(allResources, resources...)
		}
	}

	toDelete := make(map[templatesv1alpha1.ObjectRef]templatesv1alpha1.ObjectRef)
	for _, n := range rt.Status.AppliedResources {
		gvk, err := n.Ref.GroupVersionKind()
		if err != nil {
			return err
		}
		ref := n.Ref
		ref.APIVersion = gvk.Group
		toDelete[ref] = n.Ref
	}

	rt.Status.AppliedResources = nil

	var errs *multierror.Error
	for _, resource := range allResources {
		ref := templatesv1alpha1.ObjectRefFromObject(resource)
		gvk, err := ref.GroupVersionKind()
		if err != nil {
			return err
		}

		ari := templatesv1alpha1.AppliedResourceInfo{
			Ref:     ref,
			Success: true,
		}

		ref.APIVersion = gvk.Group
		delete(toDelete, ref)

		err = r.applyTemplate(ctx, rt, resource)
		if err != nil {
			ari.Success = false
			ari.Error = err.Error()
			errs = multierror.Append(errs, err)
		}

		rt.Status.AppliedResources = append(rt.Status.AppliedResources, ari)
	}

	for _, ref := range toDelete {
		gvk, err := ref.GroupVersionKind()
		if err != nil {
			return err
		}
		m := metav1.PartialObjectMetadata{}
		m.SetGroupVersionKind(gvk)
		m.SetNamespace(ref.Namespace)
		m.SetName(ref.Name)

		err = r.Delete(ctx, &m)
		if err != nil {
			if !errors.IsNotFound(err) {
				errs = multierror.Append(errs, err)
			}
		}
	}

	return errs.ErrorOrNil()
}

func (r *ResourceTemplateReconciler) applyTemplate(ctx context.Context, rt *templatesv1alpha1.ResourceTemplate, rendered *unstructured.Unstructured) error {
	log := ctrl.LoggerFrom(ctx)

	x := rendered.DeepCopy()

	mres, err := controllerutil.CreateOrUpdate(ctx, r.Client, x, func() error {
		if err := controllerutil.SetControllerReference(rt, x, r.Scheme); err != nil {
			return err
		}
		controllers.MergeMap(x.Object, rendered.Object)
		return nil
	})
	if err != nil {
		return err
	}

	if mres != controllerutil.OperationResultNone {
		log.Info(fmt.Sprintf("CreateOrUpdate returned %v", mres), "ref", templatesv1alpha1.ObjectRefFromObject(rendered))
	}
	return nil
}

func (r *ResourceTemplateReconciler) renderTemplates(ctx context.Context, j2 *jinja2.Jinja2, rt *templatesv1alpha1.ResourceTemplate, vars map[string]any) ([]*unstructured.Unstructured, error) {
	var ret []*unstructured.Unstructured
	for _, t := range rt.Spec.Templates {
		if t.Object != nil {
			x := t.Object.DeepCopy()
			_, err := j2.RenderStruct(x, jinja2.WithGlobals(vars))
			if err != nil {
				return nil, err
			}
			ret = append(ret, x)
		} else if t.Raw != nil {
			r, err := j2.RenderString(*t.Raw, jinja2.WithGlobals(vars))
			if err != nil {
				return nil, err
			}
			d := yaml.NewYAMLToJSONDecoder(strings.NewReader(r))
			for {
				var u unstructured.Unstructured
				err = d.Decode(&u)
				if err != nil {
					if err == io.EOF {
						break
					}
					return nil, err
				}
				ret = append(ret, &u)
			}
		} else {
			return nil, fmt.Errorf("no template specified")
		}
	}
	return ret, nil
}

func (r *ResourceTemplateReconciler) buildBaseVars(ctx context.Context, rt *templatesv1alpha1.ResourceTemplate) (map[string]any, error) {
	vars := map[string]any{}

	u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(rt)
	if err != nil {
		return nil, err
	}

	vars["resourceTemplate"] = u
	return vars, nil
}

func (r *ResourceTemplateReconciler) buildGenerator(ctx context.Context, namespace string, spec templatesv1alpha1.Generator) (generators2.Generator, error) {
	if spec.PullRequest != nil {
		return generators2.BuildPullRequestGenerator(ctx, r.Client, namespace, *spec.PullRequest)
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
