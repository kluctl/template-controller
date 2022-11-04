package webgit

import (
	"context"
	"fmt"
	"github.com/kluctl/template-controller/api/v1alpha1"
	"github.com/xanzy/go-gitlab"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"sync"
	"time"
)

type Gitlab struct {
	client *gitlab.Client
}

type GitlabProject struct {
	gitlab    *Gitlab
	projectId string
}

type GitlabMergeRequest struct {
	project *GitlabProject
	mrId    int

	mr    *gitlab.MergeRequest
	mutex sync.Mutex
}

func NewGitlab(baseUrl *string, token string) (WebgitInterface, error) {
	var opts []gitlab.ClientOptionFunc
	if baseUrl != nil {
		opts = append(opts, gitlab.WithBaseURL(*baseUrl))
	}
	client, err := gitlab.NewClient(token, opts...)
	if err != nil {
		return nil, err
	}
	return &Gitlab{
		client: client,
	}, nil
}

func (g *Gitlab) CurrentUserId() (string, error) {
	u, _, err := g.client.Users.CurrentUser()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", u.ID), nil
}

func (g *Gitlab) GetProject(projectId string) (ProjectInterface, error) {
	return &GitlabProject{
		gitlab:    g,
		projectId: projectId,
	}, nil
}

func (p *GitlabProject) ListMergeRequests(state v1alpha1.MergeRequestState) ([]MergeRequestInterface, error) {
	opts := &gitlab.ListProjectMergeRequestsOptions{
		State: (*string)(&state),
	}

	mrs, _, err := p.gitlab.client.MergeRequests.ListProjectMergeRequests(p.projectId, opts)
	if err != nil {
		return nil, err
	}
	var ret []MergeRequestInterface
	for _, mr := range mrs {
		ret = append(ret, &GitlabMergeRequest{
			project: p,
			mrId:    mr.ID,
			mr:      mr,
		})
	}
	return ret, nil
}

func (p *GitlabProject) GetMergeRequest(mrId string) (MergeRequestInterface, error) {
	mrIdInt, err := strconv.ParseInt(mrId, 0, 32)
	if err != nil {
		return nil, err
	}
	return &GitlabMergeRequest{
		project: p,
		mrId:    int(mrIdInt),
	}, nil
}

func (g *GitlabMergeRequest) convertNote(n *gitlab.Note) Note {
	return &GitlabNote{
		g:    g,
		note: n,
	}
}

func (g *GitlabMergeRequest) Info() (*v1alpha1.MergeRequestInfo, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.mr == nil {
		opt := &gitlab.GetMergeRequestsOptions{}
		mr, _, err := g.project.gitlab.client.MergeRequests.GetMergeRequest(g.project.projectId, g.mrId, opt)
		if err != nil {
			return nil, err
		}
		g.mr = mr
	}

	state, err := v1alpha1.StateFromString(g.mr.State)
	if err != nil {
		return nil, err
	}

	return &v1alpha1.MergeRequestInfo{
		ID:           g.mr.IID,
		TargetBranch: g.mr.TargetBranch,
		SourceBranch: g.mr.SourceBranch,
		Title:        g.mr.Title,
		State:        state,
		CreatedAt:    metav1.NewTime(*g.mr.CreatedAt),
		UpdatedAt:    metav1.NewTime(*g.mr.UpdatedAt),
		Author:       g.mr.Author.Username,
		Labels:       g.mr.Labels,
		Draft:        g.mr.Draft,
	}, nil
}

func (g *GitlabMergeRequest) HasApproved() (bool, error) {
	userId, err := g.project.gitlab.CurrentUserId()
	if err != nil {
		return false, err
	}
	mc, _, err := g.project.gitlab.client.MergeRequestApprovals.GetConfiguration(g.project.projectId, g.mrId)
	if err != nil {
		return false, err
	}
	for _, u := range mc.ApprovedBy {
		if fmt.Sprintf("%d", u.User.ID) == userId {
			return true, nil
		}
	}
	return false, nil
}

func (g *GitlabMergeRequest) Approve() error {
	opt := &gitlab.ApproveMergeRequestOptions{}
	_, _, err := g.project.gitlab.client.MergeRequestApprovals.ApproveMergeRequest(g.project.projectId, g.mrId, opt)
	return err
}

