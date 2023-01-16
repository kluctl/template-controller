<!-- This comment is uncommented when auto-synced to www-kluctl.io

---
title: TextTemplate
linkTitle: TextTemplate
description: TextTemplate documentation
weight: 30
---
-->

# GithubComment

The `TextTemplate` API allows to define text templates that are rendered into the status of the TextTemplate.
The result can for example be used in `GitlabComment`/`GithubComment`.

## Example

For the below example to work, you will also have to deploy the RBAC resources documented in
[ObjectTemplate](./objecttemplate.md#serviceaccountname).

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-configmap
  namespace: default
data:
  mykey: input-value
---
apiVersion: templates.kluctl.io/v1alpha1
kind: TextTemplate
metadata:
  name: example
  namespace: default
spec:
  serviceAccountName: example-template-service-account
  inputs:
    - name: input1
      object:
        ref:
          apiVersion: v1
          kind: ConfigMap
          name: my-configmap
  template: |
    This template text can use variables from the inputs defined above, for example this: {{ inputs.input1.data.mykey }}.
```

The above example will render the given template text and write it into the result of the object:

```yaml
apiVersion: templates.kluctl.io/v1alpha1
kind: TextTemplate
...
status:
  conditions:
  - lastTransitionTime: "2023-01-16T11:24:15Z"
    message: Success
    observedGeneration: 2
    reason: Success
    status: "True"
    type: Ready
  result: 'This template text can use variables from the inputs defined above, for example this: input-value.'
```

## Spec fields

### suspend

If set to `true`, reconciliation of this TextTemplate is suspended.

### serviceAccountName

The service account to use while retrieving template inputs. See the [ObjectTemplate](./objecttemplate.md#serviceaccountname)
documentation for details.

### inputs

List of template inputs which are then available while rendering the text template. At the moment, only Kubernetes
objects are supported as inputs, but other types of inputs might be supported in the future.

Example:

```yaml
apiVersion: templates.kluctl.io/v1alpha1
kind: TextTemplate
metadata:
  name: example
  namespace: default
spec:
  serviceAccountName: example-template-service-account
  inputs:
    - name: input1
      object:
        ref:
          apiVersion: v1
          kind: ConfigMap
          name: my-configmap
          namespace: default
        jsonPath: data
  template: |
    This template text can use variables from the inputs defined above, for example this: {{ inputs.input1.mykey }}.
```

#### inputs.name

Specifies the name of the input, which is then used to refer to the input inside the text template.

#### inputs.object

Specifies the object to load as input. The specified [service account](#serviceaccountname) must have proper permissions
to access this object.

### template

Specifies the raw template text to be rendered in the reconciliation loop. While rendering, each input is available
via the global `inputs` variable and the specified name of the input, e.g. `{{ inputs.my_input.sub_field }}.

See [templating](../../templating.md) for more details on the templating engine.

### templateRef

Specifies another object to load the template text from. Currently only ConfigMaps are supported.

#### templateRef.configMap:

Specifies a ConfigMap to load the template from.

Example:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-configmap
  namespace: default
data:
  mykey: input-value
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-template
  namespace: default
data:
  template: |
    This template text can use variables from the inputs defined above, for example this: {{ inputs.input1.data.mykey }}.
---
apiVersion: templates.kluctl.io/v1alpha1
kind: TextTemplate
metadata:
  name: example
  namespace: default
spec:
  serviceAccountName: example-template-service-account
  inputs:
    - name: input1
      object:
        ref:
          apiVersion: v1
          kind: ConfigMap
          name: my-configmap
  templateRef:
    configMap:
      name: my-template
      key: template
```

## Resulting status

The resulting rendered template is written into the status and can then be used by other objects, e.g. `GitlabComment`/`GithubComment`.

Example:

```yaml
...
status:
  conditions:
    - lastTransitionTime: "2023-01-16T11:24:15Z"
      message: Success
      observedGeneration: 3
      reason: Success
      status: "True"
      type: Ready
  result: 'This template text can use variables from the inputs defined above,
    for example this: input-value.'
```