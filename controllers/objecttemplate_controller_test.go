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
	"math/rand/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	templatesv1alpha1 "github.com/kluctl/template-controller/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("ObjectTemplate controller", func() {
	const (
		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("Template without permissions to write object", func() {
		ns := fmt.Sprintf("test-%d", rand.Int64())

		key := client.ObjectKey{Name: "t1", Namespace: ns}
		cmKey := client.ObjectKey{Name: "cm1", Namespace: ns}

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
			}, "2s", interval).Should(MatchError("configmaps \"cm1\" not found"))

			assertFailedConfigMaps(key, cmKey)
		})
		It("Should succeed when RBAC is added", func() {
			createRoleWithBinding("default", ns, []string{"configmaps"})

			triggerReconcile(key)
			waitUntiReconciled(key, timeout)

			assertAppliedConfigMaps(key, cmKey)
			assertConfigMapData(cmKey, map[string]string{
				"k1": "3",
			})
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
			assertConfigMapData(cmKey, map[string]string{
				"k1": "3",
			})
		})
		It("Should cleanup", func() {
			Expect(k8sClient.Delete(ctx, t)).To(Succeed())
			waitUntilDeleted(key, &templatesv1alpha1.ObjectTemplate{}, timeout)
		})
	})
	Context("Template without permissions to read matrix object", func() {
		ns := fmt.Sprintf("test-%d", rand.Int64())

		key := client.ObjectKey{Name: "t1", Namespace: ns}
		cmKey := client.ObjectKey{Name: "cm1", Namespace: ns}

		t := buildObjectTemplate(key.Name, key.Namespace,
			[]templatesv1alpha1.MatrixEntry{
				buildMatrixObjectEntry("m1", "m1", ns, "Secret", "", false),
			},
			[]templatesv1alpha1.Template{
				{Object: buildTestConfigMap(cmKey.Name, cmKey.Namespace, map[string]string{
					"k1": `{{ matrix.m1.data.k1 + matrix.m1.data.k2 }}`,
				})},
			})

		It("Should fail initially", func() {
			createNamespace(ns)
			createServiceAccount("default", ns)
			createRoleWithBinding("default", ns, []string{"configmaps"})

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
			Expect(c.Message).To(ContainSubstring("Secret \"m1\" is forbidden"))
		})
		It("Should succeed when RBAC is created", func() {
			createRoleWithBinding("default", ns, []string{"secrets"})
			triggerReconcile(key)
			waitUntiReconciled(key, timeout)
			assertAppliedConfigMaps(key, cmKey)
			assertConfigMapData(cmKey, map[string]string{
				"k1": "MQ==Mg==", // two base64 encoded strings got added
			})
		})
		It("Should fail with non-existing SA", func() {
			updateObjectTemplate(key, func(t *templatesv1alpha1.ObjectTemplate) {
				t.Spec.ServiceAccountName = "non-existent"
			})
			waitUntiReconciled(key, timeout)

			t2 := getObjectTemplate(key)
			c := getReadyCondition(t2.GetConditions())
			Expect(c.Status).To(Equal(metav1.ConditionFalse))
			Expect(c.Message).To(ContainSubstring("Secret \"m1\" is forbidden"))
		})
		It("Should succeed after the SA is being created", func() {
			createServiceAccount("non-existent", ns)
			createRoleWithBinding("non-existent", ns, []string{"configmaps", "secrets"})
			triggerReconcile(key)
			waitUntiReconciled(key, timeout)
			t2 := getObjectTemplate(key)
			c := getReadyCondition(t2.GetConditions())
			Expect(c.Status).To(Equal(metav1.ConditionTrue))
		})
		It("Should still succeed if the matrix input has no namespace", func() {
			Expect(k8sClient.Create(ctx, buildTestSecret("m2", ns, map[string]string{
				"k3": "3",
			}))).To(Succeed())
			updateObjectTemplate(key, func(t *templatesv1alpha1.ObjectTemplate) {
				// we test the existing matrix entry to be changed
				t.Spec.Matrix[0].Object.Ref.Namespace = ""
				// and a new one being added
				t.Spec.Matrix = append(t.Spec.Matrix, buildMatrixObjectEntry("m2", "m2", "", "Secret", "", false))
				t.Spec.Templates[0].Object = buildTestConfigMap(cmKey.Name, cmKey.Namespace, map[string]string{
					"k1": `{{ matrix.m1.data.k1 + matrix.m1.data.k2 + matrix.m2.data.k3 }}`,
				})
			})
			waitUntiReconciled(key, timeout)
			t2 := getObjectTemplate(key)
			c := getReadyCondition(t2.GetConditions())
			Expect(c.Status).To(Equal(metav1.ConditionTrue))
			assertConfigMapData(cmKey, map[string]string{
				"k1": "MQ==Mg==Mw==", // three base64 encoded strings got added
			})
		})
		It("Should cleanup", func() {
			Expect(k8sClient.Delete(ctx, t)).To(Succeed())
			waitUntilDeleted(key, &templatesv1alpha1.ObjectTemplate{}, timeout)
		})
	})
	Context("Things get modified", func() {
		ns := fmt.Sprintf("test-%d", rand.Int64())

		key := client.ObjectKey{Name: "t1", Namespace: ns}
		cmKey := client.ObjectKey{Name: "cm1", Namespace: ns}
		sKey := client.ObjectKey{Name: "m1", Namespace: ns}

		t := buildObjectTemplate(key.Name, key.Namespace,
			[]templatesv1alpha1.MatrixEntry{
				buildMatrixObjectEntry("m1", "m1", ns, "Secret", "", false),
			},
			[]templatesv1alpha1.Template{
				{Object: buildTestConfigMap(cmKey.Name, cmKey.Namespace, map[string]string{
					"k1": `{{ matrix.m1.data.k1 + matrix.m1.data.k2 }}`,
				})},
			})
		// disable periodic reconciliation
		t.Spec.Interval.Duration = 100 * time.Hour

		It("Should succeed initially", func() {
			createNamespace(ns)
			createServiceAccount("default", ns)
			createRoleWithBinding("default", ns, []string{"configmaps"})
			createRoleWithBinding("default", ns, []string{"secrets"})

			err := k8sClient.Create(ctx, buildTestSecret("m1", ns, map[string]string{
				"k1": "1",
				"k2": "2",
			}))
			Expect(err).To(Succeed())

			Expect(k8sClient.Create(ctx, t)).Should(Succeed())
			waitUntiReconciled(key, timeout)
			assertAppliedConfigMaps(key, cmKey)

			assertConfigMapData(cmKey, map[string]string{
				"k1": "MQ==Mg==", // two base64 encoded strings got added
			})
		})
		It("Should update the CM after the matrix object got updated", func() {
			assertConfigMapData(cmKey, map[string]string{
				"k1": "MQ==Mg==", // two base64 encoded strings got added
			})

			var s v1.Secret
			Expect(k8sClient.Get(ctx, sKey, &s)).To(Succeed())
			s.StringData = map[string]string{
				"k1": "3",
			}
			Expect(k8sClient.Update(ctx, &s)).To(Succeed())

			Eventually(func() map[string]string {
				var cm v1.ConfigMap
				Expect(k8sClient.Get(ctx, cmKey, &cm)).To(Succeed())
				return cm.Data
			}, timeout, interval).Should(Equal(map[string]string{
				"k1": "Mw==Mg==", // two base64 encoded strings got added
			}))
		})
		It("Should update the CM after the ObjectTemplate got updated", func() {
			var t2 templatesv1alpha1.ObjectTemplate
			Expect(k8sClient.Get(ctx, key, &t2)).To(Succeed())
			t2.Spec.Templates[0].Object = buildTestConfigMap(cmKey.Name, cmKey.Namespace, map[string]string{
				"k1": `{{ (matrix.m1.data.k1 | b64decode | int) + (matrix.m1.data.k2 | b64decode | int) }}`,
			})
			Expect(k8sClient.Update(ctx, &t2)).To(Succeed())
			Eventually(func() map[string]string {
				var cm v1.ConfigMap
				Expect(k8sClient.Get(ctx, cmKey, &cm)).To(Succeed())
				return cm.Data
			}, timeout, interval).Should(Equal(map[string]string{
				"k1": "5",
			}))
		})
		It("Should cleanup", func() {
			Expect(k8sClient.Delete(ctx, t)).To(Succeed())
			waitUntilDeleted(key, &templatesv1alpha1.ObjectTemplate{}, timeout)
		})
	})
	Context("Pruning", func() {
		ns := fmt.Sprintf("test-%d", rand.Int64())

		key1 := client.ObjectKey{Name: "t1", Namespace: ns}
		key2 := client.ObjectKey{Name: "t2", Namespace: ns}
		cmKey1 := client.ObjectKey{Name: "cm1", Namespace: ns}
		cmKey2 := client.ObjectKey{Name: "cm2", Namespace: ns}
		cmKey3 := client.ObjectKey{Name: "cm3", Namespace: ns}
		cmKey4 := client.ObjectKey{Name: "cm4", Namespace: ns}

		assertAllExist := func() {
			assertConfigMapData(cmKey1, map[string]string{
				"k1": "3",
			})
			assertConfigMapData(cmKey2, map[string]string{
				"k1": "4",
			})
			assertConfigMapData(cmKey3, map[string]string{
				"k1": "5",
			})
			assertConfigMapData(cmKey4, map[string]string{
				"k1": "6",
			})
		}

		t1 := buildObjectTemplate(key1.Name, key1.Namespace,
			[]templatesv1alpha1.MatrixEntry{buildMatrixListEntry("m1")},
			[]templatesv1alpha1.Template{
				{Object: buildTestConfigMap(cmKey1.Name, cmKey1.Namespace, map[string]string{
					"k1": `{{ matrix.m1.k1 + matrix.m1.k2 }}`,
				})},
				{Object: buildTestConfigMap(cmKey2.Name, cmKey2.Namespace, map[string]string{
					"k1": `{{ matrix.m1.k1 + matrix.m1.k2 + 1 }}`,
				})},
			})
		t2 := buildObjectTemplate(key2.Name, key2.Namespace,
			[]templatesv1alpha1.MatrixEntry{buildMatrixListEntry("m1")},
			[]templatesv1alpha1.Template{
				{Object: buildTestConfigMap(cmKey3.Name, cmKey3.Namespace, map[string]string{
					"k1": `{{ matrix.m1.k1 + matrix.m1.k2 + 2}}`,
				})},
				{Object: buildTestConfigMap(cmKey4.Name, cmKey4.Namespace, map[string]string{
					"k1": `{{ matrix.m1.k1 + matrix.m1.k2 + 3 }}`,
				})},
			})
		t2.Spec.Prune = true

		It("Should create the objects initially", func() {
			createNamespace(ns)
			createServiceAccount("default", ns)
			createRoleWithBinding("default", ns, []string{"configmaps"})

			Expect(k8sClient.Create(ctx, t1)).Should(Succeed())
			Expect(k8sClient.Create(ctx, t2)).Should(Succeed())
			waitUntiReconciled(key1, timeout)
			waitUntiReconciled(key2, timeout)

			assertAppliedConfigMaps(key1, cmKey1, cmKey2)
			assertAppliedConfigMaps(key2, cmKey3, cmKey4)
			assertAllExist()
		})
		It("Should not delete the objects when prune is false", func() {
			Expect(k8sClient.Delete(ctx, t1)).To(Succeed())
			waitUntilDeleted(key1, &templatesv1alpha1.ObjectTemplate{}, timeout)
			assertAllExist()
		})
		It("Should delete the objects when prune is true", func() {
			Expect(k8sClient.Delete(ctx, t2)).To(Succeed())
			waitUntilDeleted(key2, &templatesv1alpha1.ObjectTemplate{}, timeout)
			waitUntilDeleted(cmKey3, &v1.ConfigMap{}, timeout)
			waitUntilDeleted(cmKey4, &v1.ConfigMap{}, timeout)
		})
	})
})
