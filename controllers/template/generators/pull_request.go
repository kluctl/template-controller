package generators

import (
	"context"
	"fmt"
	templatesv1alpha1 "github.com/kluctl/template-controller/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"regexp"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MergeRequestInfo struct {
	ID           int         `json:"id"`
	TargetBranch string      `json:"targetBranch"`
	SourceBranch string      `json:"sourceBranch"`
	Title        string      `json:"title"`
	State        string      `json:"state"`
	CreatedAt    metav1.Time `json:"createdAt"`
	UpdatedAt    metav1.Time `json:"updatedAt"`
	Author       string      `json:"author"`
	Labels       []string    `json:"labels"`
	Draft        bool        `json:"draft"`
}

func BuildPullRequestGenerator(ctx context.Context, client client.Client, namespace string, spec templatesv1alpha1.PullRequestGenerator) (Generator, error) {
	if spec.Gitlab != nil {
		return buildPullRequestGeneratorGitlab(ctx, client, namespace, spec)
	} else {
		return nil, fmt.Errorf("no pullRequest generator specified")
	}
}

func filterPullRequests(mrs []MergeRequestInfo, filters []templatesv1alpha1.PullRequestGeneratorFilter) ([]MergeRequestInfo, error) {
	var branchMatchPatterns []*regexp.Regexp

	for _, f := range filters {
		var err error
		var p *regexp.Regexp
		if f.BranchMatch != nil {
			p, err = regexp.Compile(*f.BranchMatch)
			if err != nil {
				return nil, err
			}
		}

		branchMatchPatterns = append(branchMatchPatterns, p)
	}

	ret := make([]MergeRequestInfo, 0, len(mrs))
	for _, mr := range mrs {
		match := true
		for i, f := range filters {
			if f.BranchMatch != nil {
				if !branchMatchPatterns[i].MatchString(mr.SourceBranch) {
					match = false
					break
				}
			}
		}
		if match {
			ret = append(ret, mr)
		}
	}
	return ret, nil
}
