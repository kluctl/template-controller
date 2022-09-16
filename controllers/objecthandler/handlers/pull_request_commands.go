package handlers

import (
	"context"
	"fmt"
	"github.com/kluctl/go-jinja2"
	"github.com/kluctl/template-controller/api/v1alpha1"
	"github.com/kluctl/template-controller/controllers"
	"github.com/kluctl/template-controller/controllers/objecthandler/comments/templates"
	"github.com/kluctl/template-controller/controllers/objecthandler/webgit"
	"k8s.io/apimachinery/pkg/runtime"
	"regexp"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type PullRequestCommandHandler struct {
	mr   webgit.MergeRequestInterface
	spec v1alpha1.PullRequestCommandHandler

	clusterId string
}

func BuildPullRequestCommandHandler(ctx context.Context, client client.Client, namespace string, spec v1alpha1.PullRequestCommandHandler, defaults *v1alpha1.ObjectHandlerDefaultsSpec) (Handler, error) {
	mr, err := webgit.BuildWebgitMergeRequest(ctx, client, namespace, &spec, defaults)
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
	j2, err := controllers.NewJinja2()
	if err != nil {
		return err
	}
	defer j2.Close()

	if status.PullRequestCommand == nil {
		status.PullRequestCommand = &v1alpha1.PullRequestCommandHandlerStatus{}
	}

	err = p.reconcileHelpComment(j2, obj, status)
	if err != nil {
		return err
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
		err = p.processCommand(ctx, j2, client, n, obj)
		if err != nil {
			updateStatus()
			break
		}
		newLastTime = n.GetCreatedAt()
	}
	updateStatus()

	return nil
}

var helpCommandTemplate = templates.MustGetTemplate("commandhelp.md.jinja2")

func (p *PullRequestCommandHandler) reconcileHelpComment(j2 *jinja2.Jinja2, obj client.Object, status *v1alpha1.HandlerStatus) error {
	if !p.spec.PostHelpComment {
		return nil
	}

	vars := map[string]any{}

	vars["spec"] = &p.spec

	comment, err := j2.RenderString(helpCommandTemplate, jinja2.WithGlobals(vars))
	if err != nil {
		return err
	}

	err = reconcileComment(p.clusterId, p.mr, "pull-request-command-help", obj, comment, &status.PullRequestCommand.HelpNoteId, &status.PullRequestCommand.HelpNoteBodyHash)
	if err != nil {
		return err
	}

	return nil
}

var commandRegex = regexp.MustCompile("^/([a-zA-Z][a-zA-Z0-9]*)$")

func (p *PullRequestCommandHandler) processCommand(ctx context.Context, j2 *jinja2.Jinja2, c client.Client, n webgit.Note, obj client.Object) error {
	body := n.GetBody()
	if hasMarkerComment(body, "pull-request-command-processed", p.clusterId, obj.GetNamespace(), obj.GetName()) {
		return nil
	}

	m := commandRegex.FindStringSubmatch(body)
	if m == nil {
		return nil
	}
	commandName := m[1]

	found := false
	var err error
	for _, command := range p.spec.Commands {
		if command.Name == commandName {
			found = true
			err = p.handleCommand(ctx, j2, c, obj, command)
			break
		}
	}
	if !found {
		return nil
	}

	newBody := body
	newBody += fmt.Sprintf("\n\n:robot: Command has been processed at %s\n", time.Now().Format(time.RFC3339))
	if err != nil {
		newBody += fmt.Sprintf("<br>:boom: Command failed with error: %s\n", err.Error())
	}
	newBody += generateMarkerComment("pull-request-command-processed", p.clusterId, obj.GetNamespace(), obj.GetName())

	err = n.UpdateBody(newBody)
	if err != nil {
		return err
	}

	return nil
}

func (p *PullRequestCommandHandler) handleCommand(ctx context.Context, j2 *jinja2.Jinja2, c client.Client, obj client.Object, command v1alpha1.PullRequestCommandHandlerCommandSpec) error {
	for _, action := range command.Actions {
		if action.Annotate != nil {
			err := p.handleActionAnnotate(j2, obj, action.Annotate)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("action is missing in command %s", command.Name)
		}
	}
	return nil
}

func (p *PullRequestCommandHandler) handleActionAnnotate(j2 *jinja2.Jinja2, obj client.Object, action *v1alpha1.PullRequestCommandHandlerActionAnnotateSpec) error {
	vars := map[string]any{}

	u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return err
	}

	vars["object"] = u

	name, err := j2.RenderString(action.Annotation, jinja2.WithGlobals(vars))
	if err != nil {
		return err
	}
	value, err := j2.RenderString(action.Value, jinja2.WithGlobals(vars))
	if err != nil {
		return err
	}

	a := obj.GetAnnotations()
	if a == nil { //asd
		a = map[string]string{}
	}
	a[name] = value
	obj.SetAnnotations(a)

	return nil
}
