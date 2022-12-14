<!-- This comment is uncommented when auto-synced to www-kluctl.io

---
title: ObjectTemplate
description: ObjectTemplate documentation
weight: 10
---
-->

# ObjectTemplate

The `ObjectTemplate` API defines templates that are rendered based on a matrix of input values.

## Example

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: input-configmap
  namespace: default
data:
  x: someValue
---
apiVersion: templates.kluctl.io/v1alpha1
kind: ObjectTemplate
metadata:
  name: example-template
  namespace: default
spec:
  serviceAccountName: example-template-service-account
  prune: true
  matrix:
    - name: input1
      object:
        ref:
          apiVersion: v1
          kind: ConfigMap
          name: input-configmap
  templates:
    - object:
        apiVersion: v1
        kind: ConfigMap
        metadata:
          name: "templated-configmap"
        data:
          y: "{{ matrix.input1.x }}"
    - raw: |
        apiVersion: v1
        kind: ConfigMap
        metadata:
          name: "templated-configmap-from-raw"
        data:
          z: "{{ matrix.input1.x }}"
```

The above manifests show a simple example that will create two ConfigMaps from one input ConfigMap. The individual fields
possible in `ObjectTemplate` are described further down.

## Spec fields

The following fields are supported in `spec`.

### serviceAccountName

`ObjectTemplate` requires a service account to access cluster objects. This is required when it gathers input objects
for the matrix and when it applies rendered objects. Please see [security](../../security.md) for some important notes!

For this to work, the referenced service account must have at least `GET`, `CREATE` and `UPDATE` permissions for
the involved objects and kinds. For the above example, the following service account would be enough:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: example-template-service-account
  namespace: default
---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: example-template-service-account
  namespace: default
rules:
  - apiGroups: [""]
    resources: ["configmaps"]
    resourceNames: ["templated-configmap-from", "templated-configmap-from-raw"]
    verbs: ["*"]
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: example-template-service-account
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: example-template-service-account
subjects:
  - kind: ServiceAccount
    name: example-template-service-account
    namespace: default
```

### interval

Specifies the interval at which the `ObjectTemplate` is reconciled.

### suspend

If set to `true`, reconciliation is suspended.

### prune

If `true`, the Template Controller will delete rendered objects when either the `ObjectTemplate` gets deleted or when
the rendered object disappears from the rendered objects list.

### matrix

The `matrix` defines a list of matrix entries, which are then used as inputs into the templates. Each entry results in
a list of values associated with the entry name. All lists are then multiplied together to form the actual matrix of
input values.

Each matrix entry has a `name`, which is later used to identify the value in the template.

As an example, if you have two entries with simple lists with the following values:

```yaml
matrix:
- name: input1
  list:
    - a: v1
      b: v2
- name: input2
  list:
    - c: v3
      d: v4
```

It will lead to the following matrix:

```yaml
- input1:
    a: v1
    b: v2
  input2:
    c: v3
    d: v4
```

Now take the following matrix example with an entry with two list items:

```yaml
matrix:
- name: input1
  list:
    - a: v1
      b: v2
    - a: v1_2
      b: v2_2
- name: input2
  list:
    - c: v3
      d: v4
```

It will lead to the following matrix:

```yaml
- input1:
    a: v1
    b: v2
  input2:
    c: v3
    d: v4
- input1:
    a: v1_2
    b: v2_2
  input2:
    c: v3
    d: v4
```

Each input value is then used as input when rendering the templates. In the above examples, it means that all templates
are rendered twice, once with `matrix.input1` set to the first input value and the second time with the second input
value.

The following matrix entry types are supported:

#### list

This is the simplest form and represents a list of arbitrary objects. See the above examples.

Due to the use of [controller-gen](https://github.com/kubernetes-sigs/controller-tools) and an internal
[limitation](https://github.com/kubernetes-sigs/controller-tools/issues/461) in regard to validation and CRD generation,
list elements must be objects at the moment. A future version of the Template Controller will support arbitrary values
(e.g. numbers and strings) as elements.

#### object

This refers an object on the cluster. The object is read by the controller and then used as an input value for the
matrix. Example:

```yaml
matrix:
- name: input1
  object:
    ref:
      apiVersion: v1
      kind: ConfigMap
      name: input-configmap
```

The referenced object can be of any kind, but the used [service account](#serviceaccountname) must have access to the
referenced object. The read object is then wholly used as matrix input.

To only use a sub-part of the referenced object, set `jsonPath` to a valid [JSON Path](https://goessner.net/articles/JsonPath/)
pointing to the subfield(s) that you want to use. Example:

```yaml
matrix:
- name: input1
  object:
    ref:
      apiVersion: v1
      kind: ConfigMap
      name: input-configmap
      jsonPath: .data
```

This will make the data field available as input instead of the full object, meaning that values can be used inside the
templates by simply referring `{{ matrix.input1.my_key }}` (no `.data` required).

In case you want to interpret a subfield as an input list instead of a single value, set `expandLists` to `true`.
Example:

```yaml
matrix:
- name: input1
  object:
    ref:
      apiVersion: templates.kluctl.io/v1alpha1
      kind: ListGithubPullRequests
      name: list-gh-prs
      jsonPath: status.pullRequests
      expandLists: true
```

This will lead to one matrix input per list element at `status.pullRequests` instead of a single matrix input that
represents the list.

### templates

`templates` is a list of template objects. Each template object is rendered and applied once per entry from the
multiplied matrix inputs. When rendering, the context contains the global variable `matrix` representing the current
entry. `matrix` has one member field per named matrix input.

In the lists example from above, this would for example give `matrix.input1` and `matrix.input2` for each render
invocation.

In case a template object is missing the namespace, it is set to the namespace of the `ObjectTemplate` object.

The [service account](#serviceaccountname) used for the `ObjectTemplate` must have permissions to get and apply the
resulting objects.

There are currently two forms of template objects supported, `object` and `raw`. `object` is an inline object where
each string field is treated as independent template to render. `raw` represents one large (multi-line) string that
is rendered in one-go and then unmarshalled as yaml/json.

It is recommended to prefer `object` over `raw` and only revert to `raw` templates when you need to perform advanced
templating (e.g. `{% if ... %}` or other control structures) or when it is important to treat a field as non-string
(e.g. boolean or number) when unmarshalled into an object. An example for such case would be if you want to use a
template value for `replicas` of a `Deployment`, which MUST be a number.

Example for an `object`:

```yaml
templates:
- object:
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: "templated-configmap"
    data:
      y: "{{ matrix.input1.x }}"
```

Example for a `raw` template object:

```yaml
templates:
- raw: |
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: "templated-configmap-from-raw"
    data:
      z: "{{ matrix.input1.x }}"
```

See [templating](../../templating.md) for more details on the templating engine.