package generators

import (
	"context"
	"fmt"
	"github.com/kluctl/kluctl/v2/pkg/utils/uo"
	"github.com/xanzy/go-gitlab"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	templatesv1alpha1 "kluctl/template-controller/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PullRequestGitlabGenerator struct {
	client  *gitlab.Client
	project string
	labels  []string
	state   string
	filters []templatesv1alpha1.PullRequestGeneratorFilter
}

func (g *PullRequestGitlabGenerator) convertMR(mr *gitlab.MergeRequest) MergeRequestInfo {
	return MergeRequestInfo{
		ID:           mr.ID,
		IID:          mr.IID,
		TargetBranch: mr.TargetBranch,
		SourceBranch: mr.SourceBranch,
		Title:        mr.Title,
		State:        mr.State,
		CreatedAt:    metav1.NewTime(*mr.CreatedAt),
		UpdatedAt:    metav1.NewTime(*mr.UpdatedAt),
		Author:       mr.Author.Username,
		Labels:       mr.Labels,
		Draft:        mr.Draft,
	}
}

func (g *PullRequestGitlabGenerator) BuildContexts() ([]*GeneratedContext, error) {
	opt := &gitlab.ListProjectMergeRequestsOptions{}
	if g.state != "" {
		opt.State = &g.state
	}
	if len(g.labels) != 0 {
		l := gitlab.Labels(g.labels)
		opt.Labels = &l
	}
	mrs, _, err := g.client.MergeRequests.ListProjectMergeRequests(g.project, opt)
	if err != nil {
		return nil, err
	}

	mrs2 := make([]MergeRequestInfo, 0, len(mrs))
	for _, mr := range mrs {
		mrs2 = append(mrs2, g.convertMR(mr))
	}
	mrs2, err = filterPullRequests(mrs2, g.filters)
	if err != nil {
		return nil, err
	}

	var ret []*GeneratedContext
	for _, mr := range mrs2 {
		vars := uo.New()
		_ = vars.SetNestedField(mr, "mergeRequest")

		ret = append(ret, &GeneratedContext{
			Vars: vars,
		})
	}
	return ret, nil
}

func buildPullRequestGeneratorGitlab(ctx context.Context, client client.Client, namespace string, spec templatesv1alpha1.PullRequestGenerator) (Generator, error) {
	var secret v1.Secret
	err := client.Get(ctx, types.NamespacedName{
		Namespace: namespace,
		Name:      spec.Gitlab.TokenRef.SecretName,
	}, &secret)
	if err != nil {
		return nil, err
	}

	token, ok := secret.Data[spec.Gitlab.TokenRef.Key]
	if !ok {
		return nil, fmt.Errorf("key %s not found in secret %s", spec.Gitlab.TokenRef.Key, spec.Gitlab.TokenRef.SecretName)
	}

	var opts []gitlab.ClientOptionFunc
	if spec.Gitlab.API != "" {
		opts = append(opts, gitlab.WithBaseURL(spec.Gitlab.API))
	}

	gitlabClient, err := gitlab.NewClient(string(token), opts...)
	if err != nil {
		return nil, err
	}
	ret := &PullRequestGitlabGenerator{
		client:  gitlabClient,
		project: spec.Gitlab.Project,
		labels:  spec.Gitlab.Labels,
		state:   spec.Gitlab.PullRequestState,
		filters: spec.Filters,
	}
	return ret, nil
}
