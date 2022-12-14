package comments

import (
	"context"
	"github.com/kluctl/template-controller/controllers/objecthandler/comments/templates"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CommentGenerator interface {
	GenerateComment(ctx context.Context, obj client.Object) (string, error)
}

var byGroupKind = map[schema.GroupKind]CommentGenerator{
	schema.GroupKind{Group: "flux.kluctl.io", Kind: "KluctlDeployment"}: &TemplateComment{template: templates.MustGetTemplate("kluctldeployment.md.jinja2")},
}
var genericGenerator = &TemplateComment{template: templates.MustGetTemplate("generic.md.jinja2")}

func GetCommentGenerator(obj client.Object) (CommentGenerator, error) {
	gk := obj.GetObjectKind().GroupVersionKind().GroupKind()
	generator, ok := byGroupKind[gk]
	if ok {
		return generator, nil
	}
	return genericGenerator, nil
}
