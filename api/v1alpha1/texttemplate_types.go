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
)

const (
	TextTemplateFinalizer = "finalizers.templates.kluctl.io"
)

// TextTemplateSpec defines the desired state of TextTemplate
type TextTemplateSpec struct {
	// Suspend can be used to suspend the reconciliation of this object.
	// +optional
	// +kubebuilder:default:=false
	Suspend bool `json:"suspend"`

	// The name of the Kubernetes service account to impersonate
	// when reconciling this TextTemplate. If omitted, the "default" service account is used.
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`

	// +optional
	Inputs []TextTemplateInput `json:"inputs,omitempty"`

	// +optional
	Template *string `json:"template,omitempty"`

	// +optional
	TemplateRef *TemplateRef `json:"templateRef,omitempty"`
}

type TextTemplateInput struct {
	// +required
	Name string `json:"name"`

	// +optional
	Object *TextTemplateInputObject `json:"object,omitempty"`
}

type TextTemplateInputObject struct {
	// +required
	Ref ObjectRef `json:"ref"`

	// +optional
	JsonPath *string `json:"jsonPath,omitempty"`
}

type TemplateRef struct {
	// +optional
	ConfigMap *TemplateRefConfigMap `json:"configMap,omitempty"`
}

type TemplateRefConfigMap struct {
	// +required
	Name string `json:"name"`

	// +optional
	Namespace string `json:"namespace,omitempty"`

	// +required
	Key string `json:"key"`
}

// TextTemplateStatus defines the observed state of TextTemplate
type TextTemplateStatus struct {
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// +optional
	Result string `json:"result,omitempty"`
}

// GetConditions returns the status conditions of the object.
func (in *TextTemplate) GetConditions() []metav1.Condition {
	return in.Status.Conditions
}

// SetConditions sets the status conditions on the object.
func (in *TextTemplate) SetConditions(conditions []metav1.Condition) {
	in.Status.Conditions = conditions
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TextTemplate is the Schema for the texttemplates API
type TextTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TextTemplateSpec   `json:"spec,omitempty"`
	Status TextTemplateStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TextTemplateList contains a list of TextTemplate
type TextTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TextTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TextTemplate{}, &TextTemplateList{})
}
