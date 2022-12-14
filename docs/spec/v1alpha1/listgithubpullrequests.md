<!-- This comment is uncommented when auto-synced to www-kluctl.io

---
title: ListGithubPullRequests
description: ListGithubPullRequests documentation
weight: 30
---
-->

# ListGithubPullRequests

The `ListGithubPullRequests` API defines allows to query the GitHub API for a list of pull requests (PRs). These PRs
can be filtered when needed. The resulting list of PRs is written into the status of the
`ListGithubPullRequests` object.

The resulting PRs list inside the status can for example be used in `ObjectTemplate` to create objects based on
pull requests.

## Example

```yaml
apiVersion: templates.kluctl.io/v1alpha1
kind: ListGithubPullRequests
metadata:
  name: list-gh-prs
  namespace: default
spec:
  interval: 1m
  owner: podtato-head
  repo: podtato-head
  state: open
  base: main
  head: podtato-head:.*
  tokenRef:
    secretName: git-credentials
    key: github-token
```

The above example will regularly (1m interval) query the GitHub API for PRs inside the podtato-head
repository. It will filter for open PRs and for PRs against the main branch.

## Spec fields

### interval

Specifies the interval in which to query the GitHub API. Defaults to `5m`.

### owner

Specifies the user or organisation name where the repository is localed.

### repo

Specifies the repository name to query PRs for.

### tokenRef

In case of private repositories, this field can be used to specify a secret that contains a GitHub API token.

### head

Specifies the head to filter PRs for. The format must be `user:ref-name` / `organization:ref-name`. The `head`
field can also contain regular expressions.

### base

Specifies the base branch to filter PRs for. The `base` field can also contain regular expressions.

### labels

Specifies a list of labels to filter PRs for.

### state

Specifies the PR state to filter for. Can either be `open`, `closed` or `all`. Default to `all`.

### limit

Limits the number of results to accept. This is a safeguard for repositories with hundreds/thousands of PRs. It defaults
to 100.

## Resulting status

The query result is written into the `status.pullRequests` field of the `ListGithubPullRequests` object. Each entry
represents a reduced version of the [GitHub Pulls API](https://docs.github.com/en/rest/pulls/pulls#list-pull-requests)
results. The result is reduced in verbosity to avoid overloading the Kubernetes apiserver. Reduction means that all
fields containing `user`, `repo`, `orga` and `label` fields are reduced to `id`, `name`, `login`, `owner` and
`full_name`.

Please note that the resulting PR objects do not follow the typical camel case notion found in CRDs, as these represent
a copy of GitHub API objects.

Example:

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