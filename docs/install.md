# Installation

The Template Controller can currently only be installed via kustomize:

```sh
kustomize build "https://github.com/kluctl/template-controller/config/install?ref=v0.3.0" | kubectl apply -f-
```

Helm Charts will be supported in the near future.