func (g *GitlabMergeRequest) Unapprove() error {
	_, err := g.project.gitlab.client.MergeRequestApprovals.UnapproveMergeRequest(g.project.projectId, g.mrId)
	return err
}

func (g *GitlabMergeRequest) CreateMergeRequestNote(body string) (Note, error) {
	opt := &gitlab.CreateMergeRequestNoteOptions{
		Body: &body,
	}
	n, _, err := g.project.gitlab.client.Notes.CreateMergeRequestNote(g.project.projectId, g.mrId, opt)
	if err != nil {
		return nil, err
	}
	return g.convertNote(n), nil
}

func (g *GitlabMergeRequest) GetMergeRequestNote(noteId string) (Note, error) {
	noteId2, err := strconv.ParseInt(noteId, 10, 32)
	if err != nil {
		return nil, err
	}
	n, _, err := g.project.gitlab.client.Notes.GetMergeRequestNote(g.project.projectId, g.mrId, int(noteId2))
	if err != nil {
		return nil, err
	}
	return g.convertNote(n), nil
}

func (g *GitlabMergeRequest) ListMergeRequestNotes() ([]Note, error) {
	sort := "desc"
	orderBy := "created_at"

	opt := &gitlab.ListMergeRequestNotesOptions{}
	opt.Page = 1
	opt.PerPage = 100
	opt.Sort = &sort
	opt.OrderBy = &orderBy

	var ret []Note
	for true {
		notes, _, err := g.project.gitlab.client.Notes.ListMergeRequestNotes(g.project.projectId, g.mrId, opt)
		if err != nil {
			return nil, err
		}

		for _, n := range notes {
			ret = append(ret, g.convertNote(n))
		}

		if len(notes) < opt.PerPage {
			break
		}
		opt.Page++
	}
	return ret, nil
}

func (g *GitlabMergeRequest) ListMergeRequestNotesAfter(t time.Time) ([]Note, error) {
	sort := "desc"
	orderBy := "created_at"

	opt := &gitlab.ListMergeRequestNotesOptions{}
	opt.Page = 1
	opt.PerPage = 10
	opt.Sort = &sort
	opt.OrderBy = &orderBy

	if t == (time.Time{}) {
		opt.PerPage = 100
	}

	var ret []Note

outer:
	for true {
		notes, _, err := g.project.gitlab.client.Notes.ListMergeRequestNotes(g.project.projectId, g.mrId, opt)
		if err != nil {
			return nil, err
		}

		for _, n := range notes {
			if !n.CreatedAt.After(t) {
				break outer
			}
			ret = append(ret, g.convertNote(n))
		}

		if len(notes) < opt.PerPage {
			break
		}
		opt.Page++
	}
	// reverse ret
	for i, j := 0, len(ret)-1; i < j; i, j = i+1, j-1 {
		ret[i], ret[j] = ret[j], ret[i]
	}
	return ret, nil
}

type GitlabNote struct {
	g    *GitlabMergeRequest
	note *gitlab.Note
}

func (n *GitlabNote) GetId() string {
	return fmt.Sprintf("%d", n.note.ID)
}

func (n *GitlabNote) GetBody() string {
	return n.note.Body
}

func (n *GitlabNote) GetCreatedAt() time.Time {
	return *n.note.CreatedAt
}

func (n *GitlabNote) UpdateBody(body string) error {
	opt := &gitlab.UpdateMergeRequestNoteOptions{
		Body: &body,
	}
	n2, _, err := n.g.project.gitlab.client.Notes.UpdateMergeRequestNote(n.g.project.projectId, n.g.mrId, n.note.ID, opt)
	if err != nil {
		return err
	}
	n.note = n2
	return nil
}

func BuildWebgitGitlab(ctx context.Context, client client.Client, namespace string, info v1alpha1.GitlabProject) (ProjectInterface, error) {
	if info.Project == "" {
		return nil, fmt.Errorf("missing gitlab project")
	}
	if info.TokenRef == nil {
		return nil, fmt.Errorf("missing tokenRef")
	}

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
