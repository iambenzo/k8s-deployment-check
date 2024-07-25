# Kubernetes Deployment Checker

This application is designed to run as a [Kubernetes Init Container](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/).

It blocks your application from running until another deployment has reached it's target Ready state.

## Usage

There are three command line arguments that you can pass to this application:

| Command       | Notes                                                           |
| ------------- | --------------------------------------------------------------- |
| `-deployment` | Name of Deployment to monitor                                   |
| `-namespace`  | Namespace of Deployment to monitor                              |
| `-sleep`      | Number of seconds between checking the Deployment (default: 10) |

## Deployment

First, we need to ensure that the init container has privileges to access the specific parts of the Kubernetes API to do it's job:

```yaml
# You need a Service Account for each namespace this
# init container will be used in
apiVersion: v1
kind: ServiceAccount
metadata:
  name: deployment-checker
  namespace: <NAMESPACE_THAT_INIT_RUNS_IN>
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: deployment-checker
rules:
- apiGroups: ["apps", "extensions"]
  resources: ["deployments"]
  verbs: ["get"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: deployment-checker
# Subjects is a list
# You can use this to apply the same permissions
# To each Service Account you deploy
subjects:
- kind: ServiceAccount
  name: deployment-checker
  namespace: <NAMESPACE_THAT_INIT_RUNS_IN>
roleRef:
  kind: ClusterRole
  name: deployment-checker
  apiGroup: rbac.authorization.k8s.io
```

Now that we have done the prep work, we can add this init container to your deployment configuration:

```yaml
spec:
  serviceAccount: deployment-checker
  serviceAccountName: deployment-checker
  initContainers:
    - name: deployment-check
      image: <CONTAINER_REGISTRY>/deployment-check:0.0.1
      imagePullPolicy: IfNotPresent
      args:
        - -deployment=<DEPLOYMENT_TO_WATCH>
        - -namespace=<NAMESPACE_OF_DEPLOYMENT_TO_WATCH>
  containers:
	- name: ...
```
