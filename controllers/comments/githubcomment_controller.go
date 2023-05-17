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

package comments

import (
	"context"
	"github.com/kluctl/template-controller/controllers"
	"github.com/kluctl/template-controller/controllers/webgit"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	templatesv1alpha1 "github.com/kluctl/template-controller/api/v1alpha1"
)

// GithubCommentReconciler reconciles a GithubComment object
type GithubCommentReconciler struct {
	BaseCommentReconciler
}

//+kubebuilder:rbac:groups=templates.kluctl.io,resources=githubcomments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=templates.kluctl.io,resources=githubcomments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=templates.kluctl.io,resources=githubcomments/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch

func (r *GithubCommentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	logger := log.FromContext(ctx)

	logger.V(1).Info("Starting reconcile")
	defer logger.V(1).Info("Finished reconcile", "err", err)

	var gc templatesv1alpha1.GithubComment
	err = r.Get(ctx, req.NamespacedName, &gc)
	if err != nil {
		logger.Error(err, "Get failed")
		err = client.IgnoreNotFound(err)
		return
	}

	// Return early if the object is suspended.
	if gc.Spec.Suspend {
		logger.Info("Reconciliation is suspended for this object")
		return ctrl.Result{}, nil
	}

	patch := client.MergeFrom(gc.DeepCopy())
	err = r.doReconcile(ctx, &gc)
	if err != nil {
		c := metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionFalse,
			ObservedGeneration: gc.GetGeneration(),
			Reason:             "Error",
			Message:            err.Error(),
		}
		apimeta.SetStatusCondition(&gc.Status.Conditions, c)
	} else {
		c := metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionTrue,
			ObservedGeneration: gc.GetGeneration(),
			Reason:             "Success",
			Message:            "Success",
		}
		apimeta.SetStatusCondition(&gc.Status.Conditions, c)
	}
	err = r.Status().Patch(ctx, &gc, patch, controllers.SubResourceFieldOwner(r.FieldManager))
	return
}

func (r *GithubCommentReconciler) doReconcile(ctx context.Context, obj *templatesv1alpha1.GithubComment) error {
	mr, err := webgit.BuildWebgitMergeRequestGithub(ctx, r.Client, obj.GetNamespace(), obj.Spec.GithubPullRequestRef)
	if err != nil {
		return err
	}

	return r.reconcileComment(ctx, mr, "github-comment", obj.Spec.Id, obj, &obj.Status.CommentId, &obj.Status.LastPostedBodyHash)
}

// SetupWithManager sets up the controller with the Manager.
func (r *GithubCommentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return r.baseSetupWithManager(mgr, r, &templatesv1alpha1.GithubComment{}, func() ItemList {
		return &templatesv1alpha1.GithubCommentList{}
	})
}
