# ROADMAP: domovyk

## v0.0: skeleton (in progress)

- [x] Folder layout and README (`README.md`)
- [x] Spec and threat-model draft (`docs/`)
- [x] Go CLI stub with `discover`, `score`, `price` printing "not implemented" (`cmd/domovyk/`)
- [x] Read-only IAM policy, UNREVIEWED (`iam/`)
- [x] 30 seed mapping rows, all unverified (`map/`)
- [x] Reference AWS setup plan, no `.tf` files (`tf/reference-aws/README.md`)
- [ ] Review of IAM policy against AWS docs and a policy simulator
- [ ] Verify or correct the 30 seed rows
- [ ] Go module path set (`OWNER` placeholder replaced)
- [ ] Full Apache-2.0 text in `LICENSE`
- [ ] Repo initialised and first commit (owner decision)

## v0.1: discover, score, sample exit plan

- [ ] Discover EC2, RDS, S3 (metadata), ELB, EKS, Lambda (list only), IAM (list only) into `inventory.json`
- [ ] Score each resource against the mapping data
- [ ] Sample exit plan for the reference AWS setup
- [ ] Run against `wocus` and `recareo-web` read-only (owner decision)

## v0.2: price from FOCUS export

- [ ] Parse FOCUS 1.2 billing export
- [ ] Price lists for OVHcloud and Hetzner with `checked_on` on every row
- [ ] Cost comparison table with stale-price warning

## Roadmap only (not in scope)

- v0.3: generate OpenTofu modules (OVH Managed Kubernetes; Talos or k3s on Hetzner; Argo CD, cert-manager, external-dns, External Secrets)
- v0.4: measured cut-over with runbook and rollback (pgcopydb --follow, rclone, Velero, DNS TTL lowered 48 h ahead). Owner's `wocus` and `recareo-web` as first real case.
- v0.5: monthly `domovyk drift` job; provider plugin contract test suite
