package comments

import (
	"context"
	"fmt"
	templatesv1alpha1 "github.com/kluctl/template-controller/api/v1alpha1"
	"github.com/kluctl/template-controller/controllers"
	"github.com/kluctl/template-controller/controllers/objecthandler/webgit"
	"github.com/xanzy/go-gitlab"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"strings"
)

type BaseCommentReconciler struct {
	client.Client
	Scheme       *runtime.Scheme
	FieldManager string
}

type GetCommentSourceSpec interface {
	client.Object
	GetCommentSourceSpec() *templatesv1alpha1.CommentSourceSpec
}

type ItemList interface {
	client.ObjectList
	GetItems() []client.Object
}

func (r *BaseCommentReconciler) getText(ctx context.Context, spec templatesv1alpha1.CommentSourceSpec, objNs string) (string, error) {
	if spec.Text != nil {
		return *spec.Text, nil
	} else if spec.TextTemplate != nil {
		var tt templatesv1alpha1.TextTemplate
		err := r.Client.Get(ctx, client.ObjectKey{
			Name:      spec.TextTemplate.Name,
			Namespace: objNs,
		}, &tt)
		if err != nil {
			return "", err
		}
		rc := meta.FindStatusCondition(tt.GetConditions(), "Ready")
		if rc == nil {
			return "", fmt.Errorf("TextTemplate %s has no ready condition yet", spec.TextTemplate.Name)
		}
		if rc.Status != metav1.ConditionTrue {
			return "", fmt.Errorf("TextTemplate %s is not ready yet: reason=%s, message=%s", spec.TextTemplate.Name, rc.Reason, rc.Message)
		}
		return tt.Status.Result, nil
	} else if spec.ConfigMap != nil {
		var c corev1.ConfigMap
		err := r.Client.Get(ctx, client.ObjectKey{
			Name:      spec.ConfigMap.Name,
			Namespace: objNs,
		}, &c)
		if err != nil {
			return "", err
		}
		cv, ok := c.Data[spec.ConfigMap.Key]
		if !ok {
			return "", fmt.Errorf("ConfigMap %s does not contain key %s", spec.ConfigMap.Name, spec.ConfigMap.Key)
		}
		return cv, nil
	} else {
		return "", fmt.Errorf("no template specified")
	}
}

func (r *BaseCommentReconciler) reconcileComment(ctx context.Context, mr webgit.MergeRequestInterface, tag string, commentId *string, obj client.Object, noteId *string, lastPostedBodyHash *string) error {
	clusterId, err := r.getClusterId(ctx)
	if err != nil {
		return err
	}

	spec := obj.(GetCommentSourceSpec).GetCommentSourceSpec()

	comment, err := r.getText(ctx, *spec, obj.GetNamespace())
	if err != nil {
		return err
	}

	body := r.generateMarkerComment(clusterId, tag, commentId, obj.GetNamespace(), obj.GetName()) + "\n" + comment

	var existingNote webgit.Note
	if *noteId == "" {
		existingNote, err = r.findNote(mr, clusterId, commentId, tag, obj)
		if err != nil {
			return err
		}
		if existingNote != nil {
			*noteId = existingNote.GetId()
			*lastPostedBodyHash = ""
		} else {
			existingNote, err = mr.CreateMergeRequestNote(body)
			if err != nil {
				return err
			}
			*noteId = existingNote.GetId()
			*lastPostedBodyHash = controllers.Sha256String(body)
		}
	} else {
		var resp *gitlab.Response
		existingNote, err = mr.GetMergeRequestNote(*noteId)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				*noteId = ""
				*lastPostedBodyHash = ""
			}
			return err
		}
	}

	if *lastPostedBodyHash == controllers.Sha256String(body) {
		return nil
	}

	err = existingNote.UpdateBody(body)
	if err != nil {
		*noteId = ""
		*lastPostedBodyHash = ""
		return err
	}
	*lastPostedBodyHash = controllers.Sha256String(body)
	return nil
}

func (r *BaseCommentReconciler) findNote(mr webgit.MergeRequestInterface, clusterId string, commentId *string, tag string, obj client.Object) (webgit.Note, error) {
	notes, err := mr.ListMergeRequestNotes()
	if err != nil {
		return nil, err
	}
	for _, n := range notes {
		if !r.hasMarkerComment(n.GetBody(), clusterId, tag, commentId, obj.GetNamespace(), obj.GetName()) {
			continue
		}
		return n, nil
	}
	return nil, nil
}

func (r *BaseCommentReconciler) getClusterId(ctx context.Context) (string, error) {
	var ns corev1.Namespace
	err := r.Client.Get(ctx, types.NamespacedName{Name: "kube-system"}, &ns)
	if err != nil {
		return "", err
	}
	return string(ns.UID), nil
}

func (r *BaseCommentReconciler) generateMarkerComment(clusterId string, tag string, commentId *string, objNamespace string, objName string) string {
	if commentId == nil {
		x := fmt.Sprintf("%s/%s", objNamespace, objName)
		commentId = &x
	}
	return fmt.Sprintf("<!-- template-controller-%s \"%s\" \"%s\" -->", tag, clusterId, *commentId)
}

func (r *BaseCommentReconciler) hasMarkerComment(body string, clusterId string, tag string, commentId *string, objNamespace string, objName string) bool {
	expected := r.generateMarkerComment(clusterId, tag, commentId, objNamespace, objName)
	for _, line := range strings.Split(body, "\n") {
		if line == expected {
			return true
		}
	}
	return false
}

// SetupWithManager sets up the controller with the Manager.
func (r *BaseCommentReconciler) baseSetupWithManager(mgr ctrl.Manager, r2 reconcile.Reconciler, obj GetCommentSourceSpec, buildList func() ItemList) error {
	const indexKey = "spec.source.ref"

	if err := mgr.GetCache().IndexField(context.TODO(), obj, indexKey,
		func(obj client.Object) []string {
			spec := obj.(GetCommentSourceSpec).GetCommentSourceSpec()
			var ref templatesv1alpha1.ObjectRef
			if spec.ConfigMap != nil {
				ref.Kind = "ConfigMap"
				ref.Name = spec.ConfigMap.Name
			} else if spec.TextTemplate != nil {
				ref.Kind = "TextTemplate"
				ref.Name = spec.TextTemplate.Name
			} else {
				return nil
			}
			return []string{ref.String()}
		}); err != nil {
		return fmt.Errorf("failed setting index fields: %w", err)
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(obj, builder.WithPredicates(
			predicate.GenerationChangedPredicate{},
		)).
		Watches(
			&source.Kind{Type: &corev1.ConfigMap{}},
			handler.EnqueueRequestsFromMapFunc(func(obj client.Object) []reconcile.Request {
				list := buildList()

				var ref templatesv1alpha1.ObjectRef
				ref.Name = obj.GetName()
				switch obj.(type) {
				case *corev1.ConfigMap:
					ref.Kind = "ConfigMap"
				case *templatesv1alpha1.TextTemplate:
					ref.Kind = "TextTemplate"
				default:
					return nil
				}
				if err := r.List(context.Background(), list, client.MatchingFields{
					indexKey: ref.String(),
				}); err != nil {
					return nil
				}
				var ret []reconcile.Request
				for _, o := range list.GetItems() {
					ret = append(ret, reconcile.Request{NamespacedName: client.ObjectKeyFromObject(o)})
				}
				return ret
			}),
		).
		Complete(r2)
}
