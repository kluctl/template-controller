package reporters

import (
	"context"
	"github.com/kluctl/template-controller/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Reporter interface {
	Report(ctx context.Context, client client.Client, obj client.Object, status *v1alpha1.ReporterStatus) error
}
