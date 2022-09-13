package handlers

import (
	"context"
	"fmt"
	"github.com/fluxcd/pkg/apis/meta"
	"github.com/kluctl/template-controller/api/v1alpha1"
	"github.com/kluctl/template-controller/controllers/status/webgit"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const KluctlDeployRequestAnnotation = "deploy.flux.kluctl.io/requestedAt"

type PullRequestCommandHandler struct {
	mr   webgit.MergeRequestInterface
	spec v1alpha1.PullRequestCommandHandler

	clusterId string
}

func BuildPullRequestCommandHandler(ctx context.Context, client client.Client, namespace string, spec v1alpha1.PullRequestCommandHandler) (Handler, error) {
	mr, err := buildWebgitMergeRequest(ctx, client, namespace, &spec)
	if err != nil {
		return nil, err
	}

	clusterId, err := getClusterId(ctx, client)
	if err != nil {
		return nil, err
	}

	return &PullRequestCommandHandler{mr: mr, spec: spec, clusterId: clusterId}, nil
}

func (p *PullRequestCommandHandler) Handle(ctx context.Context, client client.Client, obj client.Object, status *v1alpha1.HandlerStatus) error {
	if status.PullRequestCommand == nil {
		status.PullRequestCommand = &v1alpha1.PullRequestCommandHandlerStatus{}
	}

	var origLastTime time.Time
	if status.PullRequestCommand.LastProcessedCommentTime != nil {
		x, err := time.Parse(time.RFC3339Nano, *status.PullRequestCommand.LastProcessedCommentTime)
		if err == nil {
			origLastTime = x
		}
	}

	newLastTime := origLastTime

	unprocessedNotes, err := p.mr.ListMergeRequestNotesAfter(newLastTime)
	if err != nil {
		return err
	}
	if len(unprocessedNotes) == 0 {
		return nil
	}

	updateStatus := func() {
		if origLastTime != newLastTime {
			x := newLastTime.Format(time.RFC3339Nano)
			status.PullRequestCommand.LastProcessedCommentTime = &x
		}
	}

	for _, n := range unprocessedNotes {
		err = p.processGitlabStatusCommand(ctx, client, n, obj)
		if err != nil {
			updateStatus()
			break
		}
		newLastTime = n.GetCreatedAt()
	}
	updateStatus()

	return nil
}

func (p *PullRequestCommandHandler) processGitlabStatusCommand(ctx context.Context, c client.Client, n webgit.Note, obj client.Object) error {
	body := n.GetBody()
	if hasMarkerComment(body, "pull-request-command-processed", p.clusterId, obj.GetNamespace(), obj.GetName()) {
		return nil
	}

	addTimeAnnotation := func(n string) error {
		patch := client.MergeFrom(obj.DeepCopyObject().(client.Object))

		a := obj.GetAnnotations()
		if a == nil {
			a = make(map[string]string)
		}
		a[n] = time.Now().Format(time.RFC3339)
		obj.SetAnnotations(a)
		err := c.Patch(ctx, obj, patch)
		if err != nil {
			return err
		}
		return nil
	}

	if body == "/reconcile" {
		err := addTimeAnnotation(meta.ReconcileRequestAnnotation)
		if err != nil {
			return err
		}
	} else if body == "/deploy" {
		err := addTimeAnnotation(KluctlDeployRequestAnnotation)
		if err != nil {
			return err
		}
	} else {
		return nil
	}

	newBody := body
	newBody += fmt.Sprintf("\n\n:robot: Command has been processed at %s\n", time.Now().Format(time.RFC3339))
	newBody += generateMarkerComment("pull-request-command-processed", p.clusterId, obj.GetNamespace(), obj.GetName())

	err := n.UpdateBody(newBody)
	if err != nil {
		return err
	}

	return nil
}
