package controllers

import (
	"context"
	"fmt"
	templatesv1alpha1 "github.com/kluctl/template-controller/api/v1alpha1"
	"github.com/ohler55/ojg/jp"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sync"
)

type BaseTemplateReconciler struct {
	client.Client

	Manager      manager.Manager
	Scheme       *runtime.Scheme
	FieldManager string

	controller   controller.Controller
	watchedKinds map[schema.GroupVersionKind]bool
	mutex        sync.Mutex
}

func (r *BaseTemplateReconciler) getClientForObjects(serviceAccountName string, objNamespace string) (client.Client, error) {
	restConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	name := "default"
	if serviceAccountName != "" {
		name = serviceAccountName
	}
	if name == "" {
		return nil, fmt.Errorf("empty serviceAccountName not allowed")
	}
	username := fmt.Sprintf("system:serviceaccount:%s:%s", objNamespace, name)
	restConfig.Impersonate = rest.ImpersonationConfig{UserName: username}

	c, err := client.New(restConfig, client.Options{Mapper: r.RESTMapper()})
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *BaseTemplateReconciler) addWatchForKind(ctx context.Context, gvk schema.GroupVersionKind, key string, eventHandler handler.TypedEventHandler[client.Object]) error {
	logger := log.FromContext(ctx)

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.watchedKinds == nil {
		r.watchedKinds = map[schema.GroupVersionKind]bool{}
	}

	keyGvk := gvk
	keyGvk.Kind += "+" + key
	if x, ok := r.watchedKinds[keyGvk]; ok && x {
		return nil
	}

	logger.V(1).Info("Starting watch for new kind and key", "gvk", gvk, "key", key)

	var dummy unstructured.Unstructured
	dummy.SetGroupVersionKind(gvk)

	err := r.controller.Watch(source.Kind[client.Object](r.Manager.GetCache(), &dummy, eventHandler))
	if err != nil {
		return err
	}

	r.watchedKinds[keyGvk] = true
	return nil
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
