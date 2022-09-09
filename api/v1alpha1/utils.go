package v1alpha1

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Utility struct for a reference to a secret key.
type SecretRef struct {
	SecretName string `json:"secretName"`
	Key        string `json:"key"`
}

type ResourceRef struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
}

func ResourceRefFromObject(object client.Object) ResourceRef {
	gvk := object.GetObjectKind().GroupVersionKind()
	return ResourceRef{
		APIVersion: gvk.GroupVersion().String(),
		Kind:       gvk.Kind,
		Namespace:  object.GetNamespace(),
		Name:       object.GetName(),
	}
}

func (r *ResourceRef) GroupVersionKind() (schema.GroupVersionKind, error) {
	gv, err := schema.ParseGroupVersion(r.APIVersion)
	if err != nil {
		return schema.GroupVersionKind{}, err
	}

	return schema.GroupVersionKind{
		Group:   gv.Group,
		Version: gv.Version,
		Kind:    r.Kind,
	}, nil
}

func (r *ResourceRef) String() string {
	if r.Namespace != "" {
		return fmt.Sprintf("%s/%s/%s", r.Namespace, r.Kind, r.Name)
	} else {
		if r.Name != "" {
			return fmt.Sprintf("%s/%s", r.Kind, r.Name)
		} else {
			return r.Kind
		}
	}
}
