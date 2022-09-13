package comments

import (
	"context"
	"github.com/kluctl/go-jinja2"
	"github.com/kluctl/template-controller/controllers"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type TemplateComment struct {
	template string
}

func (c *TemplateComment) GenerateComment(ctx context.Context, obj client.Object) (string, error) {
	j2, err := controllers.NewJinja2()
	if err != nil {
		return "", err
	}
	defer j2.Close()

	vars := map[string]any{}

	u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return "", err
	}

	vars["object"] = u

	rendered, err := j2.RenderString(c.template, jinja2.WithGlobals(vars))
	if err != nil {
		return "", err
	}
	return rendered, nil
}
