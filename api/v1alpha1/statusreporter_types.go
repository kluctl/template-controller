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
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StatusReporterSpec defines the desired state of StatusReporter
type StatusReporterSpec struct {
	// +kubebuilder:default:="1m"
	Interval metav1.Duration `json:"interval"`

	// +required
	ForObject ObjectRef `json:"forObject"`

	// +required
	Reporters []Reporter `json:"reporters"`
}

type Reporter struct {
	// +optional
	PullRequestComment *PullRequestCommentReporter `json:"pullRequestComment,omitempty"`
	// +optional
	PullRequestApprove *PullRequestApproveReporter `json:"pullRequestApprove,omitempty"`
}

func (r *Reporter) BuildKey() string {
	b, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	s := sha256.Sum256(b)
	return hex.EncodeToString(s[:])
}

type ReporterStatus struct {
	Key string `json:"key"`

	// +optional
	Error string `json:"error,omitempty"`

	// +optional
	PullRequestComment *PullRequestCommentReporterStatus `json:"pullRequestComment,omitempty"`
	// +optional
	PullRequestApprove *PullRequestApproveReporterStatus `json:"pullRequestApprove,omitempty"`
}

type PullRequestCommentReporter struct {
	// +optional
	Gitlab *GitlabMergeRequest `json:"gitlab,omitempty"`
}

type PullRequestCommentReporterStatus struct {
	// +optional
	LastPostedStatusHash string `json:"lastPostedStatusHash,omitempty"`

	// +optional
	NoteId string `json:"noteId,omitempty"`
}

type PullRequestApproveReporter struct {
	// +optional
	Gitlab *GitlabMergeRequest `json:"gitlab,omitempty"`
}

type PullRequestApproveReporterStatus struct {
	// +optional
	Approved *bool `json:"approved,omitempty"`
}

// StatusReporterStatus defines the observed state of StatusReporter
type StatusReporterStatus struct {
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// +optional
	ReporterStatus []*ReporterStatus `json:"reporterStatus"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// StatusReporter is the Schema for the statusreporters API
type StatusReporter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StatusReporterSpec   `json:"spec,omitempty"`
	Status StatusReporterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// StatusReporterList contains a list of StatusReporter
type StatusReporterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StatusReporter `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StatusReporter{}, &StatusReporterList{})
}
