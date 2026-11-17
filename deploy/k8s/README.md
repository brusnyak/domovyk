# Domovyk on kind (local Kubernetes) — container proof

This directory proves `domovyk` builds as a container image, runs as non-root on Kubernetes, and scores a sample AWS inventory end to end. It is a **single-node kind cluster with sample data only** — not a production deployment, not a cloud cluster, no ingress, no autoscaling.

## What this does

Builds the `domovyk:dev` image (distroless nonroot, ~36 MB), loads it into a local `kind` cluster, and runs a one-off `batch/v1 Job` that executes `domovyk score` against a 5-resource sample inventory bundled in a ConfigMap. The Job writes per-resource class counts (`portable` / `rework` / `rewrite` / `unknown`) and exits `Completed`.

## Prerequisites (verified 2026-11-17)

- macOS arm64, Colima (Docker Desktop not installed)
- `kind v0.33.0`, `kubectl 1.37.2` (brew)
- Colima started: `colima start --cpu 2 --memory 4 --disk 20`
- `export DOCKER_HOST=unix://$HOME/.colima/default/docker.sock`

## Exact commands and real output

Run from `~/Developer/domovyk`.

### 1. Build the image

```sh
docker build -t domovyk:dev .
docker images domovyk:dev
```

Real output:

```
DEPRECATED: The legacy builder is deprecated and will be removed in a future release.
            Install the buildx component to build images with BuildKit:
            https://docs.docker.com/go/buildx/
Sending build context to Docker daemon  51.02MB
Step 1/6 : FROM gcr.io/distroless/static-debian12:nonroot
 ---> afa5c872c891
Step 2/6 : WORKDIR /app
 ---> Using cache
 ---> 8bc029f27908
Step 3/6 : COPY .dist/domovyk /app/domovyk
 ---> 7eed8282a721
Step 4/6 : COPY map/domovyk-map.yaml /app/map/domovyk-map.yaml
 ---> aec6c3947831
Step 5/6 : USER nonroot:nonroot
 ---> Running in 8dd0b5e72899
 ---> Removed intermediate container 8dd0b5e72899
Step 6/6 : ENTRYPOINT ["/app/domovyk"]
 ---> Running in 0247d56c0ae6
 ---> Removed intermediate container 8dd0b5e72899
Successfully built 3d6efa9f9824
Successfully tagged domovyk:dev
IMAGE         ID             DISK USAGE   CONTENT SIZE   EXTRA
domovyk:dev   3d6efa9f9824   36.5MB         10.4MB
```

### 2. Verify the binary runs under the nonroot user

```sh
docker run --rm domovyk:dev --help
```

Real output:

```
domovyk - AWS to EU exit planner (v0.1)

Usage:
  domovyk discover [--region REGION] [--profile PROFILE] [--out FILE]
                   Discover read-only AWS inventory into inventory JSON.
  domovyk score    [--inventory FILE] [--mapping FILE] [--out FILE]
                   Score each resource as portable / rework / rewrite.
  domovyk price    [--focus FILE]
                   Compare AWS billing against OVHcloud / Hetzner (v0.2).
```

### 3. Create the kind cluster and load the image

```sh
kind create cluster --name domovyk
kind load docker-image domovyk:dev --name domovyk
kubectl config current-context
```

Real output:

```
Creating cluster "domovyk" ...
 ✓ Ensuring node image (kindest/node:v1.37.0)
 ✓ Preparing nodes
 ✓ Writing configuration
 ✓ Starting control-plane
 ✓ Installing CNI
 ✓ Installing StorageClass
Set kubectl context to "kind-domovyk"
You can now use your cluster with:

kubectl cluster-info --context kind-domovyk
kind-domovyk
Image: "domovyk:dev" with ID "sha256:3d6efa9f98240b306fed7cae5e8314ac040b98932ecf20c8ae6a0c65d4b4ac21" not yet present on node "domovyk-control-plane", loading...
```

Context must be `kind-domovyk`. Never apply to any other context.

### 4. Deploy the Job

```sh
kubectl apply -f deploy/k8s/job.yaml
kubectl wait --for=condition=complete job --all --timeout=120s
kubectl get pods -o wide
kubectl logs job/domovyk-score
kubectl get job -o yaml | grep -E "runAsNonRoot|readOnlyRootFilesystem"
```

Real output:

```
configmap/domovyk-sample-inventory created
job.batch/domovyk-score created
job.batch/domovyk-score condition met
NAME                  READY   STATUS      RESTARTS   AGE   IP           NODE                    NOMINATED NODE   READINESS GATES
domovyk-score-hb454   0/1     Completed   0          9s    10.244.0.2   domovyk-control-plane   <none>          <none>
=== JOB LOGS ===
scored 5 resources -> /tmp/score.json
{
  "portable": 2,
  "rework": 1,
  "rewrite": 1,
  "unknown": 1,
  "total": 5
}
      runAsNonRoot: true
      readOnlyRootFilesystem: true
```

The Job reaches `Complete`. The log shows the scored sample inventory: `portable: 2`, `rework: 1`, `rewrite: 1`, `unknown: 1` — `unknown` is the unmapped `glue` resource, as expected.

## Hardened container

The Job runs with `runAsNonRoot: true`, `runAsUser: 65532` (distroless `nonroot`), `readOnlyRootFilesystem: true`, `seccompProfile.type: RuntimeDefault`, `capabilities.drop: ["ALL"]`, `allowPrivilegeEscalation: false`, no service-account token, no secrets, no AWS credentials, `imagePullPolicy: Never`, resource limits (200m CPU / 64 Mi memory).

## What this does not prove

- Single-node kind on a laptop — not a production cluster, no HA, no multi-node scheduling.
- No ingress, no service exposure, no autoscaling, no cloud provider integration.
- Sample data only (5 resources in `eu-central-1`) — not a real AWS account discovery.
- Does not prove `discover` against live AWS, `price` against real billing, or `generate`.
- The image is not pushed to any registry — `imagePullPolicy: Never` runs the locally built image only.

## Teardown

```sh
kind delete cluster --name domovyk
colima stop
```
