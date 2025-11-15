# gRPC + Istio Go Demo

This repository contains a minimal Go-based gRPC Payments service and client,
plus Kubernetes & Istio manifests and a Helm chart skeleton.

## Notes

- You must generate Go protobuf code before building server/client:
  ```bash
  protoc --go_out=. --go-grpc_out=. proto/payments.proto
  ```
  (requires `protoc` and `protoc-gen-go`, `protoc-gen-go-grpc` installed)

- Build images and push them to your registry, or use local kind/tilt builds.

## Quickstart (cluster with Istio sidecar injection enabled)
```bash
kubectl apply -f deploy/base/namespace.yaml
kubectl apply -f deploy/base/serviceaccount.yaml
kubectl apply -f deploy/base/deployment.yaml
kubectl apply -f deploy/base/service.yaml
kubectl apply -f deploy/istio/
```

## Run client inside cluster:
```bash
kubectl run -n payments grpc-client --rm -it --restart=Never --image=example/client -- /client
```
