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
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"math/rand/v2"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	templatesv1alpha1 "github.com/kluctl/template-controller/api/v1alpha1"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	cfg       *rest.Config
	k8sClient client.Client // You'll be using this client in your tests.
	testEnv   *envtest.Environment
	ctx       context.Context
	cancel    context.CancelFunc
)

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	ctx, cancel = context.WithCancel(context.TODO())

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	var err error
	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = templatesv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	controllerCfg := setupTemplateControllerRBAC()

	k8sManager, err := ctrl.NewManager(controllerCfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	err = (&ObjectTemplateReconciler{
		BaseTemplateReconciler: BaseTemplateReconciler{
			Client:          k8sManager.GetClient(),
			RawWatchContext: ctx,
			Scheme:          k8sManager.GetScheme(),
			FieldManager:    "template-controller",
		},
	}).SetupWithManager(k8sManager, 1)
	Expect(err).ToNot(HaveOccurred())

	err = (&TextTemplateReconciler{
		BaseTemplateReconciler: BaseTemplateReconciler{
			Client:          k8sManager.GetClient(),
			RawWatchContext: ctx,
			Scheme:          k8sManager.GetScheme(),
			FieldManager:    "template-controller",
		},
	}).SetupWithManager(k8sManager, 1)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred(), "failed to run manager")
	}()

})

var _ = AfterSuite(func() {
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

func setupTemplateControllerRBAC() *rest.Config {
	createNamespace("system")

	user, err := testEnv.AddUser(envtest.User{Name: "template-controller"}, nil)
	Expect(err).NotTo(HaveOccurred())

	rbacPath := filepath.Join("..", "config", "rbac")

	files := []string{"role.yaml"}
	applyFromDir(rbacPath, files, "system")

	rb := rbacv1.ClusterRoleBinding{
		ObjectMeta: v1.ObjectMeta{
			Name: "template-controller",
		},
		RoleRef: rbacv1.RoleRef{
			Kind: "ClusterRole",
			Name: "manager-role",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: "User",
				Name: "template-controller",
			},
		},
	}
	Expect(k8sClient.Create(ctx, &rb)).To(Succeed())

	return user.Config()
}

func createNamespace(name string) {
	ns := corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
	}
	err := k8sClient.Create(ctx, &ns)
	Expect(err).To(Succeed())
}

func createServiceAccount(saName string, saNamespace string) {
	sa := corev1.ServiceAccount{
		ObjectMeta: v1.ObjectMeta{
			Name:      saName,
			Namespace: saNamespace,
		},
	}
	err := k8sClient.Create(ctx, &sa)
	Expect(err).To(Succeed())
}

func createRoleWithBinding(saName string, saNamespace string, resources []string) {
	roleName := fmt.Sprintf("role-%s-%d", saName, rand.Int64())

	role := rbacv1.Role{
		ObjectMeta: v1.ObjectMeta{
			Name:      roleName,
			Namespace: saNamespace,
		},
		Rules: []rbacv1.PolicyRule{
			{
				Verbs:     []string{"*"},
				APIGroups: []string{"", "*"},
				Resources: resources,
			},
		},
	}
	roleBinding := rbacv1.RoleBinding{
		ObjectMeta: v1.ObjectMeta{
			Name:      roleName,
			Namespace: saNamespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      saName,
				Namespace: saNamespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind: "Role",
			Name: roleName,
		},
	}

	err := k8sClient.Create(ctx, &role)
	Expect(err).To(Succeed())
	err = k8sClient.Create(ctx, &roleBinding)
	Expect(err).To(Succeed())
}

func getReadyCondition(conditions []v1.Condition) *v1.Condition {
	for _, c := range conditions {
		if c.Type == "Ready" {
			return &c
		}
	}
	return nil
}
