<!-- This comment is uncommented when auto-synced to www-kluctl.io

---
title: Installation
description: Installation documentation
weight: 10
---
-->

# Installation

The Template Controller can currently only be installed via kustomize:

```sh
kubectl create ns kluctl-system
kustomize build "https://github.com/kluctl/template-controller/config/install?ref=v0.5.1" | kubectl apply -f-
```

Helm Charts will be supported in the near future.