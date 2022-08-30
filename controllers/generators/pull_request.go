package generators

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	templatesv1alpha1 "kluctl/template-controller/api/v1alpha1"
	"regexp"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MergeRequestInfo struct {
	ID           int         `yaml:"id"`
	IID          int         `yaml:"IID"`
	TargetBranch string      `yaml:"targetBranch"`
	SourceBranch string      `yaml:"sourceBranch"`
	Title        string      `yaml:"title"`
	State        string      `yaml:"state"`
	CreatedAt    metav1.Time `yaml:"createdAt"`
	UpdatedAt    metav1.Time `yaml:"updatedAt"`
	Author       string      `yaml:"author"`
	Labels       []string    `yaml:"labels"`
	Draft        bool        `yaml:"draft"`
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
