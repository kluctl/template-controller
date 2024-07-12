package controllers

import (
	"context"
	"fmt"
	templatesv1alpha1 "github.com/kluctl/template-controller/api/v1alpha1"
	"github.com/ohler55/ojg/jp"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sync"
)

type BaseTemplateReconciler struct {
	client.Client

	RawWatchContext context.Context

	Manager      manager.Manager
	Scheme       *runtime.Scheme
	FieldManager string

	controller controller.Controller

	watchesUtil watchesUtil

	mutex sync.Mutex
}

func (r *BaseTemplateReconciler) getClientForObjects(serviceAccountName string, objNamespace string) (client.WithWatch, string, error) {
	restConfig := rest.CopyConfig(r.Manager.GetConfig())

	name := "default"
	if serviceAccountName != "" {
		name = serviceAccountName
	}
	username := fmt.Sprintf("system:serviceaccount:%s:%s", objNamespace, name)
	restConfig.Impersonate = rest.ImpersonationConfig{UserName: username}

	c, err := client.NewWithWatch(restConfig, client.Options{Mapper: r.RESTMapper()})
	if err != nil {
		return nil, "", err
	}
	return c, name, nil
}

func (r *BaseTemplateReconciler) buildBaseVars(templateObj runtime.Object, objVarName string) (map[string]any, error) {
	vars := map[string]any{}

	u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(templateObj)
	if err != nil {
		return nil, err
	}

	vars[objVarName] = u
	return vars, nil
}

func (r *BaseTemplateReconciler) buildObjectInput(ctx context.Context, client client.Client, objNamespace string, ref templatesv1alpha1.ObjectRef, jsonPath *string, expandLists bool, expectOne bool) ([]any, error) {
	gvk, err := ref.GroupVersionKind()
	if err != nil {
		return nil, err
	}
	namespace := objNamespace
	if ref.Namespace != "" {
		namespace = ref.Namespace
	}

	var o unstructured.Unstructured
	o.SetGroupVersionKind(gvk)

	err = client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: ref.Name}, &o)
	if err != nil {
		return nil, err
	}

	var results []any

	if jsonPath != nil {
		jp, err := jp.ParseString(*jsonPath)
		if err != nil {
			return nil, err
		}
		results = jp.Get(o.Object)
	} else {
		results = []any{o.Object}
	}

	var elems []any
	for _, x := range results {
		if expandLists {
			if l, ok := x.([]any); ok {
				elems = append(elems, l...)
			} else {
				elems = append(elems, x)
			}
		} else {
			elems = append(elems, x)
		}
	}

	if expectOne {
		if len(elems) == 0 {
			return nil, fmt.Errorf("failed to get object/subElement %s: %w", ref.String(), err)
		}
		if len(elems) > 1 {
			return nil, fmt.Errorf("more than one element returned for object %s and json path %s: %w", ref.String(), *jsonPath, err)
		}
	}

	return elems, nil
}
