package v1alpha1

type GithubProject struct {
	// +required
	Owner string `json:"owner"`

	// +required
	Repo string `json:"repo"`

	// Authentication token reference.
	// +optional
	TokenRef *SecretRef `json:"tokenRef"`
}

type GithubPullRequestRef struct {
	GithubProject `json:",inline"`

	// +required
	PullRequestId *int `json:"pullRequestId,omitempty"`
}
