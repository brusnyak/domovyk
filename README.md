# Domovyk

**Customer-controlled, auditable workflows for moving workloads from AWS to EU cloud providers.**

Domovyk is a local-first Go CLI project. Its intended job is to help operators discover a workload, produce an evidence-backed migration plan, execute approved steps, validate the target, and recover safely when something fails. It should compose established tools such as OpenTofu and workload-specific transfer utilities rather than reinvent them.

<sub>**Status: pre-alpha. Not production-ready.**</sub>

## What exists today

- AWS discovery, scoring, and pricing code is present in the repository; discovery has reportedly been run once against a real AWS account.
- The mapping dataset contains 30 seed rows that remain unverified.
- The IAM policy and threat model are drafts and have not received a full security review.
- Build and test status is not yet verified. No migration execution, target provisioning, cutover, or rollback capability should be assumed.

See [PROGRESS.md](PROGRESS.md) for the verified project state and [docs/ROADMAP.md](docs/ROADMAP.md) for release gates.

## Product direction

Domovyk is not intended to be another cloud cost calculator or generic Terraform generator. Its proposed differentiator is a customer-owned workflow with:

1. Honest discovery coverage and explicit unknowns.
2. A deterministic plan with evidence and a reviewable digest.
3. Human approval before infrastructure writes and traffic changes.
4. Durable checkpoints, safe retries, and resume after interruption.
5. Application-level validation and documented rollback limits.

The first proposed execution slice is one stateless EC2 workload to one Hetzner Cloud server. This scope is a hypothesis pending customer interviews and a safe pilot; databases, Kubernetes, Lambda rewrites, multi-account orchestration, and automatic production cutover are excluded initially.

## Safety principles

- Read-only discovery by default; execution permissions are separate and opt-in.
- No cloud writes during planning.
- No apply without approval bound to the exact plan.
- Fail closed on incomplete discovery, missing permissions, unknown resources, or failed validation.
- Keep customer artifacts local and restrictive; never store credentials in plans or logs.
- Treat rollback as a tested procedure, not a marketing promise.

## Development

Start with the [product strategy](docs/PRODUCT-STRATEGY.md), [reliability requirements](docs/RELIABILITY.md), and [roadmap](docs/ROADMAP.md).

Once Go 1.24+ is installed:

```sh
gofmt -w .
go test ./...
go vet ./...
go build ./cmd/domovyk
```

These commands are required checks; they have not yet been run successfully in this project environment. Do not use real customer accounts for integration tests without explicit approval and a disposable scope.

## Research

- [Competitive landscape and product decision](docs/COMPETITIVE-LANDSCAPE-2026-11.md)
- [Research notes](docs/RESEARCH.md)
- [Threat model](docs/THREAT-MODEL.md)

## Licensing

Code is intended to be Apache-2.0; mapping data is intended to be CC-BY-4.0. Verify the corresponding full license texts are present before distributing releases.
