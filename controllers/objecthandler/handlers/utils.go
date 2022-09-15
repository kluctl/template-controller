package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kluctl/template-controller/api/v1alpha1"
	"github.com/kluctl/template-controller/controllers"
	"github.com/kluctl/template-controller/controllers/objecthandler/webgit"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

func objectToMap(obj any) (map[string]any, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var ret map[string]any
	err = json.Unmarshal(b, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func mapToObject(m any, out any) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, out)
	if err != nil {
		return err
	}
	return nil
}

func buildWebgitMergeRequest(ctx context.Context, client client.Client, namespace string, spec any, defaults *v1alpha1.ObjectHandlerDefaultsSpec) (webgit.MergeRequestInterface, error) {
	var err error
	var mergedMap map[string]any

	if defaults == nil {
		mergedMap, err = objectToMap(spec)
		if err != nil {
			return nil, err
		}
	} else {
		mergedMap, err = objectToMap(defaults)
		if err != nil {
			return nil, err
		}
		m2, err := objectToMap(spec)
		if err != nil {
			return nil, err
		}
		controllers.MergeMap2(mergedMap, m2, true)
	}

	if m, ok := mergedMap["gitlab"]; ok {
		var gitlab v1alpha1.GitlabMergeRequest
		err = mapToObject(m, &gitlab)
		if err != nil {
			return nil, err
		}
		if gitlab.MergeRequestId == nil {
			return nil, fmt.Errorf("missing mergeRequestId")
		}
		project, err := webgit.BuildWebgitGitlab(ctx, client, namespace, gitlab.GitlabProject)
		if err != nil {
			return nil, err
		}
		return project.GetMergeRequest(fmt.Sprintf("%d", *gitlab.MergeRequestId))
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
