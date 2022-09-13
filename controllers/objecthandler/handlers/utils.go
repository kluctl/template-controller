package handlers

import (
	"context"
	"fmt"
	"github.com/kluctl/template-controller/api/v1alpha1"
	"github.com/kluctl/template-controller/controllers/objecthandler/webgit"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

func buildWebgitMergeRequest(ctx context.Context, client client.Client, namespace string, spec any) (webgit.MergeRequestInterface, error) {
	v := reflect.Indirect(reflect.ValueOf(spec))
	gitlabValue := v.FieldByName("Gitlab")

	if gitlabValue.IsValid() && !gitlabValue.IsNil() {
		gitlabI := gitlabValue.Interface()
		gitlab, ok := gitlabI.(*v1alpha1.GitlabMergeRequest)
		if !ok {
			return nil, fmt.Errorf("invalid spec")
		}
		project, err := webgit.BuildWebgitGitlab(ctx, client, namespace, gitlab.GitlabProject)
		if err != nil {
			return nil, err
		}
		return project.GetMergeRequest(fmt.Sprintf("%d", gitlab.MergeRequestId))
	} else {
		return nil, fmt.Errorf("no pullRequest spec provided")
	}
}

func getClusterId(ctx context.Context, client client.Client) (string, error) {
	var ns v1.Namespace
	err := client.Get(ctx, types.NamespacedName{Name: "kube-system"}, &ns)
	if err != nil {
		return "", err
	}
	return string(ns.UID), nil
}

func generateMarkerComment(tag string, clusterId string, objNamespace string, objName string) string {
	return fmt.Sprintf("<!-- status-controller-%s \"%s\" \"%s/%s\" -->", tag, clusterId, objNamespace, objName)
}

func hasMarkerComment(body string, tag string, clusterId string, objNamespace string, objName string) bool {
	expected := generateMarkerComment(tag, clusterId, objNamespace, objName)
	for _, line := range strings.Split(body, "\n") {
		if line == expected {
			return true
		}
	}
	return false
}
