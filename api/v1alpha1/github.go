package v1alpha1

import "k8s.io/apimachinery/pkg/util/intstr"

type GithubProject struct {
	// Owner specifies the GitHub user or organisation that owns the repository
	// +required
	Owner string `json:"owner"`

	// Repo specifies the repository name.
	// +required
	Repo string `json:"repo"`

	// TokenRef specifies a secret and key to load the GitHub API token from
	// +optional
	TokenRef *SecretRef `json:"tokenRef"`
}

type GithubPullRequestRef struct {
	GithubProject `json:",inline"`

	// PullRequestId specifies the pull request ID.
	// +required
	PullRequestId *intstr.IntOrString `json:"pullRequestId"`
}
