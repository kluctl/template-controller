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
kubectl apply -f "https://raw.githubusercontent.com/kluctl/template-controller/v0.9.1/deploy/manifests/template-controller.yaml"
```

## Helm
A Helm Chart for the controller is available as well.
To install the controller via Helm, run:
```shell
$ helm install template-controller -n template-controller --create-namespace oci://ghcr.io/kluctl/template-controller
```

The Helm Chart is only distributed as an OCI package. The old Helm Repository found at https://github.com/kluctl/charts
is not maintained anymore.
