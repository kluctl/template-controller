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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	ObjectTemplateFinalizer = "finalizers.templates.kluctl.io"
)

// ObjectTemplateSpec defines the desired state of ObjectTemplate
type ObjectTemplateSpec struct {
	// +kubebuilder:default:="30s"
	Interval metav1.Duration `json:"interval"`

	// +optional
	// +kubebuilder:default:=false
	Suspend bool `json:"suspend"`

	// The name of the Kubernetes service account to impersonate
	// when reconciling this ObjectTemplate. If omitted, the "default" service account is used.
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`

	// +kubebuilder:default:=false
	// +optional
	Prune bool `json:"prune"`

	// +required
	Matrix []*MatrixEntry `json:"matrix"`

	// +required
	Templates []Template `json:"templates"`
}

type MatrixEntry struct {
	// +required
	Name string `json:"name"`

	// +optional
	Object *MatrixEntryObject `json:"object,omitempty"`

	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	List []runtime.RawExtension `json:"list,omitempty"`
}

type MatrixEntryObject struct {
	// +required
	Ref ObjectRef `json:"ref"`

	// +optional
	JsonPath *string `json:"jsonPath,omitempty"`

	// +optional
	ExpandLists bool `json:"expandLists,omitempty"`
}

type Template struct {
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Object *unstructured.Unstructured `json:"object,omitempty"`

	// +optional
	Raw *string `json:"raw,omitempty"`
}

// ObjectTemplateStatus defines the observed state of ObjectTemplate
type ObjectTemplateStatus struct {
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// +optional
	AppliedResources []AppliedResourceInfo `json:"appliedResources,omitempty"`
}

type AppliedResourceInfo struct {
	Ref ObjectRef `json:"ref"`

	Success bool `json:"success"`

	// +optional
	Error string `json:"error,omitempty"`
}

// GetConditions returns the status conditions of the object.
func (in *ObjectTemplate) GetConditions() []metav1.Condition {
	return in.Status.Conditions
}

// SetConditions sets the status conditions on the object.
func (in *ObjectTemplate) SetConditions(conditions []metav1.Condition) {
	in.Status.Conditions = conditions
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ObjectTemplate is the Schema for the objecttemplates API
type ObjectTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ObjectTemplateSpec   `json:"spec,omitempty"`
	Status ObjectTemplateStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ObjectTemplateList contains a list of ObjectTemplate
type ObjectTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ObjectTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ObjectTemplate{}, &ObjectTemplateList{})
}
