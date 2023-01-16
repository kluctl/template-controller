package v1alpha1

type CommentSpec struct {
	// Id specifies the identifier to be used by the controller when it needs to find the actual comment when it does
	// not know the internal id. This Id is written into the comment inside a comment, so that a simple text search
	// can reveal the comment
	// +optional
	Id *string `json:"id,omitempty"`

	// Source specifies the source content for the comment. Different source types are supported:
	// Text, ConfigMap and TextTemplate
	// +required
	Source CommentSourceSpec `json:"source"`
}

type CommentSourceSpec struct {
	// Text specifies a raw text comment.
	// +optional
	Text *string `json:"text,omitempty"`

	// ConfigMap specifies a ConfigMap and a key to load the source content from
	// +optional
	ConfigMap *ConfigMapRef `json:"configMap,omitempty"`

	// TextTemplate specifies a TextTemplate to load the source content from
	// +optional
	TextTemplate *LocalObjectReference `json:"textTemplate,omitempty"`
}
