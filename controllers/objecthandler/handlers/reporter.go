package handlers

import (
	"context"
	"github.com/kluctl/template-controller/api/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Handler interface {
	Handle(ctx context.Context, client client.Client, obj *unstructured.Unstructured, status *v1alpha1.HandlerStatus) error
}
