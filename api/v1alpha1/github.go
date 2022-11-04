package v1alpha1

import (
	"encoding/json"
	"fmt"
)

type MergeRequestState string

const (
	StateAll    MergeRequestState = "all"
	StateOpened MergeRequestState = "opened"
	StateClosed MergeRequestState = "closed"
	StateMerged MergeRequestState = "merged"
)

func StateFromString(s string) (MergeRequestState, error) {
	s2 := MergeRequestState(s)
	switch s2 {
	case StateAll, StateOpened, StateClosed, StateMerged:
		break
	default:
		return "", fmt.Errorf("unsupported state %s", s2)
	}
	return s2, nil
}

func (s *MergeRequestState) MarshalJSON() ([]byte, error) {
	if s == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(string(*s))
}

func (s *MergeRequestState) UnmarshalJSON(data []byte) error {
	var s2 string
	err := json.Unmarshal(data, &s2)
	if err != nil {
		return err
	}
	*s, err = StateFromString(s2)
	return err
}

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
