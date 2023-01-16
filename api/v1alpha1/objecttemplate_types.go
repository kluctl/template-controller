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
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ms|s|m|h))+$"
	Interval metav1.Duration `json:"interval"`

	// Suspend can be used to suspend the reconciliation of this object
	// +optional
	// +kubebuilder:default:=false
	Suspend bool `json:"suspend"`

	// ServiceAccountName specifies the name of the Kubernetes service account to impersonate
	// when reconciling this ObjectTemplate. If omitted, the "default" service account is used
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`

	// Prune enables pruning of previously created objects when these disappear from the list of rendered objects
	// +kubebuilder:default:=false
	// +optional
	Prune bool `json:"prune"`

	// Matrix specifies the input matrix
	// +required
	Matrix []*MatrixEntry `json:"matrix"`

	// Templates specifies a list of templates to render and deploy
	// +required
	Templates []Template `json:"templates"`
}

type MatrixEntry struct {
	// Name specifies the name this matrix input is available while rendering templates
	// +required
	Name string `json:"name"`

	// Object specifies an object to load and make available while rendering templates. The object can be accessed
	// through the name specified above. The service account used by the ObjectTemplate must have proper permissions
	// to get this object
	// +optional
	Object *MatrixEntryObject `json:"object,omitempty"`

	// List specifies a list of plain YAML values which are made available while rendering templates. The list can be
	// accessed through the name specified above
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	List []runtime.RawExtension `json:"list,omitempty"`
}

type MatrixEntryObject struct {
	// Ref specifies the apiVersion, kind, namespace and name of the object to load. The service account used by the
	// ObjectTemplate must have proper permissions to get this object
	// +required
	Ref ObjectRef `json:"ref"`

	// JsonPath optionally specifies a sub-field to load. When specified, the sub-field (and not the whole object)
	// is made available while rendering templates
	// +optional
	JsonPath *string `json:"jsonPath,omitempty"`

	// ExpandLists enables optional expanding of list. Expanding means, that each list entry is interpreted as
	// individual matrix input instead of interpreting the whole list as one matrix input. This feature is only useful
	// when used in combination with `jsonPath`
	// +optional
	ExpandLists bool `json:"expandLists,omitempty"`
}

type Template struct {
	// Object specifies a structured object in YAML form. Each field value is rendered independently.
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Object *unstructured.Unstructured `json:"object,omitempty"`

	// Raw specifies a raw string to be interpreted/parsed as YAML. The whole string is rendered in one go, allowing to
	// use advanced Jinja2 control structures. Raw object might also be required when a templated value must not be
	// interpreted as a string (which would be done in Object).
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
