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
)

// ResourcesTemplateSpec defines the desired state of ResourcesTemplate
type ResourcesTemplateSpec struct {
	// +kubebuilder:default:="30s"
	Interval metav1.Duration `json:"interval"`

	// +required
	Generators []Generator `json:"generators"`

	// +required
	// +kubebuilder:pruning:PreserveUnknownFields
	Templates []unstructured.Unstructured `json:"templates"`
}

type Generator struct {
	// +optional
	PullRequest *PullRequestGenerator `json:"pullRequest,omitempty"`
}

type PullRequestGenerator struct {
	// +optional
	Gitlab *PullRequestGeneratorGitlab `json:"gitlab"`

	// Filters for which pull requests should be considered.
	Filters []PullRequestGeneratorFilter `json:"filters,omitempty"`
}

// PullRequestGeneratorFilter is a single pull request filter.
// If multiple filter types are set on a single struct, they will be AND'd together. All filters must
// pass for a pull request to be included.
type PullRequestGeneratorFilter struct {
	BranchMatch *string `json:"branchMatch,omitempty"`
}

type PullRequestGeneratorGitlab struct {
	// GitLab project to scan. Required.
	Project string `json:"project"`
	// The GitLab API URL to talk to. If blank, uses https://gitlab.com/.
	API string `json:"api,omitempty"`
	// Authentication token reference.
	TokenRef *SecretRef `json:"tokenRef,omitempty"`
	// Labels is used to filter the MRs that you want to target
	Labels []string `json:"labels,omitempty"`
	// PullRequestState is an additional MRs filter to get only those with a certain state. Default: "" (all states)
	PullRequestState string `json:"pullRequestState,omitempty"`
}

// ResourcesTemplateStatus defines the observed state of ResourcesTemplate
type ResourcesTemplateStatus struct {
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// GetConditions returns the status conditions of the object.
func (in *ResourcesTemplate) GetConditions() []metav1.Condition {
	return in.Status.Conditions
}

// SetConditions sets the status conditions on the object.
func (in *ResourcesTemplate) SetConditions(conditions []metav1.Condition) {
	in.Status.Conditions = conditions
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ResourcesTemplate is the Schema for the resourcestemplates API
type ResourcesTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourcesTemplateSpec   `json:"spec,omitempty"`
	Status ResourcesTemplateStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ResourcesTemplateList contains a list of ResourcesTemplate
type ResourcesTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ResourcesTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ResourcesTemplate{}, &ResourcesTemplateList{})
}
