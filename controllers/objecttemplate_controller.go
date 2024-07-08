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
	"github.com/hashicorp/go-multierror"
	"github.com/kluctl/go-jinja2"
	templatesv1alpha1 "github.com/kluctl/template-controller/api/v1alpha1"
	"io"
	"k8s.io/apimachinery/pkg/api/errors"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sort"
	"strings"
	"sync"
)

const forMatrixObjectKey = "spec.matrix.object.ref"

// ObjectTemplateReconciler reconciles a ObjectTemplate object
type ObjectTemplateReconciler struct {
	BaseTemplateReconciler
}

//+kubebuilder:rbac:groups=templates.kluctl.io,resources=objecttemplates,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=templates.kluctl.io,resources=objecttemplates/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=templates.kluctl.io,resources=objecttemplates/finalizers,verbs=update
//+kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=get;list;watch;impersonate
//+kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch

// Reconcile a resource
func (r *ObjectTemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	logger := log.FromContext(ctx)

	logger.V(1).Info("Starting reconcile")
	defer logger.V(1).Info("Finished reconcile", "err", err)

	var rt templatesv1alpha1.ObjectTemplate
	err = r.Get(ctx, req.NamespacedName, &rt)
	if err != nil {
		logger.Error(err, "Get failed")
		err = client.IgnoreNotFound(err)
		return
	}

	// Add our finalizer if it does not exist
	if !controllerutil.ContainsFinalizer(&rt, templatesv1alpha1.ObjectTemplateFinalizer) {
		patch := client.MergeFrom(rt.DeepCopy())
		controllerutil.AddFinalizer(&rt, templatesv1alpha1.ObjectTemplateFinalizer)
		if err := r.Patch(ctx, &rt, patch, client.FieldOwner(r.FieldManager)); err != nil {
			logger.Error(err, "unable to register finalizer")
			return ctrl.Result{}, err
		}
	}

	// Examine if the object is under deletion
	if !rt.GetDeletionTimestamp().IsZero() {
		return r.finalize(ctx, &rt)
	}

	// Return early if the object is suspended.
	if rt.Spec.Suspend {
		logger.Info("Reconciliation is suspended for this object")
		return ctrl.Result{}, nil
	}

	for _, me := range rt.Spec.Matrix {
		if me.Object != nil {
			gvk, err2 := me.Object.Ref.GroupVersionKind()
			if err2 != nil {
				err = err2
				return
			}
			err = r.addWatchForKind(ctx, gvk, forMatrixObjectKey, r.buildWatchEventHandler(forMatrixObjectKey))
			if err != nil {
				return
			}
		}
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
	err = r.Status().Patch(ctx, &rt, patch, SubResourceFieldOwner(r.FieldManager))
	if err != nil {
		return
	}

	result.RequeueAfter = rt.Spec.Interval.Duration
	return
}

func (r *ObjectTemplateReconciler) multiplyMatrix(matrix []map[string]any, key string, newElems []any) []map[string]any {
	var newMatrix []map[string]any

	for _, m := range matrix {
		for _, e := range newElems {
			newME := map[string]any{}
			for k, v := range m {
				newME[k] = v
			}
			newME[key] = e
			newMatrix = append(newMatrix, newME)
		}
	}

	return newMatrix
}

func (r *ObjectTemplateReconciler) buildMatrixEntries(ctx context.Context, rt *templatesv1alpha1.ObjectTemplate, client client.Client) ([]map[string]any, error) {
	var err error
	var matrixEntries []map[string]any
	matrixEntries = append(matrixEntries, map[string]any{})

	for _, me := range rt.Spec.Matrix {
		var elems []any
		if me.Object != nil {
			elems, err = r.buildObjectInput(ctx, client, rt.GetNamespace(), me.Object.Ref, me.Object.JsonPath, me.Object.ExpandLists, false)
			if err != nil {
				return nil, err
			}
		} else if me.List != nil {
			for _, le := range me.List {
				var e any
				err := yaml.Unmarshal(le.Raw, &e)
				if err != nil {
					return nil, err
				}
				elems = append(elems, e)
			}
		} else {
			return nil, fmt.Errorf("missing matrix value")
		}

		matrixEntries = r.multiplyMatrix(matrixEntries, me.Name, elems)
	}
	return matrixEntries, nil
}

func (r *ObjectTemplateReconciler) doReconcile(ctx context.Context, rt *templatesv1alpha1.ObjectTemplate) error {
	baseVars, err := r.buildBaseVars(rt, "objectTemplate")
	if err != nil {
		return err
	}

	j2, err := NewJinja2()
	if err != nil {
		return err
	}
	defer j2.Close()

	var allResources []*unstructured.Unstructured
	var errs *multierror.Error
	var wg sync.WaitGroup
	var mutex sync.Mutex

	objClient, err := r.getClientForObjects(rt.Spec.ServiceAccountName, rt.GetNamespace())
	if err != nil {
		return err
	}

	matrixEntries, err := r.buildMatrixEntries(ctx, rt, objClient)
	if err != nil {
		return err
	}

	wg.Add(len(matrixEntries))
	for _, matrix := range matrixEntries {
		matrix := matrix
		go func() {
			defer wg.Done()
			vars := runtime.DeepCopyJSON(baseVars)
			MergeMap(vars, map[string]interface{}{
				"matrix": matrix,
			})

			resources, err := r.renderTemplates(j2, rt, vars)
			mutex.Lock()
			defer mutex.Unlock()
			if err != nil {
				errs = multierror.Append(errs, err)
				return
			}

			allResources = append(allResources, resources...)
		}()
	}
	wg.Wait()
	if errs != nil {
		return errs
	}

	for _, x := range allResources {
		rm, err := r.Client.RESTMapper().RESTMapping(x.GroupVersionKind().GroupKind(), x.GroupVersionKind().Version)
		if err != nil {
			return err
		}
		if rm.Scope.Name() == apimeta.RESTScopeNameNamespace && x.GetNamespace() == "" {
			x.SetNamespace(rt.Namespace)
		}
	}

	newAppliedResources := map[templatesv1alpha1.ObjectRef]templatesv1alpha1.AppliedResourceInfo{}
	for _, n := range rt.Status.AppliedResources {
		newAppliedResources[n.Ref.WithoutVersion()] = n
	}

	wg.Add(len(allResources))
	for _, resource := range allResources {
		resource := resource

		go func() {
			defer wg.Done()
			err := r.applyRenderedObject(ctx, objClient, resource)
			mutex.Lock()
			defer mutex.Unlock()

			ari := templatesv1alpha1.AppliedResourceInfo{
				Ref:     templatesv1alpha1.ObjectRefFromObject(resource),
				Success: true,
			}

			if err != nil {
				ari.Success = false
				ari.Error = err.Error()
				errs = multierror.Append(errs, err)
			}
			newAppliedResources[ari.Ref.WithoutVersion()] = ari
		}()
	}
	wg.Wait()

	defer func() {
		rt.Status.AppliedResources = make([]templatesv1alpha1.AppliedResourceInfo, 0, len(newAppliedResources))
		for _, ari := range newAppliedResources {
			rt.Status.AppliedResources = append(rt.Status.AppliedResources, ari)
		}
		sort.Slice(rt.Status.AppliedResources, func(i, j int) bool {
			return rt.Status.AppliedResources[i].Ref.String() < rt.Status.AppliedResources[j].Ref.String()
		})
	}()

	if errs != nil {
		return errs
	}

	err = r.prune(ctx, objClient, rt, allResources, newAppliedResources)
	if err != nil {
		return err
	}

	return nil
}

func (r *ObjectTemplateReconciler) prune(ctx context.Context, objClient client.Client, rt *templatesv1alpha1.ObjectTemplate, allResources []*unstructured.Unstructured, appliedResources map[templatesv1alpha1.ObjectRef]templatesv1alpha1.AppliedResourceInfo) error {
	logger := log.FromContext(ctx)

	if !rt.Spec.Prune {
		return nil
	}

	var errs *multierror.Error
	var wg sync.WaitGroup
	var mutex sync.Mutex

	existingRefs := map[templatesv1alpha1.ObjectRef]templatesv1alpha1.ObjectRef{}
	for _, resource := range allResources {
		ref := templatesv1alpha1.ObjectRefFromObject(resource)
		existingRefs[ref.WithoutVersion()] = ref
	}

	var deleted []templatesv1alpha1.ObjectRef
	for _, ari := range appliedResources {
		ari := ari
		if _, ok := existingRefs[ari.Ref.WithoutVersion()]; ok {
			continue
		}

		logger.Info("Deleting object", "ref", ari.Ref)

		gvk, err := ari.Ref.GroupVersionKind()
		if err != nil {
			return err
		}
		m := metav1.PartialObjectMetadata{}
		m.SetGroupVersionKind(gvk)
		m.SetNamespace(ari.Ref.Namespace)
		m.SetName(ari.Ref.Name)

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := objClient.Delete(ctx, &m)
			mutex.Lock()
			defer mutex.Unlock()
			if err != nil {
				if errors.IsNotFound(err) {
					err = nil
				} else {
					errs = multierror.Append(errs, err)
				}
			}
			if err == nil {
				deleted = append(deleted, ari.Ref)
			}
		}()
	}
	wg.Wait()

	for _, ref := range deleted {
		delete(appliedResources, ref.WithoutVersion())
	}

	return errs.ErrorOrNil()
}

