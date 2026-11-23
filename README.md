# domovyk

Open-source planner for moving AWS workloads to EU providers (OVHcloud, Hetzner). A Go CLI that runs on your own machine with a read-only IAM role.

**Status: in-development (pre-alpha).** `discover` and `score` run against live AWS accounts (verified against account `004406188742` in eu-north-1; discovery collected 10 IAM roles, all scored `unknown` due to no IAM mapping entries). `discover` paginates every paginated API, threads caller context with a 2-minute per-service timeout, converts AWS responses nil-safely, and records per-service coverage in a versioned inventory manifest. `price` compares a FOCUS 1.2 billing export with sourced Hetzner and OVHcloud price lists; every price row carries a `checked_on` date and an explicit `verified`/`unverified` status. `generate` and later stages do not exist.

## What it does

1. `discover`: read-only inventory of an AWS account into `inventory.json`. Collects EC2 instances/EBS, RDS instances, S3 buckets, EKS clusters, Lambda functions, IAM roles, plus account identity via STS. Uses the AWS SDK default credential chain (env, `~/.aws` config, assume-role, IMDS); secret keys are never accepted on the command line.
2. `score`: classify each resource as `portable`, `rework`, or `rewrite` using the open mapping data (`map/domovyk-map.yaml`).
3. `price`: compare your AWS billing export (FOCUS 1.2) with OVHcloud and Hetzner price lists. Every price row carries a `checked_on` date.
4. `generate`: OpenTofu modules for the target. Roadmap only, not in scope yet.
5. `cut over` and `operate`: runbook with rollback and a monthly drift check. Roadmap only.

## Honest limits

- It does not move anything. v0.1 produces a sample exit plan, nothing more.
- It will say no to Lambda-heavy accounts. Those are `rewrite`, not a lift-and-shift.
- No compliance claims. This is not legal advice. The EU Data Act text is the reference.
- Prices are point-in-time. Read the `checked_on` date before trusting a number.
- The 30 mapping rows in `map/` are unverified seed data written from general knowledge.
- The IAM policy in `iam/` has been reviewed against the actual SDK calls (`sts:GetCallerIdentity` added) but not yet run through the AWS IAM policy simulator.

## Roadmap

In scope now:
- v0.0: skeleton, README, read-only IAM policy, 30 seed mapping rows, reference AWS setup plan.
- v0.1: discover, score, sample exit plan (implemented; discovery tested against live AWS).
- v0.2: price from a FOCUS export (implemented; not yet run against a real billing export).

Roadmap only: v0.3 generate, v0.4 measured cut-over, v0.5 drift job and plugin contract tests. See `docs/ROADMAP.md`.

## Docs

- `docs/SPEC.md`: architecture and six stages
- `docs/THREAT-MODEL.md`: draft threat model
- `docs/RESEARCH.md`: competitors and open questions, with verified and unverified labels
- `docs/ROADMAP.md`, `PROGRESS.md`

## Build

```sh
go build ./cmd/domovyk
go vet ./...
go test ./...
```

Container proof: `deploy/k8s/README.md` — builds `domovyk:dev`, runs `domovyk score` on a local kind cluster, real output captured.

## Licence

Code: Apache-2.0 (full text in `LICENSE`). Mapping data in `map/`: intended CC-BY-4.0. The licence file for the data is not added yet.
