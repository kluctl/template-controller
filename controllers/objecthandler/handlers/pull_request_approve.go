package handlers

import (
	"context"
	"github.com/kluctl/template-controller/api/v1alpha1"
	"github.com/kluctl/template-controller/controllers/webgit"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PullRequestApproveReporter struct {
	mr   webgit.MergeRequestInterface
	spec v1alpha1.PullRequestApproveReporter
}

func BuildPullRequestApproveReporter(ctx context.Context, client client.Client, namespace string, spec v1alpha1.PullRequestApproveReporter) (Handler, error) {
	mr, err := webgit.BuildWebgitMergeRequest(ctx, client, namespace, spec.PullRequestRefHolder)
	if err != nil {
		return nil, err
	}

	return &PullRequestApproveReporter{mr: mr, spec: spec}, nil
}

func (p *PullRequestApproveReporter) Handle(ctx context.Context, client client.Client, obj *unstructured.Unstructured, status *v1alpha1.HandlerStatus) error {
	if status.PullRequestApprove == nil {
		status.PullRequestApprove = &v1alpha1.PullRequestApproveReporterStatus{}
	}

	approved, err := p.mr.HasApproved()
	if err != nil {
		return err
	}
	status.PullRequestApprove.Approved = &approved

	ready, err := p.computeReady(ctx, client, obj)
	if err != nil {
		return err
	}

	if ready && !approved {
		err = p.mr.Approve()
		if err != nil {
			return err
		}
		b := true
		status.PullRequestApprove.Approved = &b
	} else if !ready && approved {
		err = p.mr.Unapprove()
		if err != nil {
			return err
		}
		b := false
		status.PullRequestApprove.Approved = &b
	}
	return nil
}

func (p *PullRequestApproveReporter) computeReady(ctx context.Context, client client.Client, obj client.Object) (bool, error) {
	sc := StatusCalculator{Client: client}
	return sc.ComputeReady(ctx, obj, p.spec.MissingReadyConditionIsError)
}
