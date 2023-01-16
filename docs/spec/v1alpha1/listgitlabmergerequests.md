<!-- This comment is uncommented when auto-synced to www-kluctl.io

---
title: ListGitlabMergeRequests
linkTitle: ListGitlabMergeRequests
description: ListGitlabMergeRequests documentation
weight: 40
---
-->

# ListGitlabMergeRequests

The `ListGitlabMergeRequests` API allows to query the Gitlab API for a list of merge requests (MRs). These MRs
can be filtered when needed. The resulting list of MRs is written into the status of the
`ListGitlabMergeRequests` object.

The resulting MRs list inside the status can for example be used in `ObjectTemplate` to create objects based on
pull requests.

## Example

```yaml
apiVersion: templates.kluctl.io/v1alpha1
kind: ListGitlabMergeRequests
metadata:
  name: list-gl-mrs
  namespace: default
spec:
  interval: 1m
  project: my-group/my-repo
  state: opened
  targetBranch: main
  sourceBranch: prefix-.*
  tokenRef:
    secretName: git-credentials
    key: gitlab-token
```

The above example will regularly (1m interval) query the Gitlab API for MRs inside the `my-group/my-repo`
project. It will filter for open MRs and for MRs against the main branch.

## Spec fields

### interval

Specifies the interval in which to query the GitHub API. Defaults to `5m`.

### project

Specifies the Gitlab project to query MRs for. Must be in the format `group/project`, where group can also contain
subgroups (e.g. `group1/group2/project`).

### tokenRef

In case of private repositories, this field can be used to specify a secret that contains a Gitlab API token.

### targetBranch

Specifies the target branch to filter MRs for. The `targetBranch` field can also contain regular expressions.

### sourceBranch

Specifies the source branch to filter MRs for. The `sourceBranch` field can also contain regular expressions.

### labels

Specifies a list of labels to filter MRs for.

### state

Specifies the PR state to filter for. Can either be `opened`, `closed`, `locked`, `merged` or `all`. Default to `all`.

### limit

Limits the number of results to accept. This is a safeguard for repositories with hundreds/thousands of MRs. It defaults
to 100.

## Resulting status

The query result is written into the `status.mergeRequests` field of the `ListGitlabMergeRequests` object. The list is
identical to what is documented in the Gitlab [Merge requests API](https://docs.gitlab.com/ee/api/merge_requests.html).

Please note that the resulting PR objects do not follow the typical camel case notion found in CRDs, as these represent
a copy of Gitlab API objects.

Example:

```yaml
apiVersion: templates.kluctl.io/v1alpha1
kind: ListGitlabMergeRequests
metadata:
  name: list-gl-mrs
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
  mergeRequests:
  - id: 1
    iid: 1
    project_id: 3
    title: test1
    description: fixed login page css paddings
    state: merged
    merged_by:
      id: 87854
      name: Douwe Maan
      username: DouweM
      state: active
      avatar_url: 'https://gitlab.example.com/uploads/-/system/user/avatar/87854/avatar.png'
      web_url: 'https://gitlab.com/DouweM'
    merge_user:
      id: 87854
      name: Douwe Maan
      username: DouweM
      state: active
      avatar_url: 'https://gitlab.example.com/uploads/-/system/user/avatar/87854/avatar.png'
      web_url: 'https://gitlab.com/DouweM'
    merged_at: '2018-09-07T11:16:17.520Z'
    closed_by: null
    closed_at: null
    created_at: '2017-04-29T08:46:00Z'
    updated_at: '2017-04-29T08:46:00Z'
    target_branch: master
    source_branch: test1
    upvotes: 0
    downvotes: 0
    author:
      id: 1
      name: Administrator
      username: admin
      state: active
      avatar_url: null
      web_url: 'https://gitlab.example.com/admin'
    assignee:
      id: 1
      name: Administrator
      username: admin
      state: active
      avatar_url: null
      web_url: 'https://gitlab.example.com/admin'
    assignees:
      - name: Miss Monserrate Beier
        username: axel.block
        id: 12
        state: active
        avatar_url: >-
          http://www.gravatar.com/avatar/46f6f7dc858ada7be1853f7fb96e81da?s=80&d=identicon
        web_url: 'https://gitlab.example.com/axel.block'
    reviewers:
      - id: 2
        name: Sam Bauch
        username: kenyatta_oconnell
        state: active
        avatar_url: >-
          https://www.gravatar.com/avatar/956c92487c6f6f7616b536927e22c9a0?s=80&d=identicon
        web_url: 'http://gitlab.example.com//kenyatta_oconnell'
    source_project_id: 2
    target_project_id: 3
    labels:
      - Community contribution
      - Manage
    draft: false
    work_in_progress: false
    milestone:
      id: 5
      iid: 1
      project_id: 3
      title: v2.0
      description: Assumenda aut placeat expedita exercitationem labore sunt enim earum.
      state: closed
      created_at: '2015-02-02T19:49:26.013Z'
      updated_at: '2015-02-02T19:49:26.013Z'
      due_date: '2018-09-22'
      start_date: '2018-08-08'
      web_url: 'https://gitlab.example.com/my-group/my-project/milestones/1'
    merge_when_pipeline_succeeds: true
    merge_status: can_be_merged
    detailed_merge_status: not_open
    sha: '8888888888888888888888888888888888888888'
    merge_commit_sha: null
    squash_commit_sha: null
    user_notes_count: 1
    discussion_locked: null
    should_remove_source_branch: true
    force_remove_source_branch: false
    allow_collaboration: false
    allow_maintainer_to_push: false
    web_url: 'http://gitlab.example.com/my-group/my-project/merge_requests/1'
    references:
      short: '!1'
      relative: my-group/my-project!1
      full: my-group/my-project!1
    time_stats:
      time_estimate: 0
      total_time_spent: 0
      human_time_estimate: null
      human_total_time_spent: null
    squash: false
    task_completion_status:
      count: 0
      completed_count: 0
```