package controllers

import (
	templatesv1alpha1 "github.com/kluctl/template-controller/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	. "github.com/onsi/gomega"
)

func toUnstructured(o client.Object) *unstructured.Unstructured {
	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(o)
	Expect(err).To(Succeed())
	return &unstructured.Unstructured{
		Object: m,
	}
}

func buildTestConfigMap(name string, namespace string, data map[string]string) *unstructured.Unstructured {
	return toUnstructured(&v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	})
}

func buildTestSecret(name string, namespace string, data map[string]string) *unstructured.Unstructured {
	return toUnstructured(&v1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		StringData: data,
	})
}

func buildObjectTemplate(name string, namespace string, matrixEntries []templatesv1alpha1.MatrixEntry, templates []templatesv1alpha1.Template) *templatesv1alpha1.ObjectTemplate {
	t := &templatesv1alpha1.ObjectTemplate{
		TypeMeta: metav1.TypeMeta{
			APIVersion: templatesv1alpha1.GroupVersion.String(),
			Kind:       "ObjectTemplate",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: templatesv1alpha1.ObjectTemplateSpec{
			Interval:  metav1.Duration{Duration: time.Second},
			Matrix:    matrixEntries,
			Templates: templates,
		},
	}
	return t
}

func buildMatrixListEntry(name string) templatesv1alpha1.MatrixEntry {
	return templatesv1alpha1.MatrixEntry{
		Name: name,
		List: []runtime.RawExtension{
			{
				Raw: []byte(`{"k1": 1, "k2": 2}`),
			},
		},
	}
}

func buildMatrixObjectEntry(name string, objName string, objNamespace string, objKind string, jsonPath string, expandLists bool) templatesv1alpha1.MatrixEntry {
	var jsonPathPtr *string
	if jsonPath != "" {
		jsonPathPtr = &jsonPath
	}
	return templatesv1alpha1.MatrixEntry{
		Name: name,
		Object: &templatesv1alpha1.MatrixEntryObject{
			Ref: templatesv1alpha1.ObjectRef{
				APIVersion: "v1",
				Kind:       objKind,
				Name:       objName,
				Namespace:  objNamespace,
			},
			JsonPath:    jsonPathPtr,
			ExpandLists: expandLists,
		},
	}
}

func updateObjectTemplate(key client.ObjectKey, fn func(t *templatesv1alpha1.ObjectTemplate)) {
	t := getObjectTemplate(key)
	fn(t)
	err := k8sClient.Update(ctx, t, client.FieldOwner("tests"))
	Expect(err).To(Succeed())
}

func triggerReconcile(key client.ObjectKey) {
	updateObjectTemplate(key, func(t *templatesv1alpha1.ObjectTemplate) {
		t.Spec.Interval.Duration += time.Millisecond
	})
}

func waitUntiReconciled(key client.ObjectKey, timeout time.Duration) {
	Eventually(func() bool {
		t := getObjectTemplate(key)
		if t == nil {
			return false
		}
		c := getReadyCondition(t.GetConditions())
		if c == nil {
			return false
		}
		return c.ObservedGeneration == t.Generation
	}, timeout, time.Millisecond*250).Should(BeTrue())
}

func getObjectTemplate(key client.ObjectKey) *templatesv1alpha1.ObjectTemplate {
	var t templatesv1alpha1.ObjectTemplate
	err := k8sClient.Get(ctx, key, &t)
	if err != nil {
		return nil
	}
	return &t
}

func assertAppliedConfigMaps(key client.ObjectKey, keys ...client.ObjectKey) {
	t := getObjectTemplate(key)
	Expect(t).ToNot(BeNil())

	var found []client.ObjectKey
	for _, as := range t.Status.AppliedResources {
		if as.Success {
			found = append(found, client.ObjectKey{Name: as.Ref.Name, Namespace: as.Ref.Namespace})
		}
	}

	Expect(found).To(ConsistOf(keys))
}

func assertFailedConfigMaps(key client.ObjectKey, keys ...client.ObjectKey) {
	t := getObjectTemplate(key)
	Expect(t).ToNot(BeNil())

	var found []client.ObjectKey
	for _, as := range t.Status.AppliedResources {
		if !as.Success {
			found = append(found, client.ObjectKey{Name: as.Ref.Name, Namespace: as.Ref.Namespace})
		}
	}

	Expect(found).To(ConsistOf(keys))
}

func assertFailedConfigMap(key client.ObjectKey, cmKey client.ObjectKey, errStr string) {
	t := getObjectTemplate(key)
	Expect(t).ToNot(BeNil())
	for _, as := range t.Status.AppliedResources {
		if !as.Success && as.Ref.Name == cmKey.Name && as.Ref.Namespace == cmKey.Namespace {
			Expect(as.Error).To(ContainSubstring(errStr))
			return
		}
	}
	Expect(false).To(BeTrue())
}
