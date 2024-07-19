package controllers

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"math/rand/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"time"

	templatesv1alpha1 "github.com/kluctl/template-controller/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TextTemplate controller", func() {
	const (
		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("Template with no inputs and simple string template", func() {
		ns := fmt.Sprintf("test-%d", rand.Int64())

		key := client.ObjectKey{Name: "t1", Namespace: ns}

		t := buildTextTemplate(key.Name, key.Namespace,
			nil,
			"{{ 'v1' }}", nil, "",
		)

		It("Should succeed with string template and no inputs", func() {
			createNamespace(ns)

			Expect(k8sClient.Create(ctx, t)).Should(Succeed())
			waitUntiTextTemplateReconciled(key, timeout)
			assertTextTemplateResult(key, "v1")
		})
		It("Should cleanup", func() {
			Expect(k8sClient.Delete(ctx, t)).To(Succeed())
			waitUntilDeleted(key, &templatesv1alpha1.TextTemplate{}, timeout)
		})
	})
	Context("Template with input", func() {
		ns := fmt.Sprintf("test-%d", rand.Int64())

		key := client.ObjectKey{Name: "t1", Namespace: ns}

		t := buildTextTemplate(key.Name, key.Namespace,
			[]templatesv1alpha1.TextTemplateInput{buildTextTemplateInputEntry("i1", "i1", ns, "ConfigMap", "")},
			"{{ inputs.i1.data.k1 }}", nil, "",
		)

		It("Should fail with missing RBAC for input", func() {
			createNamespace(ns)
			createServiceAccount("default", ns)

			cm := buildTestConfigMap("i1", ns, map[string]string{
				"k1": "v1",
				"k2": "v2",
			})
			Expect(k8sClient.Create(ctx, cm)).To(Succeed())

			Expect(k8sClient.Create(ctx, t)).Should(Succeed())
			waitUntiTextTemplateReconciled(key, timeout)
			assertTextTemplateResult(key, "")

			t2 := getTextTemplate(key)
			c := getReadyCondition(t2.GetConditions())
			Expect(c.Status).To(Equal(metav1.ConditionFalse))
			Expect(c.Message).To(ContainSubstring("ConfigMap \"i1\" is forbidden"))
		})
		It("Should succeed when RBAC is added", func() {
			createRoleWithBinding("default", ns, []string{"configmaps"})

			triggerReconcileTextTemplate(key)
			waitUntiTextTemplateReconciled(key, timeout)
			t2 := getTextTemplate(key)
			c := getReadyCondition(t2.GetConditions())
			Expect(c.Status).To(Equal(metav1.ConditionTrue))
			assertTextTemplateResult(key, "v1")
		})
		It("Should fail with non-existing SA", func() {
			updateTextTemplate(key, func(t *templatesv1alpha1.TextTemplate) {
				t.Spec.ServiceAccountName = "non-existent"
			})
			waitUntiTextTemplateReconciled(key, timeout)
			t2 := getTextTemplate(key)
			c := getReadyCondition(t2.GetConditions())
			Expect(c.Status).To(Equal(metav1.ConditionFalse))
			Expect(c.Message).To(ContainSubstring("ConfigMap \"i1\" is forbidden"))
			assertTextTemplateResult(key, "v1") // should still be the old value
		})
		It("Should succeed after the SA is being created", func() {
			createServiceAccount("non-existent", ns)
			createRoleWithBinding("non-existent", ns, []string{"configmaps"})
			triggerReconcileTextTemplate(key)
			waitUntiTextTemplateReconciled(key, timeout)
			t2 := getTextTemplate(key)
			c := getReadyCondition(t2.GetConditions())
			Expect(c.Status).To(Equal(metav1.ConditionTrue))
			assertTextTemplateResult(key, "v1")
		})
		It("Should still succeed when the input object has no namespace", func() {
			updateTextTemplate(key, func(t *templatesv1alpha1.TextTemplate) {
				t.Spec.Inputs[0].Object.Ref.Namespace = ""
			})
			waitUntiTextTemplateReconciled(key, timeout)
			t2 := getTextTemplate(key)
			c := getReadyCondition(t2.GetConditions())
			Expect(c.Status).To(Equal(metav1.ConditionTrue))
		})
		It("Should respect the jsonPath", func() {
			updateTextTemplate(key, func(t *templatesv1alpha1.TextTemplate) {
				p := ".data"
				t.Spec.Inputs[0].Object.JsonPath = &p
			})
			waitUntiTextTemplateReconciled(key, timeout)
			t2 := getTextTemplate(key)
			c := getReadyCondition(t2.GetConditions())
			Expect(c.Status).To(Equal(metav1.ConditionFalse))
			Expect(c.Message).To(ContainSubstring("UndefinedError"))
			assertTextTemplateResult(key, "v1")
			updateTextTemplate(key, func(t *templatesv1alpha1.TextTemplate) {
				newTemplate := "{{ inputs.i1.k2 }}"
				t.Spec.Template = &newTemplate
			})
			waitUntiTextTemplateReconciled(key, timeout)
			t2 = getTextTemplate(key)
			c = getReadyCondition(t2.GetConditions())
			Expect(c.Status).To(Equal(metav1.ConditionTrue))
			assertTextTemplateResult(key, "v2")
		})
		It("Should cleanup", func() {
			Expect(k8sClient.Delete(ctx, t)).To(Succeed())
			waitUntilDeleted(key, &templatesv1alpha1.TextTemplate{}, timeout)
		})
	})
	Context("Template with external template", func() {
		ns := fmt.Sprintf("test-%d", rand.Int64())

		key := client.ObjectKey{Name: "t1", Namespace: ns}

		t := buildTextTemplate(key.Name, key.Namespace,
			nil,
			"", &client.ObjectKey{Name: "t1", Namespace: ns}, "t1",
		)

		It("Should fail without permissions to read the template", func() {
			createNamespace(ns)
			createServiceAccount("default", ns)

			cm1 := buildTestConfigMap("t1", ns, map[string]string{
				"t1": "{{ 't1' }}",
			})
			Expect(k8sClient.Create(ctx, cm1)).To(Succeed())

			Expect(k8sClient.Create(ctx, t)).Should(Succeed())
			waitUntiTextTemplateReconciled(key, timeout)
			assertTextTemplateResult(key, "")

			t2 := getTextTemplate(key)
			c := getReadyCondition(t2.GetConditions())
			Expect(c.Status).To(Equal(metav1.ConditionFalse))
			Expect(c.Message).To(ContainSubstring("ConfigMap \"t1\" is forbidden"))
		})
		It("Should succeed when RBAC is added", func() {
			createRoleWithBinding("default", ns, []string{"configmaps"})

			triggerReconcileTextTemplate(key)
			waitUntiTextTemplateReconciled(key, timeout)
			t2 := getTextTemplate(key)
			c := getReadyCondition(t2.GetConditions())
			Expect(c.Status).To(Equal(metav1.ConditionTrue))
			assertTextTemplateResult(key, "t1")
		})
		It("Should cleanup", func() {
			Expect(k8sClient.Delete(ctx, t)).To(Succeed())
			waitUntilDeleted(key, &templatesv1alpha1.TextTemplate{}, timeout)
		})
	})
	Context("Things get modified", func() {
		ns := fmt.Sprintf("test-%d", rand.Int64())

		key := client.ObjectKey{Name: "t1", Namespace: ns}
		inputKey := client.ObjectKey{Name: "i1", Namespace: ns}
		templateKey := client.ObjectKey{Name: "t1", Namespace: ns}

		t := buildTextTemplate(key.Name, key.Namespace,
			[]templatesv1alpha1.TextTemplateInput{buildTextTemplateInputEntry("i1", "i1", ns, "ConfigMap", "")},
			"{{ inputs.i1.data.k1 }}", nil, "",
		)

		It("Should succeed initially", func() {
			createNamespace(ns)
			createServiceAccount("default", ns)
			createRoleWithBinding("default", ns, []string{"configmaps"})

			cm := buildTestConfigMap("i1", ns, map[string]string{
				"k1": "v1",
				"k2": "v2",
			})
			Expect(k8sClient.Create(ctx, cm)).To(Succeed())

			Expect(k8sClient.Create(ctx, t)).Should(Succeed())
			waitUntiTextTemplateReconciled(key, timeout)
			assertTextTemplateResult(key, "v1")
		})
		It("Should properly update on input change", func() {
			var cm v1.ConfigMap
			Expect(k8sClient.Get(ctx, inputKey, &cm)).To(Succeed())
			cm.Data["k1"] = "v3"
			Expect(k8sClient.Update(ctx, &cm)).To(Succeed())
			Eventually(func() string {
				t2 := getTextTemplate(key)
				return t2.Status.Result
			}, timeout, interval).Should(Equal("v3"))
		})
		It("Should properly update on template change", func() {
			updateTextTemplate(key, func(t *templatesv1alpha1.TextTemplate) {
				tmp := "{{ inputs.i1.data.k2 }}"
				t.Spec.Template = &tmp
			})
			Eventually(func() string {
				t2 := getTextTemplate(key)
				return t2.Status.Result
			}, timeout, interval).Should(Equal("v2"))
		})
		It("Should fail initially when referring a non-existent external template", func() {
			updateTextTemplate(key, func(t *templatesv1alpha1.TextTemplate) {
				t.Spec.Template = nil
				t.Spec.TemplateRef = &templatesv1alpha1.TemplateRef{
					ConfigMap: &templatesv1alpha1.TemplateRefConfigMap{
						Name:      templateKey.Name,
						Namespace: ns,
						Key:       "t1",
					},
				}
			})
			waitUntiTextTemplateReconciled(key, timeout)
			t2 := getTextTemplate(key)
			c := getReadyCondition(t2.GetConditions())
			Expect(c.Status).To(Equal(metav1.ConditionFalse))
			Expect(c.Message).To(ContainSubstring("configmaps \"t1\" not found"))
			assertTextTemplateResult(key, "v2") // nothing changed
		})
		It("Should succeed when the external template gets created", func() {
			cm := buildTestConfigMap("t1", ns, map[string]string{
				"t1": "{{ 't1' }}",
				"t2": "{{ 't2' }}",
			})
			Expect(k8sClient.Create(ctx, cm)).To(Succeed())
			Eventually(func() string {
				t2 := getTextTemplate(key)
				return t2.Status.Result
			}, timeout, interval).Should(Equal("t1"))
			t2 := getTextTemplate(key)
			c := getReadyCondition(t2.GetConditions())
			Expect(c.Status).To(Equal(metav1.ConditionTrue))
		})
		It("Should update when the external template gets updated", func() {
			var cm v1.ConfigMap
			Expect(k8sClient.Get(ctx, templateKey, &cm)).To(Succeed())
			cm.Data["t1"] = "{{ 't3' }}"
			Expect(k8sClient.Update(ctx, &cm)).To(Succeed())
			Eventually(func() string {
				t2 := getTextTemplate(key)
				return t2.Status.Result
			}, timeout, interval).Should(Equal("t3"))
		})
		It("Should update when the external template key gets updated", func() {
			updateTextTemplate(key, func(t *templatesv1alpha1.TextTemplate) {
				t.Spec.TemplateRef.ConfigMap.Key = "t2"
			})
			Eventually(func() string {
				t2 := getTextTemplate(key)
				return t2.Status.Result
			}, timeout, interval).Should(Equal("t2"))
		})
		It("Should cleanup", func() {
			Expect(k8sClient.Delete(ctx, t)).To(Succeed())
			waitUntilDeleted(key, &templatesv1alpha1.TextTemplate{}, timeout)
		})
	})
})

