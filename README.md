<!-- This comment is uncommented when auto-synced to www-kluctl.io

---
title: "Template Controller"
linkTitle: "Template Controller"
description: "Template Controller documentation."
weight: 200
---
-->

# Template Controller

The Template Controller is a controller originating from the [Kluctl](https://kluctl.io) project, but not limited to
Kluctl. It allows to define template objects which are rendered and applied into the cluster based on an input matrix.

In its easiest form, an `ObjectTemplate` takes one input object (e.g. a ConfigMap) and creates another object
(e.g. a Secret) which is then applied into the cluster.

The Template Controller also offers CRDs which allow to query external resources (e.g. GitHub Pull Requests) which can
then be used as inputs into `ObjectTemplates`.

## Use Cases

Template Controller has many use case. Some are for example:
1. [Dynamic environments for Pull Requests](./docs/use-case-dynamic-environments.md)
2. [Transformation of Secrets/Objects](./docs/use-case-transformation.md)

## Documentation

Reference documentation is available [here](./docs/spec/v1alpha1).

## Installation

Installation instructions can be found [here](./docs/install.md)
