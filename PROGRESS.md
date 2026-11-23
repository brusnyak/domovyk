# PROGRESS: domovyk (for Yegor)

Status: pre-alpha. Scope: v0.0 to v0.2 only.

## Done (files exist; nothing run)

- [x] `README.md`, `LICENSE` (placeholder), `go.mod` (`OWNER` placeholder)
- [x] `docs/SPEC.md`, `docs/THREAT-MODEL.md` (draft), `docs/RESEARCH.md`, `docs/ROADMAP.md`
- [x] `cmd/domovyk/main.go` stub (not compiled; go not run)
- [x] `iam/domovyk-readonly-policy.json` (UNREVIEWED), `iam/README.md`
- [x] `map/domovyk-map.yaml` (30 rows, all unverified), `map/SCHEMA.md`
- [x] `map/prices.yaml` (Hetzner + OVHcloud rows with checked_on and status; OVHcloud compute/managed_db marked unverified pending page confirmation)
- [x] `tf/reference-aws/README.md` (plan only)
- [x] `deploy/k8s/README.md` (container proof: kind + domovyk score job, real output captured)
- [x] `docs/AWS-E2E-FINDINGS.md` (real AWS eu-north-1 discovery: 10 IAM roles, all unknown/unmapped)

## Not done

- [ ] Go build check (`go build ./cmd/domovyk`). Not run.
- [ ] IAM policy review and simulator check
- [ ] Verify the 30 seed rows
- [ ] GitHub and HN search for existing tools (required before v0.1)
- [ ] Price-list source decision (bundle vs fetch)
- [ ] Full Apache-2.0 text in `LICENSE`
- [ ] `git init`, first commit (no push)
- [ ] Reference AWS `.tf` files (v0.0 item, sandbox account only)
- [ ] Portfolio drafts: owner review; lingonberry-ops has no data

## First 5 steps

1. Read `docs/RESEARCH.md`, section "Must re-check before v0.1", and run the GitHub and HN searches listed there. Record hits with URL and date.
2. Replace `OWNER` in `go.mod`, then run `go build ./cmd/domovyk` and `go run ./cmd/domovyk help`.
3. Check `iam/domovyk-readonly-policy.json` against AWS docs and the IAM policy simulator in your sandbox account.
4. Verify 5 seed rows in `map/domovyk-map.yaml` against AWS and provider docs. Set `checked_on` and `status` only for rows you checked.
5. Decide the price-list source (bundled snapshot or fetch) and note it in `docs/SPEC.md` open points.

## Unverified items written in these files

Every item below is marked UNVERIFIED in its file. Do not cite as fact.

- Terraformer archived 16 Mar 2026 (`docs/RESEARCH.md`)
- MinIO OSS repo archived Feb 2026 (`docs/RESEARCH.md`)
- Steampipe plugin licences (`docs/RESEARCH.md`)
- Former2 about 49% coverage (`docs/RESEARCH.md`)
- Hetzner cloud-controller v1.30.0 crash (`docs/RESEARCH.md`)
- Hetzner price rises in 2026 (`docs/RESEARCH.md`)
- Any price number: none written; none should be added without `checked_on`
- 30 seed mapping rows, from general knowledge (`map/domovyk-map.yaml`)
- IAM policy simulator check pending (`iam/`)
- Threat model is a draft, not reviewed (`docs/THREAT-MODEL.md`)
- Reference AWS plan has no cost figure (`tf/reference-aws/README.md`)

Verified by the owner (cited as verified in `docs/RESEARCH.md`): EU Data Act switching charges end 12 Jan 2027; parallel multi-cloud egress allowed at cost.
