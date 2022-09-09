package v1alpha1

type GitlabProject struct {
	// GitLab project to scan. Required.
	Project string `json:"project"`

	// The GitLab API URL to talk to. If blank, uses https://gitlab.com/.
	// +optional
	API string `json:"api,omitempty"`

	// Authentication token reference.
	TokenRef SecretRef `json:"tokenRef"`
}

type GitlabMergeRequest struct {
	GitlabProject `json:",inline"`

	// The merge request id
	MergeRequestId int `json:"mergeRequestId"`
}
