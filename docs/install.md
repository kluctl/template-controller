# Installation

The Template Controller can currently only be installed via kustomize:

```sh
kubectl create ns kluctl-system
kustomize build "https://github.com/kluctl/template-controller/config/install?ref=v0.4.1" | kubectl apply -f-
```

Helm Charts will be supported in the near future.