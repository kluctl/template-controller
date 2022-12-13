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
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/kluctl/kluctl/v2/pkg/git/messages"
	"io"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/fluxcd/go-git/v5/plumbing/object"
	"github.com/gobwas/glob"
	"github.com/kluctl/kluctl/v2/pkg/git"
	"github.com/kluctl/kluctl/v2/pkg/git/auth"
	git_url "github.com/kluctl/kluctl/v2/pkg/git/git-url"
	ssh_pool "github.com/kluctl/kluctl/v2/pkg/git/ssh-pool"
	templatesv1alpha1 "github.com/kluctl/template-controller/api/v1alpha1"
	yaml3 "gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// GitProjectorReconciler reconciles a GitProjector object
type GitProjectorReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	FieldManager string
	TmpBaseDir   string

	sshPool ssh_pool.SshPool
}

//+kubebuilder:rbac:groups=templates.kluctl.io,resources=gitprojectors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=templates.kluctl.io,resources=gitprojectors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=templates.kluctl.io,resources=gitprojectors/finalizers,verbs=update

func (r *GitProjectorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	logger := log.FromContext(ctx)

	logger.V(1).Info("Starting reconcile")
	defer logger.V(1).Info("Finished reconcile", "err", err)

	var obj templatesv1alpha1.GitProjector
	err = r.Get(ctx, req.NamespacedName, &obj)
	if err != nil {
		logger.Error(err, "Get failed")
		err = client.IgnoreNotFound(err)
		return
	}

	// Add our finalizer if it does not exist
	if !controllerutil.ContainsFinalizer(&obj, templatesv1alpha1.ObjectTemplateFinalizer) {
		patch := client.MergeFrom(obj.DeepCopy())
		controllerutil.AddFinalizer(&obj, templatesv1alpha1.ObjectTemplateFinalizer)
		if err := r.Patch(ctx, &obj, patch, client.FieldOwner(r.FieldManager)); err != nil {
			logger.Error(err, "unable to register finalizer")
			return ctrl.Result{}, err
		}
	}

	// Examine if the object is under deletion
	if !obj.GetDeletionTimestamp().IsZero() {
		return r.finalize(ctx, &obj)
	}

	// Return early if the object is suspended.
	if obj.Spec.Suspend {
		logger.Info("Reconciliation is suspended for this object")
		return ctrl.Result{}, nil
	}

	patch := client.MergeFrom(obj.DeepCopy())
	err = r.doReconcile(ctx, &obj)
	if err != nil {
		c := metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionFalse,
			ObservedGeneration: obj.GetGeneration(),
			Reason:             "Error",
			Message:            err.Error(),
		}
		apimeta.SetStatusCondition(&obj.Status.Conditions, c)
	} else {
		c := metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionTrue,
			ObservedGeneration: obj.GetGeneration(),
			Reason:             "Success",
			Message:            "Success",
		}
		apimeta.SetStatusCondition(&obj.Status.Conditions, c)
	}
	err = r.Status().Patch(ctx, &obj, patch, client.FieldOwner(r.FieldManager))
	if err != nil {
		return
	}

	result.RequeueAfter = obj.Spec.Interval.Duration
	return
}

// SetupWithManager sets up the controller with the Manager.
func (r *GitProjectorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&templatesv1alpha1.GitProjector{}).
		Complete(r)
}

func (r *GitProjectorReconciler) doReconcile(ctx context.Context, obj *templatesv1alpha1.GitProjector) error {
	url, err := git_url.Parse(obj.Spec.RepoUrl)
	if err != nil {
		return err
	}

	auth, err := r.buildGitAuth(ctx, obj)
	if err != nil {
		return err
	}

	mr, err := git.NewMirroredGitRepo(ctx, *url, filepath.Join(r.TmpBaseDir, "git-mirrors"), &r.sshPool, auth)
	if err != nil {
		return err
	}

	err = mr.Lock()
	if err != nil {
		return err
	}
	defer mr.Unlock()

	err = mr.Update()
	if err != nil {
		return err
	}

	matchingRefs, err := r.filterRefs(obj, mr)
	if err != nil {
		return err
	}
	sortedMatchedRefNames := make([]string, 0, len(matchingRefs))
	for name, _ := range matchingRefs {
		sortedMatchedRefNames = append(sortedMatchedRefNames, name)
	}
	sort.Strings(sortedMatchedRefNames)

	allRefsHash := sha256.New()
	for _, name := range sortedMatchedRefNames {
		_, _ = fmt.Fprintf(allRefsHash, "%s=%s\n", name, matchingRefs[name])
	}
	allRefsHashStr := hex.EncodeToString(allRefsHash.Sum(nil))
	if allRefsHashStr == obj.Status.AllRefsHash {
		// nothing to do
		return nil
	}

	var globs []glob.Glob
	for _, f := range obj.Spec.Files {
		g, err := glob.Compile(f.Glob, '/')
		if err != nil {
			return err
		}
		globs = append(globs, g)
	}

	type matchedFile struct {
		gitFile templatesv1alpha1.GitFile
		file    object.File
	}

	newResults := make([]templatesv1alpha1.GitProjectorResult, 0, len(matchingRefs))
	for name, hash := range matchingRefs {
		t, err := mr.GetGitTreeByCommit(hash)
		if err != nil {
			return err
		}

		var matchedFiles []matchedFile

		err = t.Files().ForEach(func(file *object.File) error {
			for i, g := range globs {
				if g.Match(file.Name) {
					matchedFiles = append(matchedFiles, matchedFile{
						gitFile: obj.Spec.Files[i],
						file:    *file,
					})
					break
				}
			}
			return nil
		})
		if err != nil {
			return err
		}

		var ref templatesv1alpha1.GitRef
		if strings.HasPrefix(name, "refs/heads/") {
			ref.Branch = strings.TrimPrefix(name, "refs/heads/")
		} else if strings.HasSuffix(name, "refs/tags/") {
			ref.Tag = strings.TrimPrefix(name, "refs/tags/")
		} else {
			return fmt.Errorf("could not determine ref type for %s", name)
		}

		result := templatesv1alpha1.GitProjectorResult{
			Reference: ref,
		}

		for _, mf := range matchedFiles {
			rawContent, err := mf.file.Contents()
			if err != nil {
				return err
			}
			resultFile := templatesv1alpha1.GitProjectorResultFile{
				Path: mf.file.Name,
			}
			if mf.gitFile.ParseYaml {
				d := yaml3.NewDecoder(strings.NewReader(rawContent))
				for {
					var a any
					err = d.Decode(&a)
					if err == io.EOF {
						break
					}
					if err != nil {
						return fmt.Errorf("failed to parse %s as yaml", mf.file.Name)
					}
					b, err := json.Marshal(a)
					if err != nil {
						return fmt.Errorf("failed to marshal %s as json", mf.file.Name)
					}
					resultFile.Parsed = append(resultFile.Parsed, &runtime.RawExtension{Raw: b})
				}
			} else {
				resultFile.Raw = &rawContent
			}
			result.Files = append(result.Files, resultFile)
		}
		newResults = append(newResults, result)
	}

	sort.Slice(newResults, func(i, j int) bool {
		return newResults[i].Reference.Less(newResults[j].Reference)
	})

	obj.Status.Result = newResults
	obj.Status.AllRefsHash = allRefsHashStr

	return nil
}

