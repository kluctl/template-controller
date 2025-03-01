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

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/kluctl/template-controller/controllers"
	"github.com/kluctl/template-controller/controllers/comments"
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	templatesv1alpha1 "github.com/kluctl/template-controller/api/v1alpha1"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(apiextensionsv1.AddToScheme(scheme))

	utilruntime.Must(templatesv1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var testJinja2 bool
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var watchAllNamespaces bool
	var concurrent int
	flag.BoolVar(&testJinja2, "test-jinja2", false, "Perform a simple Jinja2 rendering test and exit.")
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&watchAllNamespaces, "watch-all-namespaces", true,
		"Watch for custom resources in all namespaces, if set to false it will only watch the runtime namespace.")
	flag.IntVar(&concurrent, "concurrent", 4, "The number of concurrent reconciliations for each type.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	if testJinja2 {
		j2, err := controllers.NewJinja2()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to create jinja2 renderer: %s", err.Error())
			os.Exit(1)
		}
		defer j2.Close()
		_, err = j2.RenderString(`{{ "ok" }}`)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed render test string: %s", err.Error())
			os.Exit(1)
		}
		return
	}

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	ctx := ctrl.SetupSignalHandler()

	watchNamespace := ""
	if !watchAllNamespaces {
		watchNamespace = os.Getenv("RUNTIME_NAMESPACE")
	}

	var cacheNamespaces map[string]cache.Config
	if watchNamespace != "" {
		cacheNamespaces = map[string]cache.Config{
			watchNamespace: {},
		}
	}

	cfg := ctrl.GetConfigOrDie()
	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress: metricsAddr,
		},
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "3ab68de8.kluctl.io",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
		Cache: cache.Options{
			DefaultNamespaces: cacheNamespaces,
		},
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	rawWatchContext, rawWatchContextCancel := context.WithCancel(ctx)
	defer rawWatchContextCancel()

	fieldManager := "template-controller"

	if err = (&controllers.ObjectTemplateReconciler{
		BaseTemplateReconciler: controllers.BaseTemplateReconciler{
			Client:          mgr.GetClient(),
			RawWatchContext: rawWatchContext,
			Scheme:          mgr.GetScheme(),
			FieldManager:    fieldManager,
		},
	}).SetupWithManager(mgr, concurrent); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ObjectTemplate")
		os.Exit(1)
	}
	if err = (&controllers.TextTemplateReconciler{
		BaseTemplateReconciler: controllers.BaseTemplateReconciler{
			Client:          mgr.GetClient(),
			RawWatchContext: rawWatchContext,
			Scheme:          mgr.GetScheme(),
			FieldManager:    fieldManager,
		},
	}).SetupWithManager(mgr, concurrent); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "TextTemplate")
		os.Exit(1)
	}
	if err = (&controllers.ListGitlabMergeRequestsReconciler{
		Client:       mgr.GetClient(),
		Scheme:       mgr.GetScheme(),
		FieldManager: fieldManager,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ListGitlabMergeRequests")
		os.Exit(1)
	}
	if err = (&controllers.ListGithubPullRequestsReconciler{
		Client:       mgr.GetClient(),
		Scheme:       mgr.GetScheme(),
		FieldManager: fieldManager,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ListGithubPullRequests")
		os.Exit(1)
	}
	if err = (&controllers.GitProjectorReconciler{
		Client:       mgr.GetClient(),
		Scheme:       mgr.GetScheme(),
		FieldManager: fieldManager,
		TmpBaseDir:   filepath.Join(os.TempDir(), "template-controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "GitProjector")
		os.Exit(1)
	}
	if err = (&comments.GitlabCommentReconciler{
		BaseCommentReconciler: comments.BaseCommentReconciler{
			Client:       mgr.GetClient(),
			Scheme:       mgr.GetScheme(),
			FieldManager: fieldManager,
		},
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "GitlabComment")
		os.Exit(1)
	}
	if err = (&comments.GithubCommentReconciler{
		BaseCommentReconciler: comments.BaseCommentReconciler{
			Client:       mgr.GetClient(),
			Scheme:       mgr.GetScheme(),
			FieldManager: fieldManager,
		},
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "GithubComment")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
