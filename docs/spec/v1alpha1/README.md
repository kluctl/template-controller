<!-- This comment is uncommented when auto-synced to www-kluctl.io

---
title: v1alpha1 specs
linkTitle: v1alpha1 specs
description: templates.kluctl.io/v1alpha1 documentation
weight: 10
---
-->

# templates.kluctl.io/v1alpha1

This is the v1alpha1 API specification for defining templating related resources.

## Specification

- [ObjectTemplate CRD](objecttemplate.md)
    + [Spec fields](objecttemplate.md#spec-fields)
- [GitProjector CRD](gitprojector.md)
    + [Spec fields](gitprojector.md#spec-fields)
- [ListGithubPullRequests CRD](listgithubpullrequests.md)
    + [Spec fields](listgithubpullrequests.md#spec-fields)
- [ListGitlabMergeRequests CRD](listgitlabmergerequests.md)
    + [Spec fields](listgitlabmergerequests.md#spec-fields)
- [GithubComment CRD](githubcomment.md)
    + [Spec fields](githubcomment.md#spec-fields)
- [GitlabComment CRD](gitlabcomment.md)
    + [Spec fields](gitlabcomment.md#spec-fields)

## Implementation

* [template-controller](https://github.com/kluctl/template-controller)
