package v1alpha1

type CommentSpec struct {
	// +optional
	Id *string `json:"id,omitempty"`

	// +required
	Source CommentSourceSpec `json:"source"`
}

type CommentSourceSpec struct {
	// +optional
	Text *string `json:"text,omitempty"`

	// +optional
	ConfigMap *ConfigMapRef `json:"configMap,omitempty"`

	// +optional
	TextTemplate *LocalObjectReference `json:"textTemplate,omitempty"`
}
