# PROGRESS: domovyk (for Yegor)

Status: in-development (pre-alpha). Scope: v0.0 to v0.3 implemented. Build, vet, and tests pass; GitHub Actions CI on push/PR.

## Done
- [x] `README.md`, `LICENSE` (full Apache-2.0 text), `go.mod` (`OWNER` placeholder)
- [x] `docs/SPEC.md`, `docs/THREAT-MODEL.md` (draft), `docs/RESEARCH.md`, `docs/ROADMAP.md`, `docs/PRICE.md`
- [x] `cmd/domovyk/main.go` — compiled, `go build` and `go vet` pass; CLI with `discover`, `score`, `price`, `generate`
- [x] `internal/inventory` — Resource + Inventory types, JSON save/load (mode 0600), versioned manifest with account ID, region, per-service coverage, pagination completion, unit tested
- [x] `internal/discover` — read-only AWS discovery (EC2, EBS, RDS, S3, EKS, Lambda, IAM) with injected mockable clients, pagination on every paginated API, caller context + 2-min per-service timeout, nil-safe conversion, no credential CLI flags (issue #4), unit tested (nil fields, pagination, access denied, partial failures)
- [x] `internal/score` — classifies resources as portable/rework/rewrite using map/domovyk-map.yaml, unit tested (bug fix: classify match-length comparison)
- [x] `internal/price` — v0.2 price comparison, unit tested (FOCUS parsing, FX, no_price, provenance, atomic write 0600)
- [x] `internal/generate` — v0.3 OpenTofu HCL generation (`hcloud_*` / `ovh_*`) for portable/rework resources, unit tested
- [x] `go build ./cmd/domovyk`, `go vet ./...`, `go test ./...` all pass
- [x] Git initialized: first commit 74d043a, .gitignore added, compiled binary excluded
- [x] GitHub Actions CI: `.github/workflows/ci.yml` build/vet/test on push and PR
- [x] `iam/domovyk-readonly-policy.json` — reviewed against actual SDK calls; `sts:GetCallerIdentity` added; still needs simulator check against real AWS
- [x] `map/domovyk-map.yaml` — 30 rows verified against MoveInfra EU equivalents + Hetzner live pricing page; `checked_on` 2026-11-23, status `reviewed`, `map/SCHEMA.md`
- [x] `map/prices.yaml` — Hetzner + OVHcloud rows with checked_on and status; OVHcloud compute/managed_db marked unverified pending page confirmation
- [x] `tf/reference-aws/` — OpenTofu reference setup: main.tf, variables.tf, vpc.tf, ec2.tf, rds.tf, s3.tf, alb.tf, README.md
- [x] `deploy/k8s/README.md` — container proof: kind + domovyk score job, real output captured
- [x] `docs/AWS-E2E-FINDINGS.md` — real AWS eu-north-1 discovery: 10 IAM roles, all unknown/unmapped

## Not done
- [ ] Run `discover` against live AWS with valid credentials (previous session credentials expired)
- [ ] IAM policy simulator check against real AWS (needs owner's AWS sandbox account)
- [ ] Run `price` against a real FOCUS billing export (synthetic testdata exists)
- [ ] Reference AWS `.tf` apply into a sandbox account (owner decision; spend is owner's decision)
- [ ] Portfolio drafts: owner review; lingonberry-ops has no data

## Next steps
1. Add valid AWS credentials, then run:
   `go run ./cmd/domovyk discover --region eu-central-1 --out inventory.json`
2. `go run ./cmd/domovyk score --inventory inventory.json --out score.json`
3. `go run ./cmd/domovyk generate --inventory inventory.json --out generated --target hetzner`
4. Review `generated/*.tf`, right-size, then `tofu apply` in a sandbox.
5. Run `domovyk price` with a real FOCUS 1.2 billing export.
6. Verify seed rows and IAM policy before claiming anything on a CV.

## Unverified items
- IAM policy simulator check pending (`iam/`)
- Threat model is a draft, not reviewed (`docs/THREAT-MODEL.md`)
- Reference AWS plan has no cost figure (`tf/reference-aws/README.md`)
- Any price number: none written; none should be added without `checked_on`

Verified by the owner (cited as verified in `docs/RESEARCH.md`): EU Data Act switching charges end 12 Jan 2027; parallel multi-cloud egress allowed at cost.
