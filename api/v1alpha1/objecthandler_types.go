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

// ObjectHandlerSpec defines the desired state of ObjectHandler
type ObjectHandlerSpec struct {
	// +kubebuilder:default:="1m"
	Interval metav1.Duration `json:"interval"`

	// +required
	ForObject ObjectRef `json:"forObject"`

	// +optional
	Defaults *ObjectHandlerDefaultsSpec `json:"defaults,omitempty"`

	// +required
	Handlers []Handler `json:"handlers"`
}

type ObjectHandlerDefaultsSpec struct {
	// +optional
	Gitlab *GitlabMergeRequest `json:"gitlab,omitempty"`
}

type Handler struct {
	// +optional
	PullRequestComment *PullRequestCommentReporter `json:"pullRequestComment,omitempty"`
	// +optional
	PullRequestApprove *PullRequestApproveReporter `json:"pullRequestApprove,omitempty"`
	// +optional
	PullRequestCommand *PullRequestCommandHandler `json:"pullRequestCommand,omitempty"`
}

func (r *Handler) BuildKey() string {
	b, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	s := sha256.Sum256(b)
	return hex.EncodeToString(s[:])
}

type HandlerStatus struct {
	Key string `json:"key"`

	// +optional
	Error string `json:"error,omitempty"`

	// +optional
	PullRequestComment *PullRequestCommentReporterStatus `json:"pullRequestComment,omitempty"`
	// +optional
	PullRequestApprove *PullRequestApproveReporterStatus `json:"pullRequestApprove,omitempty"`
	// +optional
	PullRequestCommand *PullRequestCommandHandlerStatus `json:"pullRequestCommand,omitempty"`
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

	// +optional
	// +kubebuilder:default:=false
	MissingReadyConditionIsError bool `json:"missingReadyConditionIsError"`
}

type PullRequestApproveReporterStatus struct {
	// +optional
	Approved *bool `json:"approved,omitempty"`
}

type PullRequestCommandHandler struct {
	// +optional
	Gitlab *GitlabMergeRequest `json:"gitlab,omitempty"`

	// +optional
	PostHelpComment bool `json:"postHelpComment"`

	// +required
	Commands []PullRequestCommandHandlerCommandSpec `json:"commands"`
}

type PullRequestCommandHandlerCommandSpec struct {
	// +required
	Name string `json:"name"`
	// +optional
	Description string `json:"description,omitempty"`
	// +required
	Actions []PullRequestCommandHandlerActionSpec `json:"actions"`
}

type PullRequestCommandHandlerActionSpec struct {
	// +optional
	Annotate *PullRequestCommandHandlerActionAnnotateSpec `json:"annotate"`
}

type PullRequestCommandHandlerActionAnnotateSpec struct {
	// +required
	Annotation string `json:"annotation"`
	// +required
	Value string `json:"value"`
}

type PullRequestCommandHandlerStatus struct {
	// +optional
	LastProcessedCommentTime *string `json:"lastProcessedCommentTime"`

	// +optional
	HelpNoteId string `json:"helpNoteId,omitempty"`

	// +optional
	HelpNoteBodyHash string `json:"helpNoteBodyHash,omitempty"`
}

// ObjectHandlerStatus defines the observed state of ObjectHandler
type ObjectHandlerStatus struct {
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// +optional
	HandlerStatus []*HandlerStatus `json:"handlerStatus"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ObjectHandler is the Schema for the objecthandlers API
type ObjectHandler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ObjectHandlerSpec   `json:"spec,omitempty"`
	Status ObjectHandlerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ObjectHandlerList contains a list of ObjectHandler
type ObjectHandlerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ObjectHandler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ObjectHandler{}, &ObjectHandlerList{})
}
