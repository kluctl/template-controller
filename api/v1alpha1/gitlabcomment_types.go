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
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GitlabCommentSpec defines the desired state of GitlabComment
type GitlabCommentSpec struct {
	GitlabMergeRequestRef `json:"gitlab"`
	CommentSpec           `json:"comment"`

	// +optional
	// +kubebuilder:default:=false
	Suspend bool `json:"suspend"`
}

// GitlabCommentStatus defines the observed state of GitlabComment
type GitlabCommentStatus struct {
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// +optional
	NoteId string `json:"noteId,omitempty"`

	// +optional
	LastPostedBodyHash string `json:"lastPostedBodyHash,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GitlabComment is the Schema for the gitlabcomments API
type GitlabComment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GitlabCommentSpec   `json:"spec,omitempty"`
	Status GitlabCommentStatus `json:"status,omitempty"`
}

func (gc *GitlabComment) GetCommentSourceSpec() *CommentSourceSpec {
	return &gc.Spec.Source
}

//+kubebuilder:object:root=true

// GitlabCommentList contains a list of GitlabComment
type GitlabCommentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GitlabComment `json:"items"`
}

func (l *GitlabCommentList) GetItems() []client.Object {
	ret := make([]client.Object, len(l.Items))
	for i, _ := range l.Items {
		ret[i] = &l.Items[i]
	}
	return ret
}

func init() {
	SchemeBuilder.Register(&GitlabComment{}, &GitlabCommentList{})
}
