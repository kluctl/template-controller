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
		duration = time.Second * 10
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
			}, "2s").Should(MatchError("configmaps \"cm1\" not found"))

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
			Expect(c.Message).To(ContainSubstring("secrets \"m1\" is forbidden"))
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
			Expect(c.Message).To(ContainSubstring("secrets \"m1\" is forbidden"))
		})
		It("Should succeed after the SA is being created", func() {
			createServiceAccount("non-existent", ns)
			createRoleWithBinding("non-existent", ns, []string{"configmaps"})
			triggerReconcile(key)
			waitUntiReconciled(key, timeout)
			assertAppliedConfigMaps(key, cmKey)
		})
		It("Should cleanup", func() {
			Expect(k8sClient.Delete(ctx, t)).To(Succeed())
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
			}, timeout, time.Millisecond*250).Should(Equal(map[string]string{
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
			}, timeout, time.Millisecond*250).Should(Equal(map[string]string{
				"k1": "5",
			}))
		})
		It("Should cleanup", func() {
			Expect(k8sClient.Delete(ctx, t)).To(Succeed())
		})
	})
})
