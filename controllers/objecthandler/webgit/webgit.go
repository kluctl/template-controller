package webgit

import (
	"context"
	"fmt"
	"github.com/kluctl/template-controller/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type Note interface {
	GetId() string
	GetBody() string

	UpdateBody(body string) error
	GetCreatedAt() time.Time
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

func BuildWebgitMergeRequest(ctx context.Context, client client.Client, namespace string, holder v1alpha1.PullRequestRefHolder) (MergeRequestInterface, error) {
	if holder.Gitlab != nil {
		return BuildWebgitMergeRequestGitlab(ctx, client, namespace, *holder.Gitlab)
	} else if holder.Github != nil {
		return BuildWebgitMergeRequestGithub(ctx, client, namespace, *holder.Github)
	} else {
		return nil, fmt.Errorf("no git merge request spec provided")
	}
}
