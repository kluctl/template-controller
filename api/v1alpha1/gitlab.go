package v1alpha1

type GitlabProject struct {
	// GitLab project to scan. Required.
	// +required
	Project string `json:"project"`

	// The GitLab API URL to talk to. If blank, uses https://gitlab.com/.
	// +optional
	API *string `json:"api,omitempty"`

	// Authentication token reference.
	// +optional
	TokenRef *SecretRef `json:"tokenRef"`
}

type GitlabMergeRequestRef struct {
	GitlabProject `json:",inline"`

	// The merge request id
	// +required
	MergeRequestId int `json:"mergeRequestId,omitempty"`
}
