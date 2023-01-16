<!-- This comment is uncommented when auto-synced to www-kluctl.io

---
title: GitlabComment
linkTitle: GitlabComment
description: GitlabComment documentation
weight: 30
---
-->

# GitlabComment

The `GitlabComment` API allows to post a comment to a Gitlab Merge Request.

## Example

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-configmap
  namespace: default
data:
  my-key: |
    This can by **any** form of [Markdown](https://en.wikipedia.org/wiki/Markdown) supported by Gitlab.
---
apiVersion: templates.kluctl.io/v1alpha1
kind: GitlabComment
metadata:
  name: comment-gl
  namespace: default
spec:
  gitlab:
    project: my-group/my-repo
    mergeRequestId: 1234
    tokenRef:
      secretName: git-credentials
      key: gitlab-token
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

### gitlab

Specifies which Gitlab project and merge request to post the comment to.

#### gitlab.project

Specifies the user or organisation name where the repository is localed.

#### gitlab.repo

Specifies the repository name to query PRs for.

#### gitlab.tokenRef

In case of private repositories, this field can be used to specify a secret that contains a Gitlab API token.

#### github.pullRequestId

Specifies the ID of the pull request.

### comment

Same as in [GithubComment](./githubcomment.md#comment)
