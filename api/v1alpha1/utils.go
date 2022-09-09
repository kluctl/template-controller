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

type ObjectRef struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`

	// +optional
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name"`
}

func ObjectRefFromObject(object client.Object) ObjectRef {
	gvk := object.GetObjectKind().GroupVersionKind()
	return ObjectRef{
		APIVersion: gvk.GroupVersion().String(),
		Kind:       gvk.Kind,
		Namespace:  object.GetNamespace(),
		Name:       object.GetName(),
	}
}

func (r *ObjectRef) GroupVersionKind() (schema.GroupVersionKind, error) {
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

func (r *ObjectRef) String() string {
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
