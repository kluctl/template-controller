<!-- This comment is uncommented when auto-synced to www-kluctl.io

---
title: Installation
description: Installation documentation
weight: 10
---
-->

# Installation

The Template Controller can currently be installed via static manifests or via Helm.

## Static Manifests
```sh
kubectl apply -f "https://raw.githubusercontent.com/kluctl/template-controller/v0.8.3/deploy/manifests/template-controller.yaml"
```

## Helm
A Helm Chart for the controller is available as well.
To install the controller via Helm, run:
```shell
$ helm repo add kluctl https://kluctl.github.io/charts
$ helm install template-controller -n template-controller --create-namespace kluctl/template-controller
```
