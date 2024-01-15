package v1alpha1

import "k8s.io/apimachinery/pkg/util/intstr"

type GitlabProject struct {
	// Project specifies the Gitlab group and project (separated by slash) to
	// use, or the numeric project id
	// +required
	Project *intstr.IntOrString `json:"project"`

	// API specifies the GitLab API URL to talk to.
	// If blank, uses https://gitlab.com/.
	// +optional
	API *string `json:"api,omitempty"`

	// TokenRef specifies a secret and key to load the Gitlab API token from
	// +optional
	TokenRef *SecretRef `json:"tokenRef"`
}

type GitlabMergeRequestRef struct {
	GitlabProject `json:",inline"`

	// MergeRequestId specifies the Gitlab merge request internal ID
	// +required
	MergeRequestId *intstr.IntOrString `json:"mergeRequestId"`
}
