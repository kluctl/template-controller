# Use Case: Transformation of Secrets/Objects

There are cases where an object can not be created before another object is created by some other component inside the
cluster, meaning that you have no control over the input object.

A simple example is the [Zalando Postgres Operator](https://github.com/zalando/postgres-operator), which allows you
to create a Postgres database with a Custom Resource. Inside the CR, you can define databases and users to be
auto-created. When the operator creates these databases and users, it also auto-creates Kubernetes secrets with the
credentials allowing you to access the databases.

These secrets can however not be used directly when connecting to the databases, as you'd usually have to build some
connection urls (e.g. JDBC urls). Usually, one would create some kind of init script or something like that to 
build this url and then pass it to the application that wants to use it.

The Template Controller allows an alternative solution.

## Using ObjectTemplate to transform secrets

Let's assume you have a [sample](https://github.com/zalando/postgres-operator/blob/master/docs/user.md#create-a-manifest-for-a-new-postgresql-cluster)
Postgres database deployed via the Zalando Postgres Operator. The operator has also created the following secret:

```yaml
apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: foo-user.acid-minimal-cluster.credentials.postgresql.acid.zalan.do
  namespace: default
data:
  password: aHNiSVF6MFJJa0hTd2ZxS1NiTG5YV3dUQUVqcUtTNFpvU2dyOXp4b3pzMmJvTE02WWl0eTE0YjJTZlNFTHExdw==
  username: Zm9vX3VzZXI=
```

Based on that secret, you'd like to create a new secret with the JDBC url generated.

## RBAC

The ObjectTemplate requires a service account with proper access rights for the involved secrets:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: postgres-secret-transformer
  namespace: default
---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: postgres-secret-transformer
  namespace: default
rules:
  - apiGroups: [""]
    resources: ["secrets"]
    # give the ObjectTemplate access to the two involved secrets
    resourceNames: ["zalando.acid-minimal-cluster.credentials.postgresql.acid.zalan.do", "transformed-postgres-secret"]
    verbs: ["*"]
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: postgres-secret-transformer
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: postgres-secret-transformer
subjects:
  - kind: ServiceAccount
    name: postgres-secret-transformer
    namespace: default
```

## ObjectTemplate

Use the following [`ObjectTemplate`](./spec/v1alpha1/objecttemplate.md) to perform the transformation:

```yaml
apiVersion: templates.kluctl.io/v1alpha1
kind: ObjectTemplate
metadata:
  name: postgres-secret-transformer
  namespace: default
spec:
  serviceAccountName: postgres-secret-transformer
  prune: true
  matrix:
    - name: secret
      object:
        ref:
          apiVersion: v1
          kind: Secret
          name: zalando.acid-minimal-cluster.credentials.postgresql.acid.zalan.do
  templates:
  - object:
      apiVersion: v1
      kind: Secret
      metadata:
        name: "transformed-postgres-secret"
      stringData:
        jdbc_url: "jdbc:postgresql://acid-minimal-cluster/zalando?user={{ matrix.secret.data.username | b64decode }}&password={{ matrix.secret.data.password | b64decode }}"
        # sometimes the key names inside a secret are not what another component requires, so we can simply use different names if we want
        username_with_different_key: "{{ matrix.secret.data.username | b64decode }}"
        password_with_different_key: "{{ matrix.secret.data.password | b64decode }}"
```

This will lead to the following `transformed-postgres-secret`

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: transformed-postgres-secret
  namespace: default
type: Opaque
data:
  jdbc_url: amRiYzpwb3N0Z3Jlc3FsOi8vaG9zdC9kYXRhYmFzZT91c2VyPWZvb191c2VyJnBhc3N3b3JkPWJVUU52Zkd4amduQUdiaEhOWkZkamtwZFFYbnk1aDdXNGlFU1YyWUxVNnVrRHdXWjBPMjdRb0NBdUJTTnF3TVk=
  password_with_different_key: YlVRTnZmR3hqZ25BR2JoSE5aRmRqa3BkUVhueTVoN1c0aUVTVjJZTFU2dWtEd1daME8yN1FvQ0F1QlNOcXdNWQ==
  username_with_different_key: Zm9vX3VzZXI=
```

Base64 decoding the secret data will show:

```yaml
jdbc_url: jdbc:postgresql://host/database?user=foo_user&password=bUQNvfGxjgnAGbhHNZFdjkpdQXny5h7W4iESV2YLU6ukDwWZ0O27QoCAuBSNqwMY                                                                                                                                                                                      │
password_with_different_key: bUQNvfGxjgnAGbhHNZFdjkpdQXny5h7W4iESV2YLU6ukDwWZ0O27QoCAuBSNqwMY                                                                                                                                                                                                                          │
username_with_different_key: foo_user
```
