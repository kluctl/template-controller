package webgit

import (
	"context"
	"fmt"
	"github.com/google/go-github/v47/github"
	"github.com/kluctl/template-controller/api/v1alpha1"
	"golang.org/x/oauth2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"sync"
	"time"
)

type GithubMergeRequest struct {
	ctx context.Context

	client *github.Client
	owner  string
	repo   string
	prId   int

	currentUserCache *github.User
	currentUserMutex sync.Mutex

	pr     *github.PullRequest
	review *github.PullRequestReview
	mutex  sync.Mutex
}

func (g *GithubMergeRequest) getCurrentUser() (*github.User, error) {
	g.currentUserMutex.Lock()
	defer g.currentUserMutex.Unlock()

	if g.currentUserCache != nil {
		return g.currentUserCache, nil
	}

	req, err := g.client.NewRequest("GET", "/user", nil)
	if err != nil {
		return nil, err
	}

	var user github.User
	_, err = g.client.Do(g.ctx, req, &user)
	if err != nil {
		return nil, err
	}

	g.currentUserCache = &user

	return g.currentUserCache, nil
}

func (g *GithubMergeRequest) convertComment(n *github.IssueComment) Note {
	return &GithubNote{
		g:       g,
		comment: n,
	}
}

func (g *GithubMergeRequest) findReview() (*github.PullRequestReview, error) {
	currentUser, err := g.getCurrentUser()
	if err != nil {
		return nil, err
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.review != nil {
		return g.review, nil
	}

	opts := &github.ListOptions{}
	opts.Page = 0
	opts.PerPage = 100

	for {
		reviews, _, err := g.client.PullRequests.ListReviews(g.ctx, g.owner, g.repo, g.prId, opts)
		if err != nil {
			return nil, err
		}
		for _, r := range reviews {
			if r.User.ID == currentUser.ID {
				g.review = r
				return g.review, nil
			}
		}
		if len(reviews) < opts.PerPage {
			break
		}
		opts.Page++
	}

	return nil, nil
}

func (g *GithubMergeRequest) HasApproved() (bool, error) {
	review, err := g.findReview()
	if err != nil {
		return false, err
	}
	if review == nil {
		return false, nil
	}
	if *review.State == "APPROVED" {
		return true, nil
	}
	return false, nil
}

func (g *GithubMergeRequest) Approve() error {
	event := "APPROVE"
	body := "Approved"

	req := &github.PullRequestReviewRequest{
		Body:  &body,
		Event: &event,
	}
	_, _, err := g.client.PullRequests.CreateReview(g.ctx, g.owner, g.repo, g.prId, req)
	if err != nil {
		return err
	}

	return nil
}

func (g *GithubMergeRequest) Unapprove() error {
	review, err := g.findReview()
	if err != nil {
		return err
	}
	if review == nil {
		return nil
	}

	message := "Not approved"
	req := &github.PullRequestReviewDismissalRequest{
		Message: &message,
	}
	_, _, err = g.client.PullRequests.DismissReview(g.ctx, g.owner, g.repo, g.prId, *review.ID, req)
	if err != nil {
		return err
	}
	return nil
}

func (g *GithubMergeRequest) CreateMergeRequestNote(body string) (Note, error) {
	comment := &github.IssueComment{
		Body: &body,
	}
	n, _, err := g.client.Issues.CreateComment(g.ctx, g.owner, g.repo, g.prId, comment)
	if err != nil {
		return nil, err
	}
	return g.convertComment(n), nil
}

func (g *GithubMergeRequest) GetMergeRequestNote(noteId string) (Note, error) {
	noteId2, err := strconv.ParseInt(noteId, 10, 32)
	if err != nil {
		return nil, err
	}
	n, resp, err := g.client.Issues.GetComment(g.ctx, g.owner, g.repo, noteId2)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}
	return g.convertComment(n), nil
}

func (g *GithubMergeRequest) ListMergeRequestNotes() ([]Note, error) {
	sort := "created"
	direction := "desc"

	opt := &github.IssueListCommentsOptions{}
	opt.Page = 1
	opt.PerPage = 100
	opt.Sort = &sort
	opt.Direction = &direction

	var ret []Note
	for true {
		notes, _, err := g.client.Issues.ListComments(g.ctx, g.owner, g.repo, g.prId, opt)
		if err != nil {
			return nil, err
		}

		for _, n := range notes {
			ret = append(ret, g.convertComment(n))
		}

		if len(notes) < opt.PerPage {
			break
		}
		opt.Page++
	}
	return ret, nil
}

func (g *GithubMergeRequest) ListMergeRequestNotesAfter(t time.Time) ([]Note, error) {
	opt := &github.IssueListCommentsOptions{}
	opt.Page = 1
	opt.PerPage = 100

	if t != (time.Time{}) {
		opt.Since = &t
	}

	var ret []Note
	for true {
		notes, _, err := g.client.Issues.ListComments(g.ctx, g.owner, g.repo, g.prId, opt)
		if err != nil {
			return nil, err
		}

		for _, n := range notes {
			ret = append(ret, g.convertComment(n))
		}

		if len(notes) < opt.PerPage {
			break
		}
		opt.Page++
	}
	return ret, nil
}

type GithubNote struct {
	g       *GithubMergeRequest
	comment *github.IssueComment
}

func (n *GithubNote) GetId() string {
	return fmt.Sprintf("%d", *n.comment.ID)
}

func (n *GithubNote) GetBody() string {
	return *n.comment.Body
}

func (n *GithubNote) GetCreatedAt() time.Time {
	return *n.comment.CreatedAt
}

func (n *GithubNote) UpdateBody(body string) error {
	updateComment := github.IssueComment{
		Body: &body,
	}

	newComment, _, err := n.g.client.Issues.EditComment(n.g.ctx, n.g.owner, n.g.repo, *n.comment.ID, &updateComment)
	if err != nil {
		return err
	}
	n.comment = newComment
	return nil
}

func BuildWebgitMergeRequestGithub(ctx context.Context, client client.Client, namespace string, info v1alpha1.GithubPullRequestRef) (*GithubMergeRequest, error) {
	if info.Owner == "" {
		return nil, fmt.Errorf("missing github owner")
	}
	if info.Repo == "" {
		return nil, fmt.Errorf("missing github owner")
	}
	if info.TokenRef == nil {
		return nil, fmt.Errorf("missing github tokenRef")
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
		return nil, fmt.Errorf("github token is missing in secret")
	}
	token := string(tokenBytes)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &GithubMergeRequest{
		ctx:    ctx,
		client: github.NewClient(tc),
		owner:  info.Owner,
		repo:   info.Repo,
		prId:   info.PullRequestId,
	}, nil
}
