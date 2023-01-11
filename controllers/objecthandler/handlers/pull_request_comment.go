package handlers

import (
	"context"
	"github.com/kluctl/template-controller/api/v1alpha1"
	"github.com/kluctl/template-controller/controllers/objecthandler/comments"
	"github.com/kluctl/template-controller/controllers/webgit"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PullRequestCommentReporter struct {
	mr        webgit.MergeRequestInterface
	clusterId string
}

func BuildPullRequestCommentReporter(ctx context.Context, client client.Client, namespace string, spec v1alpha1.PullRequestCommentReporter) (Handler, error) {
	mr, err := webgit.BuildWebgitMergeRequest(ctx, client, namespace, spec.PullRequestRefHolder)
	if err != nil {
		return nil, err
	}

	clusterId, err := getClusterId(ctx, client)
	if err != nil {
		return nil, err
	}

	return &PullRequestCommentReporter{
		mr:        mr,
		clusterId: clusterId,
	}, nil
}

func (p *PullRequestCommentReporter) Handle(ctx context.Context, client client.Client, obj *unstructured.Unstructured, status *v1alpha1.HandlerStatus) error {
	if status.PullRequestComment == nil {
		status.PullRequestComment = &v1alpha1.PullRequestCommentReporterStatus{}
	}

	generator, err := comments.GetCommentGenerator(obj)
	if err != nil {
		return err
	}

	comment, err := generator.GenerateComment(ctx, obj)
	if err != nil {
		return err
	}

	err = p.reconcileComment(obj, comment, status.PullRequestComment)
	if err != nil {
		return err
	}

	return nil
}

func (p *PullRequestCommentReporter) reconcileComment(obj client.Object, statusComment string, status *v1alpha1.PullRequestCommentReporterStatus) error {
	return reconcileComment(p.clusterId, p.mr, "pull-request-comment", obj, statusComment, &status.NoteId, &status.LastPostedStatusHash)
}
