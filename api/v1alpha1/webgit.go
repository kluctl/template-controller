package v1alpha1

type GitlabProject struct {
	// GitLab project to scan. Required.
	// +optional
	Project *string `json:"project"`

	// The GitLab API URL to talk to. If blank, uses https://gitlab.com/.
	// +optional
	API *string `json:"api,omitempty"`

	// Authentication token reference.
	// +optional
	TokenRef *SecretRef `json:"tokenRef"`
}

type GitlabMergeRequest struct {
	GitlabProject `json:",inline"`

	// The merge request id
	// +optional
	MergeRequestId *int `json:"mergeRequestId,omitempty"`
}

type GithubProject struct {
	// +optional
	Owner *string `json:"owner"`

	// +optional
	Repo *string `json:"repo"`

	// Authentication token reference.
	// +optional
	TokenRef *SecretRef `json:"tokenRef"`
}

type GithubPullRequest struct {
	GithubProject `json:",inline"`

	// +optional
	PullRequestId *int `json:"pullRequestId,omitempty"`
}
