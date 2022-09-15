package handlers

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

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
