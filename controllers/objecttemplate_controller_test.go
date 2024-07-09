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
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"math/rand/v2"
	client2 "sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	templatesv1alpha1 "github.com/kluctl/template-controller/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func updateObjectTemplate(key client2.ObjectKey, fn func(t *templatesv1alpha1.ObjectTemplate)) {
	t := getObjectTemplate(key)
	fn(t)
	err := k8sClient.Update(ctx, t, client2.FieldOwner("tests"))
	Expect(err).To(Succeed())
}

func triggerReconcile(key client2.ObjectKey) {
	updateObjectTemplate(key, func(t *templatesv1alpha1.ObjectTemplate) {
		t.Spec.Interval.Duration += time.Millisecond
	})
}

func waitUntiReconciled(key client2.ObjectKey, timeout time.Duration) {
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

func getObjectTemplate(key client2.ObjectKey) *templatesv1alpha1.ObjectTemplate {
	var t templatesv1alpha1.ObjectTemplate
	err := k8sClient.Get(ctx, key, &t)
	if err != nil {
		return nil
	}
	return &t
}

func assertAppliedConfigMaps(key client2.ObjectKey, keys ...client2.ObjectKey) {
	t := getObjectTemplate(key)
	Expect(t).ToNot(BeNil())

	var found []client2.ObjectKey
	for _, as := range t.Status.AppliedResources {
		if as.Success {
			found = append(found, client2.ObjectKey{Name: as.Ref.Name, Namespace: as.Ref.Namespace})
		}
	}

	Expect(found).To(ConsistOf(keys))
}

func assertFailedConfigMaps(key client2.ObjectKey, keys ...client2.ObjectKey) {
	t := getObjectTemplate(key)
	Expect(t).ToNot(BeNil())

	var found []client2.ObjectKey
	for _, as := range t.Status.AppliedResources {
		if !as.Success {
			found = append(found, client2.ObjectKey{Name: as.Ref.Name, Namespace: as.Ref.Namespace})
		}
	}

	Expect(found).To(ConsistOf(keys))
}

func assertFailedConfigMap(key client2.ObjectKey, cmKey client2.ObjectKey, errStr string) {
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

var _ = Describe("ObjectTemplate controller", func() {
	const (
		timeout  = time.Second * 1000
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("Template without permissions to write object", func() {
		ns := fmt.Sprintf("test-%d", rand.Int64())

		key := client2.ObjectKey{Name: "t1", Namespace: ns}
		cmKey := client2.ObjectKey{Name: "cm1", Namespace: ns}

		t := buildObjectTemplate(key.Name, key.Namespace,
			[]templatesv1alpha1.MatrixEntry{buildMatrixListEntry("m1")},
			[]templatesv1alpha1.Template{
				{Object: buildTestConfigMap(cmKey.Name, cmKey.Namespace, map[string]string{
					"k1": `{{ matrix.m1.k1 + matrix.m1.k2 }}`,
				})},
			})

		It("Should fail initially", func() {
			createNamespace(ns)
			createServiceAccount("default", ns)

			Expect(k8sClient.Create(ctx, t)).Should(Succeed())
			waitUntiReconciled(key, timeout)

			Consistently(func() error {
				return k8sClient.Get(ctx, cmKey, &v1.ConfigMap{})
			}, "2s").Should(MatchError("configmaps \"cm1\" not found"))

			assertFailedConfigMaps(key, cmKey)
		})
		It("Should succeed when RBAC is added", func() {
			createRoleWithBinding("default", ns, []string{"configmaps"})

			triggerReconcile(key)
			waitUntiReconciled(key, timeout)

			assertAppliedConfigMaps(key, cmKey)

			var cm v1.ConfigMap
			err := k8sClient.Get(ctx, cmKey, &cm)
			Expect(err).To(Succeed())

			Expect(cm.Data).To(Equal(map[string]string{
				"k1": "3",
			}))
		})
		It("Should fail with non-existing SA", func() {
			updateObjectTemplate(key, func(t *templatesv1alpha1.ObjectTemplate) {
				t.Spec.ServiceAccountName = "non-existent"
			})
			waitUntiReconciled(key, timeout)

			assertFailedConfigMap(key, cmKey, "configmaps \"cm1\" is forbidden")
		})
		It("Should succeed after the SA is being created", func() {
			createServiceAccount("non-existent", ns)
			createRoleWithBinding("non-existent", ns, []string{"configmaps"})
			triggerReconcile(key)
			waitUntiReconciled(key, timeout)
			assertAppliedConfigMaps(key, cmKey)
		})
	})
	Context("Template without permissions to read matrix object", func() {
		ns := fmt.Sprintf("test-%d", rand.Int64())

		key := client2.ObjectKey{Name: "t1", Namespace: ns}
		cmKey := client2.ObjectKey{Name: "cm1", Namespace: ns}

		It("Should fail initially", func() {
			createNamespace(ns)
			createServiceAccount("default", ns)
			createRoleWithBinding("default", ns, []string{"configmaps"})

			t := buildObjectTemplate(key.Name, key.Namespace,
				[]templatesv1alpha1.MatrixEntry{
					buildMatrixObjectEntry("m1", "m1", ns, "Secret", "", false),
				},
				[]templatesv1alpha1.Template{
					{Object: buildTestConfigMap(cmKey.Name, cmKey.Namespace, map[string]string{
						"k1": `{{ matrix.m1.k1 + matrix.m1.k2 }}`,
					})},
				})

			err := k8sClient.Create(ctx, buildTestSecret("m1", ns, map[string]string{
				"k1": "1",
				"k2": "2",
			}))
			Expect(err).To(Succeed())

			Expect(k8sClient.Create(ctx, t)).Should(Succeed())
			waitUntiReconciled(key, timeout)

			t2 := getObjectTemplate(key)
			c := getReadyCondition(t2.GetConditions())
			Expect(c.Status).To(Equal(metav1.ConditionFalse))
			Expect(c.Message).To(ContainSubstring("secrets \"m1\" is forbidden"))
		})
		It("Should succeed when RBAC is created", func() {
			createRoleWithBinding("default", ns, []string{"secrets"})
			triggerReconcile(key)
			waitUntiReconciled(key, timeout)
			assertAppliedConfigMaps(key, cmKey)
		})
		It("Should fail with non-existing SA", func() {
			updateObjectTemplate(key, func(t *templatesv1alpha1.ObjectTemplate) {
				t.Spec.ServiceAccountName = "non-existent"
			})
			waitUntiReconciled(key, timeout)

			t2 := getObjectTemplate(key)
			c := getReadyCondition(t2.GetConditions())
			Expect(c.Status).To(Equal(metav1.ConditionFalse))
			Expect(c.Message).To(ContainSubstring("secrets \"m1\" is forbidden"))
		})
		It("Should succeed after the SA is being created", func() {
			createServiceAccount("non-existent", ns)
			createRoleWithBinding("non-existent", ns, []string{"configmaps"})
			triggerReconcile(key)
			waitUntiReconciled(key, timeout)
			assertAppliedConfigMaps(key, cmKey)
		})
	})
})
