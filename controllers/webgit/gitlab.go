package webgit

import (
	"context"
	"fmt"
	"github.com/kluctl/template-controller/api/v1alpha1"
	"github.com/xanzy/go-gitlab"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"sync"
	"time"
)

type GitlabMergeRequest struct {
	client *gitlab.Client

	projectId string
	mrId      int

	currentUserCache *gitlab.User
	currentUserMutex sync.Mutex

	mr    *gitlab.MergeRequest
	mutex sync.Mutex
}

func (g *GitlabMergeRequest) convertNote(n *gitlab.Note) Note {
	return &GitlabNote{
		g:    g,
		note: n,
	}
}

func (g *GitlabMergeRequest) currentUser() (*gitlab.User, error) {
	g.currentUserMutex.Lock()
	defer g.currentUserMutex.Unlock()
	if g.currentUserCache != nil {
		return g.currentUserCache, nil
	}
	var err error
	g.currentUserCache, _, err = g.client.Users.CurrentUser()
	if err != nil {
		return nil, err
	}
	return g.currentUserCache, nil
}

func (g *GitlabMergeRequest) HasApproved() (bool, error) {
	cu, err := g.currentUser()
	if err != nil {
		return false, err
	}

	mc, _, err := g.client.MergeRequestApprovals.GetConfiguration(g.projectId, g.mrId)
	if err != nil {
		return false, err
	}
	for _, u := range mc.ApprovedBy {
		if u.User.ID == cu.ID {
			return true, nil
		}
	}
	return false, nil
}

func (g *GitlabMergeRequest) Approve() error {
	opt := &gitlab.ApproveMergeRequestOptions{}
	_, _, err := g.client.MergeRequestApprovals.ApproveMergeRequest(g.projectId, g.mrId, opt)
	return err
}

func (g *GitlabMergeRequest) Unapprove() error {
	_, err := g.client.MergeRequestApprovals.UnapproveMergeRequest(g.projectId, g.mrId)
	return err
}

func (g *GitlabMergeRequest) CreateMergeRequestNote(body string) (Note, error) {
	opt := &gitlab.CreateMergeRequestNoteOptions{
		Body: &body,
	}
	n, _, err := g.client.Notes.CreateMergeRequestNote(g.projectId, g.mrId, opt)
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
	n, resp, err := g.client.Notes.GetMergeRequestNote(g.projectId, g.mrId, int(noteId2))
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, nil
		}
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
		notes, _, err := g.client.Notes.ListMergeRequestNotes(g.projectId, g.mrId, opt)
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
		notes, _, err := g.client.Notes.ListMergeRequestNotes(g.projectId, g.mrId, opt)
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
	n2, _, err := n.g.client.Notes.UpdateMergeRequestNote(n.g.projectId, n.g.mrId, n.note.ID, opt)
	if err != nil {
		return err
	}
	n.note = n2
	return nil
}

func BuildWebgitMergeRequestGitlab(ctx context.Context, client client.Client, namespace string, info v1alpha1.GitlabMergeRequestRef) (*GitlabMergeRequest, error) {
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

	var opts []gitlab.ClientOptionFunc
	if info.API != nil {
		opts = append(opts, gitlab.WithBaseURL(*info.API))
	}
	glClient, err := gitlab.NewClient(token, opts...)
	if err != nil {
		return nil, err
	}

	return &GitlabMergeRequest{
		client:    glClient,
		projectId: info.Project,
		mrId:      info.MergeRequestId,
	}, nil
}
