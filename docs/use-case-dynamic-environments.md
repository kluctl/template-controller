# Use Case: Dynamic environments for Pull Requests

This use case was the initial and first use case why the Template Controller was created. You can use `ListGithubPullRequests`
to query the GitHub API for a list of pull requests on a GitHub Repo and then use the result inside a `ObjectTemplate`
to generate GitOps environments for new pull requests.

## Flux

This example will create templated [Kustomization](https://fluxcd.io/flux/components/kustomize/kustomization/)
objects. The means, that you should first [install Flux](https://fluxcd.io/flux/installation/) on your cluster. The
[dev install](https://fluxcd.io/flux/installation/#dev-install) variant should be sufficient.

## podtato-head as example

This example uses the [podtato-head](https://github.com/podtato-head/podtato-head) demo project to demonstrate the
Template Controller. You must fork the repository and replace all occurrences of `podtato-head` as `owner` with your
own username. It is not recommended to blindly use the public repository as you this will lead to unverified and
potentially dangerous environments being deployed into your cluster!

## GitHub credentials

In case you want to listen for PRs from a private repository (e.g. because you've forked podtato-head), you'll need to store a
[GitHub personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)
inside a Kubernetes Secret.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: git-credentials
  namespace: default
stringData:
  github-token: "<your-github-token>"
```

WARNING: Of course, in a real setup you would NOT store the plain token inside a manifest, but instead use
[Sealed Secrets](https://github.com/bitnami-labs/sealed-secrets) or [SOPS](https://github.com/mozilla/sops).

## A dedicated ServiceAccount

The Template Controller uses service accounts to query matrix inputs and apply rendered objects. These service accounts
determine what the template can access and what not. In this example, we'll create a service account with the
`cluster-admin` role, which you should NOT do in production. Instead, define your own `Role` or `ClusterRole` and
attach it to the service account. This role should have read/write access to all objects references in the matrix and
the rendered objects.

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: podtato-head-envs-objecttemplate
  namespace: default
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: podtato-head-envs-objecttemplate
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  # WARNING, this is only for demo purposes. You should use a more restricted role for the ObjectTemplate
  name: cluster-admin
subjects:
  - kind: ServiceAccount
    name: podtato-head-envs-objecttemplate
    namespace: default
```

The above serviceAccount is then later referenced inside the `ObjectTemplate` object.

## Listing GitHub pull requests

Listing pull requests from a GitHub repository can be done through the
[`ListGithubPullRequests`](./spec/v1alpha1/listgithubpullrequests.md) CRD. It specifies the GitHub repository to use and
some filter options.

```yaml
apiVersion: templates.kluctl.io/v1alpha1
kind: ListGithubPullRequests
metadata:
  name: list-gh-prs
  namespace: default
spec:
  interval: 1m
  # Replace the owner with your username in case you forked podtato-head
  owner: podtato-head
  repo: podtato-head
  # Ignore closed PRs
  state: open
  # Only PR's that go against the main branch
  base: main
  # Replace `podtato-head` with your username. This will only allows heads from your own fork!
  # Otherwise, you risk deploying unsafe environments into your cluster!
  head: podtato-head:.*
  tokenRef:
    secretName: git-credentials
    key: github-token
```

After applying this resource, the Template Controller will start to query the GitHub API for matching pull requests and
then store the results inside the status of the `ListGithubPullRequests` CR. Example:

```yaml
apiVersion: templates.kluctl.io/v1alpha1
kind: ListGithubPullRequests
metadata:
  name: list-gh-prs
  namespace: default
spec:
  ...
status:
  conditions:
  - lastTransitionTime: "2022-11-07T14:55:36Z"
    message: Success
    observedGeneration: 3
    reason: Success
    status: "True"
    type: Ready
  # The pullRequests list contains much more detailed info, but to keep it short I've reduced verbosity here
  pullRequests:
  - base:
      label: podtato-head:main
      ref: main
      repo:
        full_name: podtato-head/podtato-head
        name: podtato-head
      sha: de7e66af16d41b0ef83de9a0b3be6f5cf0caf942
    body: "..."
    created_at: "2022-02-02T23:06:28Z"
    head:
      label: vivek:issue-79_implement_ms_ketch
      ref: issue-79_implement_ms_ketch
      repo:
        full_name: vivek/podtato-head
        name: podtato-head
      sha: 6379b4c8f413dae70daa03a5a13de4267486fd59
    number: 151
    state: open
    title: '...'
    updated_at: "2022-02-04T03:53:03Z"
```

## The ObjectTemplate

The `pullRequests` field from the above status can then be used as an input into the an
[`ObjectTemplate`](./spec/v1alpha1/objecttemplate.md).

```yaml
apiVersion: templates.kluctl.io/v1alpha1
kind: ObjectTemplate
metadata:
  name: pr-envs
  namespace: default
spec:
  serviceAccountName: podtato-head-envs-objecttemplate
  # This causes removal of templated objects in case they disappear from the rendered list of objects
  prune: true
  matrix:
    - name: pr
      object:
        ref:
          apiVersion: templates.kluctl.io/v1alpha1
          kind: ListGithubPullRequests
          name: list-gh-prs
        jsonPath: status.pullRequests
        expandLists: true
  templates:
    - object:
        apiVersion: v1
        kind: Namespace
        metadata:
          # Give each one its own namespace
          name: "podtato-head-{{ matrix.pr.head.label | slugify }}"
    - object:
        apiVersion: source.toolkit.fluxcd.io/v1beta2
        kind: GitRepository
        metadata:
          # The pullRequests status field from the ListGithubPullRequests is a reduced form of the REST API result
          # of https://docs.github.com/en/rest/pulls/pulls#list-pull-requests, meaning that fields like `head` and `base`
          # are also available.
          name: "podtato-head-{{ matrix.pr.head.label | slugify }}"
          namespace: default
        spec:
          interval: 5m
          url: "https://github.com/{{ matrix.pr.head.repo.full_name }}.git"
          ref:
            branch: "{{ matrix.pr.head.ref }}"
    - object:
        apiVersion: kustomize.toolkit.fluxcd.io/v1beta2
        kind: Kustomization
        metadata:
          name: "podtato-head-env-{{ matrix.pr.head.label | slugify }}"
          namespace: default
        spec:
          interval: 10m
          targetNamespace: "podtato-head-{{ matrix.pr.head.label | slugify }}"
          sourceRef:
            kind: GitRepository
            # refers to the same GitRepository created above
            name: "podtato-head-{{ matrix.pr.head.label | slugify }}"
          path: "./delivery/kustomize/base"
          prune: true
```

The above `ObjectTemplate` will create 3 objects per pull request:
1. A namespace with the name `podtato-head-{{ matrix.pr.head.label | slugify }}`. Please note the use of 
Jinja2 templating. Details about what can be done can be found in the
[`ObjectTemplate`](./spec/v1alpha1/objecttemplate.md) documentation.
2. A [Flux GitRepository](https://fluxcd.io/flux/components/source/gitrepositories/) that points to repository
and branch of the current pull request.
3. A [Flux Kustomization](https://fluxcd.io/flux/components/kustomize/kustomization/) that is deployed into
the above namespace.
