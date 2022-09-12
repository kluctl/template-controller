package reporters

import (
	"context"
	"github.com/kluctl/template-controller/api/v1alpha1"
	"github.com/kluctl/template-controller/controllers"
	"github.com/kluctl/template-controller/controllers/status/comments"
	"github.com/kluctl/template-controller/controllers/status/webgit"
	"github.com/xanzy/go-gitlab"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PullRequestCommentReporter struct {
	mr        webgit.MergeRequestInterface
	clusterId string
}

func BuildPullRequestCommentReporter(ctx context.Context, client client.Client, namespace string, spec v1alpha1.PullRequestCommentReporter) (Reporter, error) {
	mr, err := buildWebgitMergeRequest(ctx, client, namespace, &spec)
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

func (p *PullRequestCommentReporter) Report(ctx context.Context, client client.Client, obj client.Object, status *v1alpha1.ReporterStatus) error {
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

	err = p.reconcileComment(ctx, obj, comment, status.PullRequestComment)
	if err != nil {
		return err
	}

	return nil
}

func (p *PullRequestCommentReporter) reconcileComment(ctx context.Context, obj client.Object, statusComment string, status *v1alpha1.PullRequestCommentReporterStatus) error {
	var err error
	body := generateMarkerComment("pull-request-comment", p.clusterId, obj.GetNamespace(), obj.GetName()) + "\n" + statusComment

	var existingNote webgit.Note
	if status.NoteId == "" {
		existingNote, err = p.findNote(obj)
		if err != nil {
			return err
		}
		if existingNote != nil {
			status.NoteId = existingNote.GetId()
			status.LastPostedStatusHash = ""
		} else {
			existingNote, err = p.mr.CreateMergeRequestNote(body)
			if err != nil {
				return err
			}
			status.NoteId = existingNote.GetId()
			status.LastPostedStatusHash = controllers.Sha256String(body)
		}
	} else {
		var resp *gitlab.Response
		existingNote, err = p.mr.GetMergeRequestNote(status.NoteId)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				status.NoteId = ""
				status.LastPostedStatusHash = ""
			}
			return err
		}
	}

	if status.LastPostedStatusHash == controllers.Sha256String(body) {
		return nil
	}

	err = existingNote.UpdateBody(body)
	if err != nil {
		status.NoteId = ""
		status.LastPostedStatusHash = ""
		return err
	}
	status.LastPostedStatusHash = controllers.Sha256String(body)
	return nil
}

func (p *PullRequestCommentReporter) findNote(obj client.Object) (webgit.Note, error) {
	notes, err := p.mr.ListMergeRequestNotes()
	if err != nil {
		return nil, err
	}
	for _, n := range notes {
		if !hasMarkerComment(n.GetBody(), "pull-request-comment", p.clusterId, obj.GetNamespace(), obj.GetName()) {
			continue
		}
		return n, nil
	}
	return nil, nil
}
