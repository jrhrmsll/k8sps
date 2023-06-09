# k8sps

`k8sps` allows to know all open ports across the nodes of a Kubernetes cluster.

Designed to run as a HTTP application, `k8sps` listen in the port 8080 providing the `/api/scan` endpoint for POST requests.

When a request is made a new port scan is triggered for all Kubernetes nodes in the cluster; producing a text report with multiple lines with the following format:
`<node>: [<port>, <port>,<port>, ..., <port>]`.

A second endpoint `/api/report` allows the user to download the port scan results or report.

## Build
Created with the Go programming language for a containerized environment, it is easy to create a Docker image with `docker build -t portscan .`.

## Helm Chart
A Helm Chart is provided, with defaults values.

## Testing
Any solution for creating a local Kubernetes cluster is enought. The next steps allow testing with Minikube.

1. Create the minikube cluster
```
minikube start
```

2. Create the Docker image
```
docker build -t portscan .
```

3. Load the image inside Minikube
```
minikube image load portscan
```

4. Install the Helm Chart
```
helm install portscan helm/portscan
```

5. Create a port forward for the portscan service
```
kubectl port-forward deployment/portscan 8080:8080
```

6. Trigger a port scan execution
```
curl -X POST  localhost:8080/api/scan
```

7. Get the port scan results or report
```
curl localhost:8080/api/report
```