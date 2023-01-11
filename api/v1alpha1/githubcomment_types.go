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

// GithubCommentSpec defines the desired state of GithubComment
type GithubCommentSpec struct {
	GithubPullRequestRef `json:",inline"`
	CommentSpec          `json:",inline"`

	// +optional
	// +kubebuilder:default:=false
	Suspend bool `json:"suspend"`
}

// GithubCommentStatus defines the observed state of GithubComment
type GithubCommentStatus struct {
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// +optional
	CommentId string `json:"commentId,omitempty"`

	// +optional
	LastPostedBodyHash string `json:"LastPostedBodyHash,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GithubComment is the Schema for the githubcomments API
type GithubComment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GithubCommentSpec   `json:"spec,omitempty"`
	Status GithubCommentStatus `json:"status,omitempty"`
}

func (gc *GithubComment) GetCommentSourceSpec() *CommentSourceSpec {
	return &gc.Spec.Source
}

//+kubebuilder:object:root=true

// GithubCommentList contains a list of GithubComment
type GithubCommentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GithubComment `json:"items"`
}

func (l *GithubCommentList) GetItems() []client.Object {
	ret := make([]client.Object, len(l.Items))
	for i, _ := range l.Items {
		ret[i] = &l.Items[i]
	}
	return ret
}

func init() {
	SchemeBuilder.Register(&GithubComment{}, &GithubCommentList{})
}
