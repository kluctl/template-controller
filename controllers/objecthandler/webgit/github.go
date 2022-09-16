package webgit

import (
	"context"
	"fmt"
	"github.com/google/go-github/v47/github"
	"github.com/kluctl/template-controller/api/v1alpha1"
	"golang.org/x/oauth2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Github struct {
	ctx    context.Context
	client *github.Client

	currentUser *github.User
	mutex       sync.Mutex
}

type GithubProject struct {
	github *Github
	owner  string
	repo   string
}

type GithubMergeRequest struct {
	project *GithubProject
	mrId    int

	pr     *github.PullRequest
	review *github.PullRequestReview
	mutex  sync.Mutex
}

func NewGithub(ctx context.Context, token string) (WebgitInterface, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	g := &Github{
		ctx:    ctx,
		client: client,
	}

	return g, nil
}

func (g *Github) getCurrentUser() (*github.User, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.currentUser != nil {
		return g.currentUser, nil
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

	g.currentUser = &user

	return g.currentUser, nil
}

func (g *Github) GetProject(projectId string) (ProjectInterface, error) {
	p := strings.SplitN(projectId, "/", 2)
	if len(p) < 2 {
		return nil, fmt.Errorf("invalid project id %s", projectId)
	}
	return &GithubProject{
		github: g,
		owner:  p[0],
		repo:   p[1],
	}, nil
}

func (g *GithubProject) ListMergeRequests(state MergeRequestState) ([]MergeRequestInterface, error) {
	opts := &github.PullRequestListOptions{
		State: string(state),
	}
	prs, _, err := g.github.client.PullRequests.List(g.github.ctx, g.owner, g.repo, opts)
	if err != nil {
		return nil, err
	}
	var ret []MergeRequestInterface
	for _, pr := range prs {
		ret = append(ret, &GithubMergeRequest{
			project: g,
			mrId:    *pr.Number,
			pr:      pr,
		})
	}
	return ret, nil
}

func (g *GithubProject) GetMergeRequest(mrId string) (MergeRequestInterface, error) {
	mrIdInt, err := strconv.ParseInt(mrId, 0, 32)
	if err != nil {
		return nil, err
	}
	return &GithubMergeRequest{
		project: g,
		mrId:    int(mrIdInt),
	}, nil
}

func (g *GithubMergeRequest) convertComment(n *github.IssueComment) Note {
	return &GithubNote{
		g:       g,
		comment: n,
	}
}

func (g *GithubMergeRequest) convertStateToGithub(state MergeRequestState) (string, error) {
	switch state {
	case StateAll:
		return "all", nil
	case StateOpened:
		return "open", nil
	case StateClosed:
		return "closed", nil
	case StateMerged:
		return "closed", nil
	}
	return "", fmt.Errorf("invalid state %s", state)
}

func (g *GithubMergeRequest) convertStateFromGithub(state string) (MergeRequestState, error) {
	switch state {
	case "all":
		return StateAll, nil
	case "open":
		return StateOpened, nil
	case "closed":
		return StateClosed, nil
	}
	return "", fmt.Errorf("invalid state %s", state)
}

func (g *GithubMergeRequest) Info() (*MergeRequestInfo, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.pr == nil {
		pr, _, err := g.project.github.client.PullRequests.Get(g.project.github.ctx, g.project.owner, g.project.repo, int(g.mrId))
		if err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}
		g.pr = pr
	}

	state, err := g.convertStateFromGithub(*g.pr.State)
	if err != nil {
		return nil, err
	}

	labels := make([]string, 0, len(g.pr.Labels))
	for _, l := range g.pr.Labels {
		labels = append(labels, *l.Name)
	}

	return &MergeRequestInfo{
		ID:           *g.pr.Number,
		TargetBranch: *g.pr.Base.Ref,
		SourceBranch: *g.pr.Head.Ref,
		Title:        *g.pr.Title,
		State:        state,
		CreatedAt:    metav1.NewTime(*g.pr.CreatedAt),
		UpdatedAt:    metav1.NewTime(*g.pr.UpdatedAt),
		Author:       *g.pr.User.Login,
		Labels:       labels,
		Draft:        *g.pr.Draft,
	}, nil
}

func (g *GithubMergeRequest) findReview() (*github.PullRequestReview, error) {
	currentUser, err := g.project.github.getCurrentUser()
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
		reviews, _, err := g.project.github.client.PullRequests.ListReviews(g.project.github.ctx, g.project.owner, g.project.repo, g.mrId, opts)
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
	_, _, err := g.project.github.client.PullRequests.CreateReview(g.project.github.ctx, g.project.owner, g.project.repo, g.mrId, req)
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
	_, _, err = g.project.github.client.PullRequests.DismissReview(g.project.github.ctx, g.project.owner, g.project.repo, g.mrId, *review.ID, req)
	if err != nil {
		return err
	}
	return nil
}

func (g *GithubMergeRequest) CreateMergeRequestNote(body string) (Note, error) {
	comment := &github.IssueComment{
		Body: &body,
	}
	n, _, err := g.project.github.client.Issues.CreateComment(g.project.github.ctx, g.project.owner, g.project.repo, g.mrId, comment)
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
	n, _, err := g.project.github.client.Issues.GetComment(g.project.github.ctx, g.project.owner, g.project.repo, noteId2)
	if err != nil {
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
		notes, _, err := g.project.github.client.Issues.ListComments(g.project.github.ctx, g.project.owner, g.project.repo, g.mrId, opt)
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
	sort := "created"
	direction := "desc"

	opt := &github.IssueListCommentsOptions{}
	opt.Page = 1
	opt.PerPage = 10
	opt.Sort = &sort
	opt.Direction = &direction

	if t == (time.Time{}) {
		opt.PerPage = 100
	}

	var ret []Note

outer:
	for true {
		notes, _, err := g.project.github.client.Issues.ListComments(g.project.github.ctx, g.project.owner, g.project.repo, g.mrId, opt)
		if err != nil {
			return nil, err
		}

		for _, n := range notes {
			if !n.CreatedAt.After(t) {
				break outer
			}
			ret = append(ret, g.convertComment(n))
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

	newComment, _, err := n.g.project.github.client.Issues.EditComment(n.g.project.github.ctx, n.g.project.owner, n.g.project.repo, *n.comment.ID, &updateComment)
	if err != nil {
		return err
	}
	n.comment = newComment
	return nil
}

func BuildWebgitGithub(ctx context.Context, client client.Client, namespace string, info v1alpha1.GithubProject) (ProjectInterface, error) {
	if info.Owner == nil {
		return nil, fmt.Errorf("missing github owner")
	}
	if info.Repo == nil {
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
		return nil, fmt.Errorf("gitlab token is missing in secret")
	}
	token := string(tokenBytes)

	g, err := NewGithub(ctx, token)
	if err != nil {
		return nil, err
	}
	return g.GetProject(fmt.Sprintf("%s/%s", info.Owner, info.Repo))
}
