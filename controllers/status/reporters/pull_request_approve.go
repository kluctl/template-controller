package reporters

import (
	"context"
	"github.com/kluctl/template-controller/api/v1alpha1"
	"github.com/kluctl/template-controller/controllers/status/webgit"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/json"
	"sigs.k8s.io/cli-utils/pkg/kstatus/status"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PullRequestApproveReporter struct {
	mr   webgit.MergeRequestInterface
	spec v1alpha1.PullRequestApproveReporter
}

func BuildPullRequestApproveReporter(ctx context.Context, client client.Client, namespace string, spec v1alpha1.PullRequestApproveReporter) (Reporter, error) {
	mr, err := buildWebgitMergeRequest(ctx, client, namespace, &spec)
	if err != nil {
		return nil, err
	}

	return &PullRequestApproveReporter{mr: mr, spec: spec}, nil
}

func (p *PullRequestApproveReporter) Report(ctx context.Context, obj client.Object, status *v1alpha1.ReporterStatus) error {
	if status.PullRequestApprove == nil {
		status.PullRequestApprove = &v1alpha1.PullRequestApproveReporterStatus{}

		approved, err := p.mr.HasApproved()
		if err != nil {
			return err
		}
		status.PullRequestApprove.Approved = &approved
	}

	ready, err := p.computeReady(obj)
	if err != nil {
		return err
	}

	if ready && !*status.PullRequestApprove.Approved {
		err = p.mr.Approve()
		if err != nil {
			return err
		}
		b := true
		status.PullRequestApprove.Approved = &b
	} else if !ready && *status.PullRequestApprove.Approved {
		err = p.mr.Unapprove()
		if err != nil {
			return err
		}
		b := false
		status.PullRequestApprove.Approved = &b
	}
	return nil
}

func (p *PullRequestApproveReporter) computeReady(obj client.Object) (bool, error) {
	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		b, err := json.Marshal(obj)
		if err != nil {
			return false, err
		}
		u = &unstructured.Unstructured{}
		err = json.Unmarshal(b, u)
		if err != nil {
			return false, err
		}
	}
	res, err := status.Compute(u)
	if err != nil {
		return false, err
	}
	if res.Status == status.CurrentStatus && p.spec.MissingStatusIsError {
		if _, ok := u.Object["status"]; !ok {
			return false, nil
		}
	}
	return res.Status == status.CurrentStatus, nil
}
