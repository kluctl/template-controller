package generators

import (
	"context"
	"fmt"
	templatesv1alpha1 "github.com/kluctl/template-controller/api/v1alpha1"
	"github.com/kluctl/template-controller/controllers/objecthandler/webgit"
	"k8s.io/apimachinery/pkg/runtime"
	"regexp"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PullRequestGenerator struct {
	project webgit.ProjectInterface
	spec    templatesv1alpha1.PullRequestGenerator
}

func BuildPullRequestGenerator(ctx context.Context, client client.Client, namespace string, spec templatesv1alpha1.PullRequestGenerator, defaults *templatesv1alpha1.ObjectTemplateDefaultsSpec) (Generator, error) {
	p, err := webgit.BuildWebgit(ctx, client, namespace, spec, defaults)
	if err != nil {
		return nil, err
	}

	return &PullRequestGenerator{
		project: p,
		spec:    spec,
	}, nil
}

func (g *PullRequestGenerator) BuildContexts() ([]*GeneratedContext, error) {
	var err error
	var mrs []webgit.MergeRequestInterface

	var state templatesv1alpha1.MergeRequestState
	var targetBranch, sourceBranch *string
	var labels []string

	if g.spec.Gitlab != nil {
		state = g.spec.Gitlab.PullRequestState
		targetBranch = g.spec.Gitlab.TargetBranch
		sourceBranch = g.spec.Gitlab.SourceBranch
		labels = g.spec.Gitlab.Labels
	} else if g.spec.Github != nil {
		state = g.spec.Github.PullRequestState
		targetBranch = g.spec.Github.TargetBranch
		sourceBranch = g.spec.Github.SourceBranch
		labels = g.spec.Github.Labels
	} else {
		return nil, fmt.Errorf("no pull request provider spec specified")
	}

	var targetBranchRE, sourceBranchRE *regexp.Regexp
	if targetBranch != nil {
		targetBranchRE, err = regexp.Compile(*targetBranch)
		if err != nil {
			return nil, err
		}
	}
	if sourceBranch != nil {
		sourceBranchRE, err = regexp.Compile(*sourceBranch)
		if err != nil {
			return nil, err
		}
	}

	mrs, err = g.project.ListMergeRequests(state)
	if err != nil {
		return nil, err
	}

	var ret []*GeneratedContext

outer:
	for _, mr := range mrs {
		info, err := mr.Info()
		if err != nil {
			return nil, err
		}
		if targetBranchRE != nil && !targetBranchRE.MatchString(info.TargetBranch) {
			continue
		}
		if sourceBranchRE != nil && !sourceBranchRE.MatchString(info.SourceBranch) {
			continue
		}
		labelsMap := map[string]bool{}
		for _, l := range info.Labels {
			labelsMap[l] = true
		}
		for _, l := range labels {
			if _, ok := labelsMap[l]; !ok {
				continue outer
			}
		}

		u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&info)
		if err != nil {
			return nil, err
		}

		vars := map[string]any{
			"mergeRequest": u,
		}

		ret = append(ret, &GeneratedContext{
			Vars: vars,
		})
	}

	return ret, nil
}
