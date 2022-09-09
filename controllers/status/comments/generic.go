package comments

import (
	"context"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type GenericComment struct {
}

func (c *GenericComment) GenerateComment(ctx context.Context, obj client.Object) (string, error) {
	return "", nil
}
