This application is designed to run as a [Kubernetes Init Container](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/).

Pass it a `-namespace` and a `-deployment` name via the command line arguments and it will ensure that the pod it's attached to won't start until the Kubernetes Deployment identified via your input has the required pods in a Ready state.

## Usage

Add to your deployment definition as an `initContainer`:

> TODO: example configuration

Ensure that the init container has access to a Service Account which has permission to get deployments and you should be good to go.

