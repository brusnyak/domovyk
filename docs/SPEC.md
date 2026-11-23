# SPEC: domovyk

Status: draft, in-development. Scope: v0.0 to v0.3.

## Purpose

Open-source exit planner: AWS to EU cloud. Go CLI that runs locally with a read-only IAM role. Produces an inventory, a portability score per service, a price comparison, and OpenTofu HCL for the EU target.

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
| generate | scored `inventory.json` | OpenTofu HCL (`hcloud_*` / `ovh_*`) for portable/rework resources | v0.3 |
| cut over | plan | runbook with rollback: `pgcopydb --follow` (Postgres), `rclone` (S3), Velero (volumes), DNS TTL lowered 48 h ahead | v0.4 (roadmap) |
| operate | target | Prometheus, Grafana, Loki, OTel; monthly `domovyk drift` job | v0.5 (roadmap) |

## Credentials and offline operation

`discover` uses the AWS SDK default credential chain (environment, `~/.aws` config, assume-role, IMDS). Named profile is optional via `--profile`. Secret keys are never accepted on the command line (issue #4 reliability gate).

**Offline constraint:** discovery requires refreshable AWS credentials reachable from the SDK chain. EC2 IMDS (`169.254.169.254`) is only available on EC2; running from a laptop with an assumed-role ARN in `~/.aws/config` requires SSO/refreshable credentials or static `~/.aws/credentials`. Live discovery is blocked until credential-chain config is fixed by the operator.

Price lists and mapping data are bundled snapshots, not fetched live. This keeps the CLI air-gap compatible.

## Provider plugin interface (declared-only)

The following provider plugin contract is documented here as a future target; nothing in code implements it yet:

- `Discover()`: list resources from the source or target.
- `Price()`: return priced rows with `checked_on`.
- `Map()`: apply mapping data to resources.
- `Generate()`: emit OpenTofu (v0.3+).

A contract test suite runs against every plugin. Only the AWS source and the mapping data are in scope before v0.3.

## Price-list source decision

Price lists are bundled as a dated snapshot in `map/prices.yaml`, each row carrying `checked_on` and `status`. They are NOT fetched live. Rationale: the CLI runs offline/local-first; fetching live would require network access, JS-rendered pages (OVHcloud), and introduce a trust/availability dependency not compatible with an air-gapped operator machine. Operators may replace `map/prices.yaml` with fresher rows at any time; each row must keep `checked_on`, `status`, and `source_url`.

## Data

- Mapping: `map/domovyk-map.yaml`, schema in `map/SCHEMA.md`. CC-BY-4.0 intended.
- Code: Apache-2.0.
- `inventory.json` contains account IDs, resource names and tags. Treat it as sensitive; write it with file mode 0600.

## Inventory schema (versioned)

`inventory.json` carries a top-level `manifest` block:

| Field | Type | Meaning |
|---|---|---|
| `manifest.schema` | string | Inventory format version (`"1"`). |
| `manifest.generated_at` | string | RFC3339 UTC timestamp of production. |
| `manifest.account_id` | string | AWS account ID from `sts:GetCallerIdentity` (best effort). |
| `manifest.region` | string | Primary region scanned. |
| `manifest.regions` | string[] | Every region discovery was attempted against. |
| `manifest.services.<name>` | object | Per-service coverage: `attempted`, `completed`, `resource_count`, `first_error`, `pagination_complete`. |
| `manifest.total_resources` | int | Sum of all resources collected. |
| `resources[]` | array | Discovered AWS resources. |

Discovery must never report an incomplete scan as complete: each service records `attempted`/`completed`/`pagination_complete` plus `first_error` on hard failure.

## Credentials

`discover` uses the AWS SDK default credential chain (environment, `~/.aws` config, assume-role, IMDS). Named profile is optional via `--profile`. Secret keys are never accepted on the command line (issue #4 reliability gate).

## Build order

v0.0 skeleton, README, read-only IAM policy, 30 seed mapping rows, reference AWS setup in OpenTofu (`tf/reference-aws/`).
v0.1 discover, score, sample exit plan. v0.2 price from a FOCUS export. v0.3 generate OpenTofu HCL. v0.4 to v0.5: roadmap only.

## Price-list source decision

Price lists are bundled as a dated snapshot in `map/prices.yaml`, each row carrying `checked_on` and `status`. They are NOT fetched live. Rationale: the CLI runs offline/local-first; fetching live would require network access, JS-rendered pages (OVHcloud), and introduce a trust/availability dependency not compatible with an air-gapped operator machine. Operators may replace `map/prices.yaml` with fresher rows at any time; each row must keep `checked_on`, `status`, and `source_url`.

## Open design points

- Mapping data format for multi-service patterns (for example, ALB plus target groups). Not decided.
