package v1alpha1

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type LocalObjectReference struct {
	// Name of the referent.
	// +required
	Name string `json:"name"`
}

// Utility struct for a reference to a secret key.
type SecretRef struct {
	SecretName string `json:"secretName"`
	Key        string `json:"key"`
}

type ConfigMapRef struct {
	Name string `json:"name"`
	Key  string `json:"key"`
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

func (r *ObjectRef) WithoutVersion() ObjectRef {
	gv, err := schema.ParseGroupVersion(r.APIVersion)
	if err != nil {
		return *r
	}
	return ObjectRef{
		APIVersion: gv.Group,
		Kind:       r.Kind,
		Namespace:  r.Namespace,
		Name:       r.Name,
	}
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

type GitRef struct {
	// Branch to filter for. Can also be a regex.
	// +optional
	Branch string `json:"branch,omitempty"`

	// Tag to filter for. Can also be a regex.
	// +optional
	Tag string `json:"tag,omitempty"`

	// Commit SHA to check out, takes precedence over all reference fields.
	// +optional
	Commit string `json:"commit,omitempty"`
}

func (gr *GitRef) Less(o GitRef) bool {
	s1 := fmt.Sprintf("%s+%s+%s", gr.Commit, gr.Tag, gr.Branch)
	s2 := fmt.Sprintf("%s+%s+%s", o.Commit, o.Tag, o.Branch)
	return s1 < s2
}
