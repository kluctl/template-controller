<!-- This comment is uncommented when auto-synced to www-kluctl.io

---
title: GitProjector
linkTitle: GitProjector
description: GitProjector documentation
weight: 20
---
-->

# GitProjector

The `GitProjector` API defines projections of Git repositories.

Projection of Git repositories means that the content of selected branches and selected files are loaded into Kubernetes,
accessible through the status of the `GitProjector`.

The projected branches and files can then be used as matrix inputs for an `ObjectTemplate`.

## Example

```yaml
apiVersion: templates.kluctl.io/v1alpha1
kind: GitProjector
metadata:
  name: preview
  namespace: default
spec:
  interval: 1m
  url: https://github.com/kluctl/kluctl-examples.git
  # In case you use a private repository
  secretRef:
    name: git-credentials
  ref:
    branch: main
  files:
    - glob: "preview-envs/preview-*.yaml"
      parseYaml: true
```

The above example creates a `GitProjector` that will periodically clone the kluctl-examples repo, look for the `main`
branch and all files matching the given glob. It will then parse all yamls and make them available through the
`GitProjector`'s status:

```yaml
apiVersion: templates.kluctl.io/v1alpha1
kind: GitProjector
metadata:
  name: preview
  namespace: default
spec:
  ...
status:
  allRefsHash: 104d3dc9b5ffabf5ba3c76532fb71da58757c494acdcb7dff3665d256f516612
  conditions:
  - lastTransitionTime: "2022-12-14T09:09:51Z"
    message: Success
    observedGeneration: 1
    reason: Success
    status: "True"
    type: Ready
  result:
  - files:
    - parsed:
      - envName: preview-env1
        replicas: 3
      path: preview-envs/preview-env1.yaml
    - parsed:
      - envName: preview-env2
        replicas: 1
      path: preview-envs/preview-env2.yaml
    ref:
      branch: main
```

## Spec fields

The following fields are supported in `spec`.

### interval

Specifies the interval at which the `GitProjector` is reconciled.

### suspend

If set to `true`, reconciliation is suspended.

### url

The git url of the repository to project. Can either be a https or a git/ssh url.

### ref

The git reference to project. Either `spec.ref.branch` or `spec.ref.tag` must be set.

Both tags and refs can be regular expressions. In case of a regular expression, the controller will include all matching
refs in the `status.result` field.

### secretRef

Same as in the Kluctl Controllers [KluctlDeployment](https://kluctl.io/docs/flux/spec/v1alpha1/kluctldeployment/#git-authentication)

### files

List of file to project into the status. Must be of the format:

```yaml
...
spec:
  ...
  files:
    - glob: "my-file.yaml"
      parseYaml: true
```

Each entry must at least contain a `glob` which is used to match files. The controller uses the https://github.com/gobwas/glob
library for pattern matching.

If `parseYaml` is set to `true`, the controller will try to parse matching files as yaml and include the parsed structured
data in the resulting status. Parsing of yaml is done with the assumption that all files possibly contain multiple yaml
documents, meaning that even yaml files with just a single document will result in a parsed list of one document.

Consider the following matching yaml file:

```yaml
envName: preview-env1
replicas: 3
```

This will result in the following projection:

```yaml
...
status:
  result:
  - files:
    - parsed:
      - envName: preview-env1
        replicas: 3
      path: preview-envs/preview-env1.yaml
    ref:
      branch: main
```

If `parseYaml` is `false`, the result will contain a raw string representation of the matching files:

```yaml
...
status:
  result:
  - files:
    - path: preview-envs/preview-env1.yaml
      raw: |-
        envName: preview-env1
        replicas: 3
    ref:
      branch: main
```