func (r *ObjectTemplateReconciler) applyRenderedObject(ctx context.Context, objClient client.Client, rendered *unstructured.Unstructured) error {
	logger := log.FromContext(ctx)

	var origMeta metav1.PartialObjectMetadata
	origObjFound := false
	origMeta.SetGroupVersionKind(rendered.GroupVersionKind())

	err := objClient.Get(ctx, client.ObjectKeyFromObject(rendered), &origMeta)
	if err == nil {
		origObjFound = true
	}

	err = objClient.Patch(ctx, rendered, client.Apply, client.FieldOwner(r.FieldManager))
	if err != nil {
		return err
	}

	if !origObjFound {
		logger.Info("Created new object", "ref", templatesv1alpha1.ObjectRefFromObject(rendered))
	} else {
		if origMeta.GetResourceVersion() != rendered.GetResourceVersion() {
			logger.Info("Updated existing object", "ref", templatesv1alpha1.ObjectRefFromObject(rendered))
		}
	}

	return nil
}

func (r *ObjectTemplateReconciler) renderTemplates(j2 *jinja2.Jinja2, rt *templatesv1alpha1.ObjectTemplate, vars map[string]any) ([]*unstructured.Unstructured, error) {
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

// SetupWithManager sets up the controller with the Manager.
func (r *ObjectTemplateReconciler) SetupWithManager(mgr ctrl.Manager, concurrent int) error {
	r.Manager = mgr

	// Index the ObjectTemplate by the objects they are for.
	if err := mgr.GetCache().IndexField(context.TODO(), &templatesv1alpha1.ObjectTemplate{}, forMatrixObjectKey,
		func(object client.Object) []string {
			o := object.(*templatesv1alpha1.ObjectTemplate)
			var ret []string
			for _, me := range o.Spec.Matrix {
				if me.Object != nil {
					ret = append(ret, BuildRefIndexValue(me.Object.Ref, o.GetNamespace()))
				}
			}
			return ret
		}); err != nil {
		return fmt.Errorf("failed setting index fields: %w", err)
	}

	c, err := ctrl.NewControllerManagedBy(mgr).
		For(&templatesv1alpha1.ObjectTemplate{}, builder.WithPredicates(
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

func (r *ObjectTemplateReconciler) buildWatchEventHandler(indexField string) handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, object client.Object) []reconcile.Request {
		var list templatesv1alpha1.ObjectTemplateList

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

func (r *ObjectTemplateReconciler) finalize(ctx context.Context, obj *templatesv1alpha1.ObjectTemplate) (ctrl.Result, error) {
	r.doFinalize(ctx, obj)

	// Remove our finalizer from the list and update it
	controllerutil.RemoveFinalizer(obj, templatesv1alpha1.ObjectTemplateFinalizer)
	if err := r.Update(ctx, obj, client.FieldOwner(r.FieldManager)); err != nil {
		return ctrl.Result{}, err
	}

	// Stop reconciliation as the object is being deleted
	return ctrl.Result{}, nil
}

func (r *ObjectTemplateReconciler) doFinalize(ctx context.Context, obj *templatesv1alpha1.ObjectTemplate) {
	log := ctrl.LoggerFrom(ctx)

	if !obj.Spec.Prune || obj.Spec.Suspend {
		return
	}

	objClient, err := r.getClientForObjects(obj.Spec.ServiceAccountName, obj.GetNamespace())
	if err != nil {
		log.Error(err, "Failed to create objClient for deletion")
		return
	}

	var wg sync.WaitGroup
	for _, ar := range obj.Status.AppliedResources {
		ar := ar
		wg.Add(1)
		go func() {
			defer wg.Done()
			gvk, err := ar.Ref.GroupVersionKind()
			if err != nil {
				return
			}

			log.Info("Deleting applied object", "ref", ar.Ref)

			var o unstructured.Unstructured
			o.SetGroupVersionKind(gvk)
			o.SetName(ar.Ref.Name)
			o.SetNamespace(ar.Ref.Namespace)
			err = objClient.Delete(ctx, &o)
			if err != nil && !errors.IsNotFound(err) {
				log.Error(err, "Failed to delete applied object", "ref", ar.Ref)
			}
		}()
	}
	wg.Wait()
}
