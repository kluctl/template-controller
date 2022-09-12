package handlers

import (
	"context"
	"fmt"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/cli-utils/pkg/kstatus/status"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

type StatusCalculator struct {
	Client client.Client
}

func (sc *StatusCalculator) ComputeReady(ctx context.Context, obj client.Object, missingReadyConditionIsError bool) (bool, error) {
	u1, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return false, err
	}
	u := &unstructured.Unstructured{Object: u1}

	if sc.hasStatus(ctx, u) {
		if _, ok := u.Object["status"]; !ok {
			// expected status, but it is missing
			return false, nil
		}
	}
	if missingReadyConditionIsError {
		uc, err := status.GetObjectWithConditions(u1)
		if err != nil {
			return false, err
		}
		found := false
		for _, cond := range uc.Status.Conditions {
			if cond.Type == "Ready" {
				found = true
				break
			}
		}
		if !found {
			return false, fmt.Errorf("missing Ready condition")
		}
	}

	res, err := status.Compute(u)
	if err != nil {
		return false, err
	}

	return res.Status == status.CurrentStatus, nil
}

func (sc *StatusCalculator) hasStatus(ctx context.Context, obj client.Object) bool {
	testObj, err := sc.Client.Scheme().New(obj.GetObjectKind().GroupVersionKind())
	if err != nil {
		return sc.hasStatusByCrd(ctx, obj)
	}
	t := reflect.TypeOf(testObj)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("json")
		if tag == "" {
			continue
		}
		tagS := strings.Split(tag, ",")
		if len(tagS) > 0 && tagS[0] == "status" {
			return true
		}
	}
	return false
}

func (sc *StatusCalculator) getSchemaForGVK(ctx context.Context, gvk schema.GroupVersionKind) (*apiextensionsv1.JSONSchemaProps, error) {
	rm, err := sc.Client.RESTMapper().RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, err
	}
	crdName := fmt.Sprintf("%s.%s", rm.Resource.Resource, gvk.Group)

	var crd apiextensionsv1.CustomResourceDefinition
	err = sc.Client.Get(ctx, types.NamespacedName{Name: crdName}, &crd)
	if err != nil {
		return nil, err
	}

	for _, v := range crd.Spec.Versions {
		if v.Name != gvk.Version {
			continue
		}
		return v.Schema.OpenAPIV3Schema, nil
	}
	return nil, fmt.Errorf("schema for %s not found", gvk.String())
}

func (sc *StatusCalculator) hasStatusByCrd(ctx context.Context, obj client.Object) bool {
	s, err := sc.getSchemaForGVK(ctx, obj.GetObjectKind().GroupVersionKind())
	if err != nil {
		return false
	}
	found := false
	for n, p := range s.Properties {
		if n != "status" {
			continue
		}
		if p.Type != "object" {
			continue
		}
		found = true
		break
	}
	if found {
		return true
	}
	return false
}
