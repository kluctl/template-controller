package webgit

import (
	"context"
	"encoding/json"
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

type WebgitInterface interface {
	GetProject(projectId string) (ProjectInterface, error)
}

type ProjectInterface interface {
	ListMergeRequests(state v1alpha1.MergeRequestState) ([]MergeRequestInterface, error)
	GetMergeRequest(mrId string) (MergeRequestInterface, error)
}

type MergeRequestInterface interface {
	Info() (*v1alpha1.MergeRequestInfo, error)

	HasApproved() (bool, error)
	Approve() error
	Unapprove() error

	CreateMergeRequestNote(body string) (Note, error)
	GetMergeRequestNote(noteId string) (Note, error)
	ListMergeRequestNotes() ([]Note, error)
	ListMergeRequestNotesAfter(t time.Time) ([]Note, error)
}

func BuildWebgit(ctx context.Context, client client.Client, namespace string, spec any) (ProjectInterface, error) {
	specMap, err := objectToMap(spec)
	if err != nil {
		return nil, err
	}

	if m, ok := specMap["gitlab"]; ok {
		var gitlab v1alpha1.GitlabProject
		err = mapToObject(m, &gitlab)
		if err != nil {
			return nil, err
		}

		project, err := BuildWebgitGitlab(ctx, client, namespace, gitlab)
		if err != nil {
			return nil, err
		}
		return project, nil
	} else if m, ok := specMap["github"]; ok {
		var github v1alpha1.GithubProject
		err = mapToObject(m, &github)
		if err != nil {
			return nil, err
		}

		project, err := BuildWebgitGithub(ctx, client, namespace, github)
		if err != nil {
			return nil, err
		}
		return project, nil
	} else {
		return nil, fmt.Errorf("no git project spec provided")
	}
}

func BuildWebgitMergeRequest(ctx context.Context, client client.Client, namespace string, spec any) (MergeRequestInterface, error) {
	project, err := BuildWebgit(ctx, client, namespace, spec)
	if err != nil {
		return nil, err
	}

	specMap, err := objectToMap(spec)

	if m, ok := specMap["gitlab"]; ok {
		var gitlab v1alpha1.GitlabMergeRequestRef
		err = mapToObject(m, &gitlab)
		if err != nil {
			return nil, err
		}
		if gitlab.MergeRequestId == nil {
			return nil, fmt.Errorf("missing mergeRequestId")
		}
		return project.GetMergeRequest(fmt.Sprintf("%d", *gitlab.MergeRequestId))
	} else if m, ok := specMap["github"]; ok {
		var github v1alpha1.GithubPullRequestRef
		err = mapToObject(m, &github)
		if err != nil {
			return nil, err
		}
		if github.PullRequestId == nil {
			return nil, fmt.Errorf("missing pullRequestId")
		}
		return project.GetMergeRequest(fmt.Sprintf("%d", *github.PullRequestId))
	} else {
		return nil, fmt.Errorf("no pullRequest spec provided")
	}
}

func objectToMap(obj any) (map[string]any, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var ret map[string]any
	err = json.Unmarshal(b, &ret)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		ret = map[string]any{}
	}
	return ret, nil
}

func mapToObject(m any, out any) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, out)
	if err != nil {
		return err
	}
	return nil
}
