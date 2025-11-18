# gRPC + Istio Mesh Demo

A production-ready demo showing how **Istio simplifies security, routing, and observability** for gRPC microservices.

### Includes:
- gRPC service and client in Go
- Helm deployment with Istio mTLS, canary routing, and authorization

### Quickstart
```bash
cd grpc-mesh-demo
kubectl apply -f deploy/base/namespace.yaml
helm install payments deploy/helm/ -n payments
kubectl apply -f deploy/istio/
```

### Deploy client using below command
```bash
kubectl apply -f ./deploy/client/client-deployment.yaml
```
