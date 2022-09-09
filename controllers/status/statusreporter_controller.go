/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package status

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/kluctl/template-controller/controllers/status/reporters"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sync"

	templatesv1alpha1 "github.com/kluctl/template-controller/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// StatusReporterReconciler reconciles a StatusReporter object
type StatusReporterReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	controller   controller.Controller
	watchedKinds map[schema.GroupVersionKind]bool
	mutex        sync.Mutex
}

//+kubebuilder:rbac:groups=status.kluctl.io,resources=statusreporters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=status.kluctl.io,resources=statusreporters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=status.kluctl.io,resources=statusreporters/finalizers,verbs=update

// Reconcile a resource
func (r *StatusReporterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var sr templatesv1alpha1.StatusReporter
	err := r.Get(ctx, req.NamespacedName, &sr)
	if err != nil {
		return ctrl.Result{}, err
	}

	patch := client.MergeFrom(sr.DeepCopy())
	err = r.doReconcile(ctx, &sr)
	if err != nil {
		c := metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionFalse,
			ObservedGeneration: sr.GetGeneration(),
			Reason:             "Error",
			Message:            err.Error(),
		}
		apimeta.SetStatusCondition(&sr.Status.Conditions, c)
	} else {
		c := metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionTrue,
			ObservedGeneration: sr.GetGeneration(),
			Reason:             "Success",
			Message:            "Success",
		}
		apimeta.SetStatusCondition(&sr.Status.Conditions, c)
	}
	err = r.Status().Patch(ctx, &sr, patch)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{
		RequeueAfter: sr.Spec.Interval.Duration,
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *StatusReporterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.watchedKinds = map[schema.GroupVersionKind]bool{}

	// Index the StatusReporter by the objects they are for.
	if err := mgr.GetCache().IndexField(context.TODO(), &templatesv1alpha1.StatusReporter{}, forObjectIndexKey,
		func(object client.Object) []string {
			sr := object.(*templatesv1alpha1.StatusReporter)
			return []string{
				buildRefIndexValue(sr.Spec.ForObject, sr.GetNamespace()),
			}
		}); err != nil {
		return fmt.Errorf("failed setting index fields: %w", err)
	}

	c, err := ctrl.NewControllerManagedBy(mgr).
		For(&templatesv1alpha1.StatusReporter{}).
		Build(r)
	if err != nil {
		return err
	}
	r.controller = c

	return nil
}

func (r *StatusReporterReconciler) doReconcile(ctx context.Context, sr *templatesv1alpha1.StatusReporter) error {
	err := r.addWatchForKind(ctx, sr)
	if err != nil {
		return err
	}

	gvk, err := sr.Spec.ForObject.GroupVersionKind()
	if err != nil {
		return err
	}

	name := types.NamespacedName{
		Name:      sr.Spec.ForObject.Name,
		Namespace: sr.GetNamespace(),
	}
	if sr.Spec.ForObject.Namespace != "" {
		name.Namespace = sr.Spec.ForObject.Namespace
	}

	var obj unstructured.Unstructured
	obj.SetGroupVersionKind(gvk)

	err = r.Client.Get(ctx, name, &obj)
	if err != nil {
		return err
	}

	existingStatuses := map[string]bool{}

	var errs *multierror.Error
	for _, spec := range sr.Spec.Reporters {
		var reporter reporters.Reporter
		if spec.PullRequestComment != nil {
			reporter, err = reporters.BuildPullRequestCommentReporter(ctx, r.Client, sr.GetNamespace(), *spec.PullRequestComment)
		} else if spec.PullRequestApprove != nil {
			reporter, err = reporters.BuildPullRequestApproveReporter(ctx, r.Client, sr.GetNamespace(), *spec.PullRequestApprove)
		} else {
			return fmt.Errorf("no reporter specified")
		}
		if err != nil {
			return err
		}

		key := spec.BuildKey()
		existingStatuses[key] = true

		var status *templatesv1alpha1.ReporterStatus
		for _, x := range sr.Status.ReporterStatus {
			if x.Key == key {
				status = x
				break
			}
		}
		if status == nil {
			status = &templatesv1alpha1.ReporterStatus{
				Key: key,
			}
			sr.Status.ReporterStatus = append(sr.Status.ReporterStatus, status)
		}

		err = reporter.Report(ctx, &obj, status)
		if err != nil {
			errs = multierror.Append(errs, err)
			status.Error = err.Error()
		} else {
			status.Error = ""
		}
	}

	old := sr.Status.ReporterStatus
	sr.Status.ReporterStatus = nil
	for _, x := range old {
		if a, _ := existingStatuses[x.Key]; a {
			sr.Status.ReporterStatus = append(sr.Status.ReporterStatus, x)
		}
	}

	return errs.ErrorOrNil()
}

func (r *StatusReporterReconciler) addWatchForKind(ctx context.Context, sr *templatesv1alpha1.StatusReporter) error {
	gvk, err := sr.Spec.ForObject.GroupVersionKind()
	if err != nil {
		return err
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if x, ok := r.watchedKinds[gvk]; ok && x {
		return nil
	}

	var dummy unstructured.Unstructured
	dummy.SetGroupVersionKind(gvk)

	err = r.controller.Watch(&source.Kind{Type: &dummy}, handler.EnqueueRequestsFromMapFunc(func(object client.Object) []reconcile.Request {
		var list templatesv1alpha1.StatusReporterList
		err := r.List(context.Background(), &list, client.MatchingFields{
			forObjectIndexKey: buildObjectIndexValue(object),
		})
		if err != nil {
			return nil
		}
		var reqs []reconcile.Request
		for _, x := range list.Items {
			reqs = append(reqs, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: x.Namespace,
					Name:      x.Name,
				},
			})
		}
		return reqs
	}))
	if err != nil {
		return err
	}

	r.watchedKinds[gvk] = true
	return nil
}
