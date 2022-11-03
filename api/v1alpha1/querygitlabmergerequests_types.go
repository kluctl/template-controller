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

// QueryGitlabMergeRequestsSpec defines the desired state of QueryGitlabMergeRequests
type QueryGitlabMergeRequestsSpec struct {
	// Interval is the interval at which to query the Gitlab API.
	// Defaults to 5m.
	// +optional
	// +kubebuilder:default:="5m"
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ms|s|m|h))+$"
	Interval metav1.Duration `json:"interval"`

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
	// +kubebuilder:validation:Enum=all;opened;closed;locked;merged
	// +kubebuilder:default:="all"
	State *string `json:"state,omitempty"`

	// Limit limits the maximum number of pull requests to fetch. Defaults to 100
	// +kubebuilder:default:=100
	Limit int `json:"limit"`
}

// QueryGitlabMergeRequestsStatus defines the observed state of QueryGitlabMergeRequests
type QueryGitlabMergeRequestsStatus struct {
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	MergeRequests []runtime.RawExtension `json:"mergeRequests,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// QueryGitlabMergeRequests is the Schema for the querygitlabmergerequests API
type QueryGitlabMergeRequests struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QueryGitlabMergeRequestsSpec   `json:"spec,omitempty"`
	Status QueryGitlabMergeRequestsStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// QueryGitlabMergeRequestsList contains a list of QueryGitlabMergeRequests
type QueryGitlabMergeRequestsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QueryGitlabMergeRequests `json:"items"`
}

func init() {
	SchemeBuilder.Register(&QueryGitlabMergeRequests{}, &QueryGitlabMergeRequestsList{})
}
