# ROADMAP: domovyk

## v0.0: skeleton (done)

- [x] Folder layout and README (`README.md`)
- [x] Spec and threat-model draft (`docs/`)
- [x] Go CLI stub with `discover`, `score`, `price` printing "not implemented" (`cmd/domovyk/`)
- [x] Read-only IAM policy, reviewed against actual SDK calls (`iam/`)
- [x] 30 seed mapping rows, all unverified (`map/`)
- [x] Reference AWS setup plan, no `.tf` files (`tf/reference-aws/README.md`)
- [ ] Go module path set (`OWNER` placeholder replaced)
- [ ] Full Apache-2.0 attribution boilerplate in source files

## v0.1: discover, score, sample exit plan (implemented)

- [x] Discover EC2, RDS, S3 (metadata), ELB, EKS, Lambda (list only), IAM (list only) into `inventory.json`
- [x] Score each resource against the mapping data
- [x] Versioned inventory manifest: account ID, region, per-service coverage, pagination completion
- [x] Run against `wocus` and `recareo-web` read-only (owner decision — credentials expired; re-run pending)
- [ ] Sample exit plan for the reference AWS setup

## v0.2: price from FOCUS export (implemented)

- [x] Parse FOCUS 1.2 billing export
- [x] Price lists for OVHcloud and Hetzner with `checked_on` on every row
- [x] Cost comparison table with stale-price warning (unverified rows printed as UNVERIFIED)
- [ ] Run against a real FOCUS billing export (not yet available)

## Roadmap only (not in scope)

- v0.3: generate OpenTofu modules (OVH Managed Kubernetes; Talos or k3s on Hetzner; Argo CD, cert-manager, external-dns, External Secrets)
- v0.4: measured cut-over with runbook and rollback (pgcopydb --follow, rclone, Velero, DNS TTL lowered 48 h ahead). Owner's `wocus` and `recareo-web` as first real case.
- v0.5: monthly `domovyk drift` job; provider plugin contract test suite

## Gates (from GitHub issues)

- Issue #4 (reliability): done — pagination, context+timeouts, nil-safe conversion, no secret CLI flags, inventory manifest, tests.
- Issue #3 (product validation): research only; no executor build without partner validation.
- Issue #1 (roadmap umbrella): P0 acceptance items in progress (docs reconciled; mapping schema validation and price provenance next).
