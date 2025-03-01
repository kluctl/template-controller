# template-controller

![Version: 0.9.4](https://img.shields.io/badge/Version-0.9.4-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: v0.9.4](https://img.shields.io/badge/AppVersion-v0.9.4-informational?style=flat-square)

A Helm chart for the template-controller

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| kluctl |  | <https://github.com/kluctl> |
| AljoschaP | <aljoscha.poertner@posteo.de> |  |
| codablock | <ablock84@gmail.com> |  |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` |  |
| crds.annotations | object | `{}` |  |
| env | list | `[]` |  |
| fullnameOverride | string | `""` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| image.repository | string | `"ghcr.io/kluctl/template-controller"` |  |
| image.tag | string | `""` |  |
| imagePullSecrets | list | `[]` |  |
| installCRDs | bool | `true` | If set, install and upgrade CRDs through helm chart. |
| nameOverride | string | `""` |  |
| nodeSelector | object | `{}` |  |
| podAnnotations."prometheus.io/port" | string | `"8080"` |  |
| podAnnotations."prometheus.io/scrape" | string | `"true"` |  |
| podSecurityContext.fsGroup | int | `1337` |  |
| resources | object | `{}` |  |
| securityContext.capabilities.drop[0] | string | `"ALL"` |  |
| securityContext.readOnlyRootFilesystem | bool | `true` |  |
| securityContext.runAsNonRoot | bool | `true` |  |
| securityContext.runAsUser | int | `1337` |  |
| securityContext.seccompProfile.type | string | `"RuntimeDefault"` |  |
| service.health.port | int | `8081` |  |
| service.prometheus.port | int | `8080` |  |
| service.type | string | `"ClusterIP"` |  |
| serviceAccount.annotations | object | `{}` |  |
| serviceAccount.create | bool | `true` |  |
| serviceAccount.name | string | `""` |  |
| tolerations | list | `[]` |  |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.14.2](https://github.com/norwoodj/helm-docs/releases/v1.14.2)