func (r *GitProjectorReconciler) filterRefs(obj *templatesv1alpha1.GitProjector, mr *git.MirroredGitRepo) (map[string]string, error) {
	refs, err := mr.RemoteRefHashesMap()
	if err != nil {
		return nil, err
	}

	matchingRefs := map[string]string{}

	if obj.Spec.Reference == nil {
		defaultRef, err := mr.DefaultRef()
		if err != nil {
			return nil, err
		}
		hash, ok := refs[defaultRef]
		if !ok {
			return nil, fmt.Errorf("default ref %s not found", defaultRef)
		}
		matchingRefs[defaultRef] = hash
		return matchingRefs, nil
	}

	if obj.Spec.Reference.Commit != "" {
		found := false
		for name, hash := range refs {
			if hash == obj.Spec.Reference.Commit {
				matchingRefs[name] = hash
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("commit %s not found", obj.Spec.Reference.Commit)
		}
		return matchingRefs, nil
	}

	var regex *regexp.Regexp
	if obj.Spec.Reference.Tag != "" {
		regex, err = regexp.Compile(fmt.Sprintf("^refs/tags/%s$", obj.Spec.Reference.Tag))
		if err != nil {
			return nil, fmt.Errorf("invalid tag regex specified: %w", err)
		}
	} else if obj.Spec.Reference.Branch != "" {
		regex, err = regexp.Compile(fmt.Sprintf("^refs/heads/%s$", obj.Spec.Reference.Branch))
		if err != nil {
			return nil, fmt.Errorf("invalid branch regex specified: %w", err)
		}
	} else {
		return nil, fmt.Errorf("ref is empty")
	}

	for name, hash := range refs {
		if regex.MatchString(name) {
			matchingRefs[name] = hash
		}
	}

	return matchingRefs, nil
}

func (r *GitProjectorReconciler) buildGitAuth(ctx context.Context, obj *templatesv1alpha1.GitProjector) (*auth.GitAuthProviders, error) {
	logger := log.FromContext(ctx)

	ga := auth.NewDefaultAuthProviders("GIT", &messages.MessageCallbacks{
		WarningFn: func(s string) {
			logger.Info(s)
		},
		TraceFn: func(s string) {
			logger.V(1).Info(s)
		},
	})

	if obj.Spec.SecretRef == nil {
		return ga, nil
	}

	var gitSecret corev1.Secret
	err := r.Client.Get(ctx, types.NamespacedName{Namespace: obj.Namespace, Name: obj.Spec.SecretRef.Name}, &gitSecret)
	if err != nil {
		return nil, err
	}

	e := auth.AuthEntry{
		Host:     "*",
		Username: "*",
	}

	if x, ok := gitSecret.Data["username"]; ok {
		e.Username = string(x)
	}
	if x, ok := gitSecret.Data["password"]; ok {
		e.Password = string(x)
	}
	if x, ok := gitSecret.Data["caFile"]; ok {
		e.CABundle = x
	}
	if x, ok := gitSecret.Data["known_hosts"]; ok {
		e.KnownHosts = x
	}
	if x, ok := gitSecret.Data["identity"]; ok {
		e.SshKey = x
	}

	var la auth.ListAuthProvider
	la.AddEntry(e)
	ga.RegisterAuthProvider(&la, false)
	return ga, nil
}

func (r *GitProjectorReconciler) finalize(ctx context.Context, obj *templatesv1alpha1.GitProjector) (ctrl.Result, error) {
	r.doFinalize(ctx, obj)

	// Remove our finalizer from the list and update it
	controllerutil.RemoveFinalizer(obj, templatesv1alpha1.ObjectTemplateFinalizer)
	if err := r.Update(ctx, obj, client.FieldOwner(r.FieldManager)); err != nil {
		return ctrl.Result{}, err
	}

	// Stop reconciliation as the object is being deleted
	return ctrl.Result{}, nil
}

func (r *GitProjectorReconciler) doFinalize(ctx context.Context, obj *templatesv1alpha1.GitProjector) {
}
