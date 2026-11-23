# PROGRESS: domovyk (for Yegor)

Status: pre-alpha → in-development. Scope: v0.0 to v0.2. Issue #4 (reliability gate) complete; build, vet, and tests pass.

## Done
- [x] `README.md`, `LICENSE` (full Apache-2.0 text), `go.mod` (`OWNER` placeholder)
- [x] `docs/SPEC.md`, `docs/THREAT-MODEL.md` (draft), `docs/RESEARCH.md`, `docs/ROADMAP.md`, `docs/PRICE.md`
- [x] `cmd/domovyk/main.go` — compiled, `go build` and `go vet` pass; CLI with `discover`, `score`, `price`
- [x] `internal/inventory` — Resource + Inventory types, JSON save/load (mode 0600), versioned manifest with account ID, region, per-service coverage, pagination completion, unit tested
- [x] `internal/discover` — read-only AWS discovery (EC2, EBS, RDS, S3, EKS, Lambda, IAM) with injected mockable clients, pagination on every paginated API, caller context + 2-min per-service timeout, nil-safe conversion, no credential CLI flags (issue #4), unit tested (nil fields, pagination, access denied, partial failures)
- [x] `internal/score` — classifies resources as portable/rework/rewrite using map/domovyk-map.yaml, unit tested (bug fix: classify match-length comparison)
- [x] `internal/price` — v0.2 price comparison, unit tested (FOCUS parsing, FX, no_price, provenance, atomic write 0600)
- [x] `go build ./cmd/domovyk`, `go vet ./...`, `go test ./...` all pass
- [x] Git initialized: first commit 74d043a, .gitignore added, compiled binary excluded
- [x] `iam/domovyk-readonly-policy.json` — reviewed against actual SDK calls; `sts:GetCallerIdentity` added; still needs simulator check against real AWS
- [x] `map/domovyk-map.yaml` (30 rows, all unverified), `map/SCHEMA.md`
- [x] `map/prices.yaml` (Hetzner + OVHcloud rows with checked_on and status; OVHcloud compute/managed_db marked unverified pending page confirmation)
- [x] `tf/reference-aws/README.md` (plan only)
- [x] `deploy/k8s/README.md` (container proof: kind + domovyk score job, real output captured)
- [x] `docs/AWS-E2E-FINDINGS.md` (real AWS eu-north-1 discovery: 10 IAM roles, all unknown/unmapped)

## Not done
- [ ] Run `discover` against live AWS with valid credentials (previous session credentials expired)
- [ ] IAM policy simulator check against real AWS (needs owner's AWS sandbox account)
- [ ] Verify seed rows in `map/domovyk-map.yaml` against primary sources; set `checked_on` and `status`
- [ ] GitHub and HN search for existing tools (required before v0.1) — see issue #3 / PR #2
- [ ] Price-list source decision (bundle vs fetch) noted in `docs/SPEC.md`
- [ ] Full Apache-2.0 attribution boilerplate in source files
- [ ] Reference AWS `.tf` files (v0.0 item, sandbox account only)
- [ ] Portfolio drafts: owner review; lingonberry-ops has no data

## Next steps
1. Add valid AWS credentials (`~/.aws/credentials` or `AWS_ACCESS_KEY_ID`/`AWS_SECRET_ACCESS_KEY`), then run:
   `go run ./cmd/domovyk discover --region eu-central-1 --out inventory.json`
2. `go run ./cmd/domovyk score --inventory inventory.json --out score.json`
3. Run `domovyk price` with a real FOCUS 1.2 billing export.
4. Verify seed rows and IAM policy before claiming anything on a CV.

## Unverified items
- 30 seed mapping rows, from general knowledge (`map/domovyk-map.yaml`)
- IAM policy simulator check pending (`iam/`)
- Threat model is a draft, not reviewed (`docs/THREAT-MODEL.md`)
- Reference AWS plan has no cost figure (`tf/reference-aws/README.md`)
- Any price number: none written; none should be added without `checked_on`

Verified by the owner (cited as verified in `docs/RESEARCH.md`): EU Data Act switching charges end 12 Jan 2027; parallel multi-cloud egress allowed at cost.
