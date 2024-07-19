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
kubectl apply -f "https://raw.githubusercontent.com/kluctl/template-controller/v0.9.2/deploy/manifests/template-controller.yaml"
```

## Helm
A Helm Chart for the controller is available as well.
To install the controller via Helm, run:
```shell
$ helm install template-controller -n template-controller --create-namespace oci://ghcr.io/kluctl/charts/template-controller
```

The Helm Chart is only distributed as an OCI package. The old Helm Repository found at https://github.com/kluctl/charts
is not maintained anymore.

## Upgrading from older Helm Charts

In case you were using the Helm Chart found at https://github.com/kluctl/charts, you'll need to perform a few extra
steps before you can upgrade to the new OCI based Helm Charts.

Run the following commands while the correct Kubectl Context is set. Please replace `<release-name>` with the release
name and `<release-namespace>` with the namespace you used when installing the old Chart.

```shell
$ rn=<release-name>
$ ns=<release-namespace>
$ for i in $(kubectl get crd -oname | grep templates.kluctl.io); do kubectl label $i app.kubernetes.io/managed-by=Helm; done
$ for i in $(kubectl get crd -oname | grep templates.kluctl.io); do kubectl annotate $i meta.helm.sh/release-name=$rn; done
$ for i in $(kubectl get crd -oname | grep templates.kluctl.io); do kubectl annotate $i meta.helm.sh/release-namespace=$ns; done
```

After this, you can perform a normal upgrade using the new OCI Chart.

```shell
$ helm upgrade -n <release-namespace> <release-name> oci://ghcr.io/kluctl/charts/template-controller
```