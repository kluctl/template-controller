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
kustomize build "https://github.com/kluctl/template-controller/config/install?ref=v0.8.0" | kubectl apply -f-
```

## Helm
A Helm Chart for the controller is also available [here](https://github.com/kluctl/charts/tree/main/charts/template-controller).
To install the controller via Helm, run:
```shell
$ helm repo add kluctl https://kluctl.github.io/charts
$ helm install template-controller kluctl/template-controller
```
