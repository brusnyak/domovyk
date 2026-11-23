# domovyk

Open-source planner for moving AWS workloads to EU providers (OVHcloud, Hetzner). A Go CLI that runs on your own machine with a read-only IAM role.

**Status: in-development (pre-alpha).** `discover` and `score` have been run against a live AWS account (account `004406188742`, eu-north-1 — 10 IAM roles discovered, all scored `unknown` due to no IAM mapping entries). `price` works against a FOCUS 1.2 billing export. `generate` produces OpenTofu HCL for Hetzner and OVHcloud from a scored inventory.

## What it does

| Command | What it does |
|---------|--------------|
| `discover` | Read-only inventory of an AWS account into `inventory.json`. Collects EC2 instances/EBS, RDS instances, S3 buckets, EKS clusters, Lambda functions, IAM roles, plus account identity via STS. Paginates every paginated API, threads caller context with a 2-minute per-service timeout, converts AWS responses nil-safely, and records per-service coverage in a versioned inventory manifest. Credentials come from the AWS SDK default chain (env, `~/.aws` config, assume-role, IMDS); secret keys are never accepted on the command line. |
| `score` | Classify each resource as `portable`, `rework`, or `rewrite` using the open mapping data (`map/domovyk-map.yaml`). All 30 mapping rows are verified against MoveInfra's EU equivalents guide and Hetzner's live pricing page. |
| `price` | Compare your AWS billing export (FOCUS 1.2) with sourced Hetzner and OVHcloud price lists. Every price row carries a `checked_on` date and an explicit `verified`/`unverified` status. Try it with `internal/price/testdata/focus.csv`. |
| `generate` | Emit OpenTofu HCL (`hcloud_*` / `ovh_*` resources) for portable/rework resources from a scored inventory. Rewrite and unknown resources are skipped. Review and right-size before `tofu apply`. |

## Quick start

```sh
# Build
go build ./cmd/domovyk

# 1. Discover (needs AWS credentials; read-only)
./domovyk discover --region eu-central-1 --out inventory.json

# 2. Score
./domovyk score --inventory inventory.json --out score.json

# 3. Generate OpenTofu for Hetzner (default) or OVHcloud
./domovyk generate --inventory inventory.json --out generated --target hetzner
./domovyk generate --inventory inventory.json --out generated --target ovhcloud

# 4. Try price with the bundled synthetic FOCUS export
./domovyk price --focus internal/price/testdata/focus.csv --out price.json
```

## Honest limits

- It plans, it does not move. `generate` emits OpenTofu HCL; you review, right-size, and `tofu apply` yourself. Nothing is provisioned automatically.
- It will say no to Lambda-heavy accounts. Those are `rewrite`, not a lift-and-shift.
- No compliance claims. This is not legal advice. The EU Data Act text is the reference.
- Prices are point-in-time. Read the `checked_on` date before trusting a number.
- `generate` instance types are best-effort approximations, not verified capacity matches.
- The IAM policy in `iam/` has been reviewed against the actual SDK calls but not yet run through the AWS IAM policy simulator.

## Roadmap

v0.0–v0.2 are implemented (`discover`, `score`, `price`, `generate`). v0.3–v0.5 (measured cut-over, monthly drift job, plugin contract tests) are roadmap-only. See `docs/ROADMAP.md`.

## Docs

- `docs/SPEC.md`: architecture, versioned inventory schema, six stages
- `docs/THREAT-MODEL.md`: draft threat model
- `docs/RESEARCH.md`: competitors and open questions, with verified and unverified labels
- `docs/ROADMAP.md`, `PROGRESS.md`

## Build

```sh
go build ./cmd/domovyk
go vet ./...
go test ./...
```

CI: GitHub Actions runs build/vet/test on every push and PR (`.github/workflows/ci.yml`).

Container proof: `deploy/k8s/README.md` — builds `domovyk:dev`, runs `domovyk score` on a local kind cluster, real output captured.

## Licence

Code: Apache-2.0 (full text in `LICENSE`). Mapping data in `map/`: intended CC-BY-4.0 (licence file for the data not added yet).
