package controllers

import (
	"fmt"
	templatesv1alpha1 "github.com/kluctl/template-controller/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func BuildRefIndexValue(ref templatesv1alpha1.ObjectRef, ns string) string {
	if ref.Namespace != "" {
		ns = ref.Namespace
	}
	return fmt.Sprintf("%s/%s/%s", ref.Kind, ns, ref.Name)
}

func BuildObjectIndexValue(obj client.Object) string {
	gvk := obj.GetObjectKind().GroupVersionKind()
	return fmt.Sprintf("%s/%s/%s", gvk.Kind, obj.GetNamespace(), obj.GetName())
}
