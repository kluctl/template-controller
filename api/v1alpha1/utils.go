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
	Group     string `json:"group"`
	Version   string `json:"version"`
	Kind      string `json:"kind"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

func ResourceRefFromObject(object client.Object) ResourceRef {
	gvk := object.GetObjectKind().GroupVersionKind()
	return ResourceRef{
		Group:     gvk.Group,
		Version:   gvk.Version,
		Kind:      gvk.Kind,
		Namespace: object.GetNamespace(),
		Name:      object.GetName(),
	}
}

func (r *ResourceRef) GroupVersionLind() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   r.Group,
		Version: r.Version,
		Kind:    r.Kind,
	}
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
