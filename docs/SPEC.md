# SPEC: domovyk

Status: draft, pre-alpha. Scope for this repo: v0.0 to v0.2.

## Purpose

Open-source exit planner: AWS to EU cloud. Go CLI. Runs locally with a read-only IAM role. Produces an inventory, a portability score per service, a price comparison and, later, a generated target.

## Non-negotiables

- Read-only against the customer's AWS account.
- Local-first: the tool sends no customer data to any Domovyk-operated service (there is none).
- Published threat model (`docs/THREAT-MODEL.md`).
- No compliance claims.
- Every price row carries a `checked_on` date.
- Say no to Lambda-heavy accounts.

## Stages

| Stage | Input | Output | Version |
|---|---|---|---|
| discover | AWS account, read-only IAM | `inventory.json` | v0.1 |
| score | `inventory.json`, `map/domovyk-map.yaml` | per-service class: `portable`, `rework`, `rewrite` | v0.1 |
| price | FOCUS 1.2 billing export, OVHcloud and Hetzner price lists | cost table with `checked_on` per row | v0.2 |
| generate | exit plan | OpenTofu modules: OVH Managed Kubernetes; Talos or k3s on Hetzner; Argo CD, cert-manager, external-dns, External Secrets | v0.3 (roadmap) |
| cut over | plan | runbook with rollback: `pgcopydb --follow` (Postgres), `rclone` (S3), Velero (volumes), DNS TTL lowered 48 h ahead | v0.4 (roadmap) |
| operate | target | Prometheus, Grafana, Loki, OTel; monthly `domovyk drift` job | v0.5 (roadmap) |

## Provider plugin interface

Each provider implements:

- `Discover()`: list resources from the source or target.
- `Price()`: return priced rows with `checked_on`.
- `Map()`: apply mapping data to resources.
- `Generate()`: emit OpenTofu (v0.3+).

A contract test suite runs against every plugin. Only the AWS source and the mapping data are in scope before v0.3.

## Data

- Mapping: `map/domovyk-map.yaml`, schema in `map/SCHEMA.md`. CC-BY-4.0 intended.
- Code: Apache-2.0.
- `inventory.json` contains account IDs, resource names and tags. Treat it as sensitive; write it with file mode 0600.

## Build order

v0.0 skeleton, README, read-only IAM policy, 30 seed rows, reference AWS setup in OpenTofu (EC2, RDS, S3, ALB). Note: the reference setup is planned only in this repo (`tf/reference-aws/README.md`), no `.tf` files yet.

v0.1 discover, score, sample exit plan. v0.2 price from a FOCUS export. v0.3 to v0.5: roadmap only.

## Open design points

- How price lists are obtained (bundled snapshot or fetched). Not decided.
- Mapping data format for multi-service patterns (for example, ALB plus target groups). Not decided.
