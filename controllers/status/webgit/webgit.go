package webgit

import (
	"context"
	"fmt"
	"github.com/kluctl/template-controller/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type Note interface {
	GetId() string
	GetBody() string

	UpdateBody(body string) error
	GetCreatedAt() time.Time
}

type WebgitInterface interface {
	GetProject(projectId string) (ProjectInterface, error)
}

type ProjectInterface interface {
	ListMergeRequests(targetBranch *string, sourceBranch *string) ([]MergeRequestInterface, error)
	GetMergeRequest(mrId string) (MergeRequestInterface, error)
}

type MergeRequestInterface interface {
	HasApproved() (bool, error)
	Approve() error
	Unapprove() error

	CreateMergeRequestNote(body string) (Note, error)
	GetMergeRequestNote(noteId string) (Note, error)
	ListMergeRequestNotes() ([]Note, error)
	ListMergeRequestNotesAfter(t time.Time) ([]Note, error)
}

func BuildWebgitGitlab(ctx context.Context, client client.Client, namespace string, info v1alpha1.GitlabProject) (ProjectInterface, error) {
	sn := types.NamespacedName{
		Namespace: namespace,
		Name:      info.TokenRef.SecretName,
	}

	var secret v1.Secret
	err := client.Get(ctx, sn, &secret)
	if err != nil {
		return nil, err
	}

	tokenBytes, ok := secret.Data[info.TokenRef.Key]
	if !ok {
		return nil, fmt.Errorf("gitlab token is missing in secret")
	}
	token := string(tokenBytes)

	g, err := NewGitlab(info.API, token)
	if err != nil {
		return nil, err
	}
	return g.GetProject(info.Project)
}
