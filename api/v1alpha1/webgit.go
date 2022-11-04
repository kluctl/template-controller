package v1alpha1

import (
	"encoding/json"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

type MergeRequestInfo struct {
	ID           int               `json:"id"`
	TargetBranch string            `json:"targetBranch"`
	SourceBranch string            `json:"sourceBranch"`
	Title        string            `json:"title"`
	State        MergeRequestState `json:"state"`
	CreatedAt    metav1.Time       `json:"createdAt"`
	UpdatedAt    metav1.Time       `json:"updatedAt"`
	Author       string            `json:"author"`
	Labels       []string          `json:"labels"`
	Draft        bool              `json:"draft"`
}
