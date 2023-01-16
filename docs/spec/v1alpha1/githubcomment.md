<!-- This comment is uncommented when auto-synced to www-kluctl.io

---
title: GithubComment
linkTitle: GithubComment
description: GithubComment documentation
weight: 30
---
-->

# GithubComment

The `GithubComment` API allows to post a comment to a GitHub Pull Request.

## Example

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-configmap
  namespace: default
data:
  my-key: |
    This can by **any** form of [Markdown](https://en.wikipedia.org/wiki/Markdown) supported by Github.
---
apiVersion: templates.kluctl.io/v1alpha1
kind: GithubComment
metadata:
  name: comment-gh
  namespace: default
spec:
  github:
    owner: my-org-or-user
    repo: my-repo
    pullRequestId: 1234
    tokenRef:
      secretName: git-credentials
      key: github-token
  comment:
    source:
      configMap:
        name: my-configmap
        key: my-key
```

The above example will post a comment to the specified pull request. The comment's content is loaded from the ConfigMap
`my-configmap`. Other sources are also supported, see the `source` field documentation for details.

The comment will be updated whenever the underlying comment source changes.

## Spec fields

### suspend

If set to `true`, reconciliation of this comment is suspended.

### github

Specifies which GitHub project and pull request to post the comment to.

#### github.owner

Specifies the user or organisation name where the repository is localed.

#### github.repo

Specifies the repository name to query PRs for.

#### github.tokenRef

In case of private repositories, this field can be used to specify a secret that contains a GitHub API token.

#### github.pullRequestId

Specifies the ID of the pull request.

### comment

This field specifies the necessary information for the comment content.

#### comment.id

This optional field specifies the identifier to mark the comment with so that the controller can identify it. It
defaults to a generated id built from the namespace and name of the comment resource.

#### comment.source

This specifies the comment source. Multiple source types are supported, specified via a sub-field.

##### comment.source.text

Raw text for the template's content. Example:

```yaml
apiVersion: templates.kluctl.io/v1alpha1
kind: GithubComment
metadata:
  name: comment-gh
  namespace: default
spec:
  github:
    owner: my-org-or-user
    repo: my-repo
    pullRequestId: 1234
    tokenRef:
      secretName: git-credentials
      key: github-token
  comment:
    source:
      text: |
        This can by **any** form of [Markdown](https://en.wikipedia.org/wiki/Markdown) supported by Github.
```

##### comment.source.configMap

Uses a ConfigMap as source for the comment's content. Example:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-configmap
  namespace: default
data:
  my-key: |
    This can by **any** form of [Markdown](https://en.wikipedia.org/wiki/Markdown) supported by Github.
---
apiVersion: templates.kluctl.io/v1alpha1
kind: GithubComment
metadata:
  name: comment-gh
  namespace: default
spec:
  github:
    owner: my-org-or-user
    repo: my-repo
    pullRequestId: 1234
    tokenRef:
      secretName: git-credentials
      key: github-token
  comment:
    source:
      configMap:
        name: my-configmap
        key: my-key
```

##### comment.source.textTemplate

Uses a [TextTemplate](./texttemplate.md) as source for the comment's content. Example:

```yaml
apiVersion: templates.kluctl.io/v1alpha1
kind: TextTemplate
metadata:
  name: my-texttemplate
  namespace: default
spec:
  inputs:
    ... # See TextTemplate documentation for details.
  template: |
    This can by **any** form of [Markdown](https://en.wikipedia.org/wiki/Markdown) supported by Github.
---
apiVersion: templates.kluctl.io/v1alpha1
kind: GithubComment
metadata:
  name: comment-gh
  namespace: default
spec:
  github:
    owner: my-org-or-user
    repo: my-repo
    pullRequestId: 1234
    tokenRef:
      secretName: git-credentials
      key: github-token
  comment:
    source:
      textTemplate:
        name: my-texttemplate
```
