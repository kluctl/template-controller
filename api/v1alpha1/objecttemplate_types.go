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

// ObjectTemplateSpec defines the desired state of ObjectTemplate
type ObjectTemplateSpec struct {
	// +kubebuilder:default:="30s"
	Interval metav1.Duration `json:"interval"`

	// +optional
	Defaults *ObjectTemplateDefaultsSpec `json:"defaults,omitempty"`

	// +required
	Generators []Generator `json:"generators"`

	// +required
	Templates []Template `json:"templates"`
}

type ObjectTemplateDefaultsSpec struct {
	// +optional
	Gitlab *GitlabProject `json:"gitlab,omitempty"`

	// +optional
	Github *GithubProject `json:"github,omitempty"`
}

type Template struct {
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Object *unstructured.Unstructured `json:"object,omitempty"`

	// +optional
	Raw *string `json:"raw,omitempty"`
}

type Generator struct {
	// +optional
	PullRequest *PullRequestGenerator `json:"pullRequest,omitempty"`
}

type PullRequestGenerator struct {
	// +optional
	Gitlab *PullRequestGeneratorGitlab `json:"gitlab,omitempty"`

	// +optional
	Github *PullRequestGeneratorGithub `json:"github,omitempty"`
}

type PullRequestGeneratorGitlab struct {
	GitlabProject `json:",inline"`

	// +optional
	TargetBranch *string `json:"targetBranch,omitempty"`

	// +optional
	SourceBranch *string `json:"sourceBranch,omitempty"`

	// Labels is used to filter the MRs that you want to target
	// +optional
	Labels []string `json:"labels,omitempty"`

	// PullRequestState is an additional MRs filter to get only those with a certain state. Default: "all"
	// +optional
	// +kubebuilder:validation:Enum=all;opened;closed;merged
	// +kubebuilder:default:="all"
	PullRequestState MergeRequestState `json:"pullRequestState,omitempty"`
}

type PullRequestGeneratorGithub struct {
	GithubProject `json:",inline"`

	// Labels is used to filter the MRs that you want to target
	Labels []string `json:"labels,omitempty"`
	// PullRequestState is an additional MRs filter to get only those with a certain state. Default: "all"
	PullRequestState MergeRequestState `json:"pullRequestState,omitempty"`
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
