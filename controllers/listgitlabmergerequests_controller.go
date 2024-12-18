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

package controllers

import (
	"context"
	"fmt"
	templatesv1alpha1 "github.com/kluctl/template-controller/api/v1alpha1"
	gitlab "gitlab.com/gitlab-org/api/client-go"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/json"
	"regexp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sort"
)

// ListGitlabMergeRequestsReconciler reconciles a ListGitlabMergeRequests object
type ListGitlabMergeRequestsReconciler struct {
	client.Client
	Scheme       *runtime.Scheme
	FieldManager string
}

//+kubebuilder:rbac:groups=templates.kluctl.io,resources=listgitlabmergerequests,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=templates.kluctl.io,resources=listgitlabmergerequests/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=templates.kluctl.io,resources=listgitlabmergerequests/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

func (r *ListGitlabMergeRequestsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var obj templatesv1alpha1.ListGitlabMergeRequests
	err := r.Get(ctx, req.NamespacedName, &obj)
	if err != nil {
		err = client.IgnoreNotFound(err)
		if err != nil {
			logger.Error(err, "Get failed")
		}
		return ctrl.Result{}, err
	}

	err = r.doReconcile(ctx, &obj)
	if err != nil {
		c := metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionFalse,
			ObservedGeneration: obj.GetGeneration(),
			Reason:             "Error",
			Message:            err.Error(),
		}
		apimeta.SetStatusCondition(&obj.Status.Conditions, c)
	} else {
		c := metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionTrue,
			ObservedGeneration: obj.GetGeneration(),
			Reason:             "Success",
			Message:            "Success",
		}
		apimeta.SetStatusCondition(&obj.Status.Conditions, c)
	}

	// TODO optimize the update as it currently causes to update all merge requests on every call
	// patching is not working very well as causes nulls to be pruned and full array replacement for every single change
	err = r.Status().Update(ctx, &obj, SubResourceFieldOwner(r.FieldManager))
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{
		RequeueAfter: obj.Spec.Interval.Duration,
	}, nil
}

func (r *ListGitlabMergeRequestsReconciler) doReconcile(ctx context.Context, obj *templatesv1alpha1.ListGitlabMergeRequests) error {
	var token string
	var err error

	if obj.Spec.TokenRef != nil {
		token, err = GetSecretToken(ctx, r.Client, obj.Namespace, *obj.Spec.TokenRef)
		if err != nil {
			return err
		}
	}

	sourceBranchRegex := regexp.MustCompile(".*")
	targetBrachRegex := regexp.MustCompile(".*")

	if obj.Spec.SourceBranch != nil {
		sourceBranchRegex, err = regexp.Compile(fmt.Sprintf("^%s$", *obj.Spec.SourceBranch))
		if err != nil {
			return err
		}
	}
	if obj.Spec.TargetBranch != nil {
		targetBrachRegex, err = regexp.Compile(fmt.Sprintf("^%s$", *obj.Spec.TargetBranch))
		if err != nil {
			return err
		}
	}

	var opts []gitlab.ClientOptionFunc
	if obj.Spec.API != nil {
		opts = append(opts, gitlab.WithBaseURL(*obj.Spec.API))
	}
	gl, err := gitlab.NewClient(token, opts...)
	if err != nil {
		return err
	}

	labels := gitlab.LabelOptions(obj.Spec.Labels)
	if len(labels) == 0 {
		labels = nil
	}

	listOpts := &gitlab.ListProjectMergeRequestsOptions{
		Labels: &labels,
		State:  obj.Spec.State,
	}
	listOpts.Page = 1
	listOpts.PerPage = 100

	var projectId any
	switch obj.Spec.Project.Type {
	case intstr.Int:
		projectId = obj.Spec.Project.IntValue()
	case intstr.String:
		projectId = obj.Spec.Project.String()
	default:
		return fmt.Errorf("invalid Project value: neither int nor string")
	}

	var result []*gitlab.MergeRequest
	for true {
		if len(result)+listOpts.PerPage > obj.Spec.Limit {
			listOpts.PerPage = obj.Spec.Limit - len(result)
		}

		page, _, err := gl.MergeRequests.ListProjectMergeRequests(projectId, listOpts, gitlab.WithContext(ctx))
		if err != nil {
			return err
		}
		result = append(result, page...)
		if len(page) != listOpts.PerPage || len(result) >= obj.Spec.Limit {
			break
		}
		listOpts.Page += 1
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	newMergeRequests := make([]runtime.RawExtension, 0, len(result))

	for _, mr := range result {
		if !sourceBranchRegex.MatchString(mr.SourceBranch) || !targetBrachRegex.MatchString(mr.TargetBranch) {
			continue
		}
		allLabelsFound := true
		for _, l := range obj.Spec.Labels {
			found := false
			for _, l2 := range mr.Labels {
				if l == l2 {
					found = true
					break
				}
			}
			if !found {
				allLabelsFound = false
				break
			}
		}
		if !allLabelsFound {
			continue
		}
		j, err := json.Marshal(mr)
		if err != nil {
			return err
		}
		newMergeRequests = append(newMergeRequests, runtime.RawExtension{Raw: j})
	}

	obj.Status.MergeRequests = newMergeRequests

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ListGitlabMergeRequestsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&templatesv1alpha1.ListGitlabMergeRequests{}).
		Complete(r)
}