func buildTextTemplate(name string, namespace string, inputs []templatesv1alpha1.TextTemplateInput, template string, templateConfigMapRef *client.ObjectKey, templateConfigMapKey string) *templatesv1alpha1.TextTemplate {
	var templatePtr *string
	if template != "" {
		templatePtr = &template
	}
	var ref *templatesv1alpha1.TemplateRef
	if templateConfigMapRef != nil {
		ref = &templatesv1alpha1.TemplateRef{
			ConfigMap: &templatesv1alpha1.TemplateRefConfigMap{
				Name:      templateConfigMapRef.Name,
				Namespace: templateConfigMapRef.Namespace,
				Key:       templateConfigMapKey,
			},
		}
	}
	t := &templatesv1alpha1.TextTemplate{
		TypeMeta: metav1.TypeMeta{
			APIVersion: templatesv1alpha1.GroupVersion.String(),
			Kind:       "TextTemplate",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: templatesv1alpha1.TextTemplateSpec{
			Inputs:      inputs,
			Template:    templatePtr,
			TemplateRef: ref,
		},
	}
	return t
}

func buildTextTemplateInputEntry(name string, objName string, objNamespace string, objKind string, jsonPath string) templatesv1alpha1.TextTemplateInput {
	var jsonPathPtr *string
	if jsonPath != "" {
		jsonPathPtr = &jsonPath
	}
	return templatesv1alpha1.TextTemplateInput{
		Name: name,
		Object: &templatesv1alpha1.TextTemplateInputObject{
			Ref: templatesv1alpha1.ObjectRef{
				APIVersion: "v1",
				Kind:       objKind,
				Name:       objName,
				Namespace:  objNamespace,
			},
			JsonPath: jsonPathPtr,
		},
	}
}

func updateTextTemplate(key client.ObjectKey, fn func(t *templatesv1alpha1.TextTemplate)) {
	t := getTextTemplate(key)
	fn(t)
	err := k8sClient.Update(ctx, t, client.FieldOwner("tests"))
	Expect(err).To(Succeed())
}

func triggerReconcileTextTemplate(key client.ObjectKey) {
	updateTextTemplate(key, func(t *templatesv1alpha1.TextTemplate) {
		a := t.GetAnnotations()
		if a == nil {
			a = map[string]string{}
		}
		if old, ok := a["test"]; ok {
			x, _ := strconv.Atoi(old)
			a["test"] = fmt.Sprintf("%d", x+1)
		} else {
			a["test"] = "1"
		}
		t.SetAnnotations(a)
	})
}

func waitUntiTextTemplateReconciled(key client.ObjectKey, timeout time.Duration) {
	Eventually(func() bool {
		t := getTextTemplate(key)
		if t == nil {
			return false
		}
		c := meta.FindStatusCondition(t.GetConditions(), "Test")
		if c == nil {
			return false
		}
		return t.Generation == c.ObservedGeneration && c.Message == t.GetAnnotations()["test"]
	}, timeout, time.Millisecond*250).Should(BeTrue())
}

func getTextTemplate(key client.ObjectKey) *templatesv1alpha1.TextTemplate {
	var t templatesv1alpha1.TextTemplate
	err := k8sClient.Get(ctx, key, &t)
	if err != nil {
		return nil
	}
	return &t
}

func assertTextTemplateResult(key client.ObjectKey, result string) {
	t := getTextTemplate(key)
	Expect(t.Status.Result).To(Equal(result))
}
