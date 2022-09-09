package status

import (
	"fmt"
	"github.com/kluctl/template-controller/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const forObjectIndexKey = "spec.forObject"

func buildRefIndexValue(ref v1alpha1.ObjectRef, ns string) string {
	if ref.Namespace != "" {
		ns = ref.Namespace
	}
	return fmt.Sprintf("%s/%s/%s", ref.Kind, ns, ref.Name)
}

func buildObjectIndexValue(obj client.Object) string {
	gvk := obj.GetObjectKind().GroupVersionKind()
	return fmt.Sprintf("%s/%s/%s", gvk.Kind, obj.GetNamespace(), obj.GetName())
}
