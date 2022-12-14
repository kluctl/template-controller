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
	"k8s.io/apimachinery/pkg/runtime"
)

// GitProjectorSpec defines the desired state of GitProjector
type GitProjectorSpec struct {
	// Interval is the interval at which to query the Gitlab API.
	// Defaults to 5m.
	// +optional
	// +kubebuilder:default:="5m"
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ms|s|m|h))+$"
	Interval metav1.Duration `json:"interval"`

	// +optional
	// +kubebuilder:default:=false
	Suspend bool `json:"suspend"`

	// +required
	URL string `json:"url"`

	// +optional
	Reference *GitRef `json:"ref,omitempty"`

	// +optional
	Files []GitFile `json:"files,omitempty"`

	// +optional
	SecretRef *LocalObjectReference `json:"secretRef,omitempty"`
}

type GitFile struct {
	// +required
	Glob string `json:"glob"`

	// +optional
	// +kubebuilder:default:=false
	ParseYaml bool `json:"parseYaml,omitempty"`
}

// GitProjectorStatus defines the observed state of GitProjector
type GitProjectorStatus struct {
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// +optional
	AllRefsHash string `json:"allRefsHash,omitempty"`

	// +optional
	Result []GitProjectorResult `json:"result"`
}

type GitProjectorResult struct {
	// +required
	Reference GitRef `json:"ref"`

	// +required
	Files []GitProjectorResultFile `json:"files"`
}

type GitProjectorResultFile struct {
	// +required
	Path string `json:"path"`

	// +optional
	Raw *string `json:"raw,omitempty"`

	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Parsed []*runtime.RawExtension `json:"parsed,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GitProjector is the Schema for the gitprojectors API
type GitProjector struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GitProjectorSpec   `json:"spec,omitempty"`
	Status GitProjectorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GitProjectorList contains a list of GitProjector
type GitProjectorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GitProjector `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GitProjector{}, &GitProjectorList{})
}
