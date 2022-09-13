package handlers

import (
	"context"
	"github.com/kluctl/template-controller/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Handler interface {
	Handle(ctx context.Context, client client.Client, obj client.Object, status *v1alpha1.HandlerStatus) error
}
