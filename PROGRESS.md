# PROGRESS: domovyk (for Yegor)

Status: in-development (pre-alpha). Scope: v0.0 to v0.3 implemented. Build, vet, and tests pass; GitHub Actions CI on push/PR.

## Done
- [x] `README.md`, `LICENSE` (full Apache-2.0 text), `go.mod` (`OWNER` placeholder)
- [x] `docs/SPEC.md`, `docs/THREAT-MODEL.md` (draft), `docs/RESEARCH.md`, `docs/ROADMAP.md`, `docs/PRICE.md`
- [x] `cmd/domovyk/main.go` — compiled, `go build` and `go vet` pass; CLI with `discover`, `score`, `price`, `generate`
- [x] `internal/inventory` — Resource + Inventory types, JSON save/load (mode 0600), versioned manifest with account ID, region, per-service coverage, pagination completion, unit tested
- [x] `internal/discover` — read-only AWS discovery (EC2, EBS, RDS, S3, EKS, Lambda, IAM) with injected mockable clients, pagination on every paginated API, caller context + 2-min per-service timeout, nil-safe conversion, no credential CLI flags (issue #4), unit tested (nil fields, pagination, access denied, partial failures). **Bug fix: EBS volumes were tagged `Service: "ec2"` and scored `portable`; now `Service: "ebs"` → `rework`. Regression test added.**
- [x] `internal/score` — classifies resources as portable/rework/rewrite using map/domovyk-map.yaml, unit tested (bug fix: classify match-length comparison; EBS regression test).
- [x] `internal/generate` — v0.3 OpenTofu HCL generation (`hcloud_*` / `ovh_*`) for portable/rework resources, unit tested. Now emits `provider.tf` with `required_providers`, provider block, API-token + region variables so `tofu apply` can run after review. Uses actual AMI from inventory. Hardcoded `0.044` volume price removed (was not from an official page per `map/prices.NOTES.md`).
- [x] `go build ./cmd/domovyk`, `go vet ./...`, `go test ./...` all pass
- [x] Git initialized: first commit 74d043a, .gitignore added, compiled binary excluded
- [x] GitHub Actions CI: `.github/workflows/ci.yml` build/vet/test on push and PR
- [x] `iam/domovyk-readonly-policy.json` — reviewed against actual SDK calls; `sts:GetCallerIdentity` added; simulator check pending live AWS
- [x] `map/domovyk-map.yaml` — 30 rows, each carrying `checked_on: "2026-11-23"` and `status: reviewed` per `map/SCHEMA.md`
- [x] `map/prices.yaml` — Hetzner + OVHcloud rows with checked_on and status; OVHcloud compute/managed_db marked unverified pending page confirmation
- [x] `tf/reference-aws/` — OpenTofu reference setup: main.tf, variables.tf, vpc.tf, ec2.tf, rds.tf, s3.tf, alb.tf, README.md
- [x] `deploy/k8s/README.md` — container proof: kind + domovyk score job, real output captured
- [x] `docs/AWS-E2E-FINDINGS.md` — real AWS eu-north-1 discovery: 10 IAM roles, all unknown/unmapped
- [x] README + SPEC + RESEARCH reconciled: verified/unverified contradiction resolved; map rows cited as source of truth; architecture diagram consolidated
- [x] `~/.agent-skills/core/domovyk-migration/` skill drafted, validated, synced (repo untouched, no push)

- [x] IAM policy simulator check against real AWS (account `004406188742`, read-only): all write actions (ec2:RunInstances/TerminateInstances, s3:PutObject/DeleteBucket, lambda:CreateFunction, rds:CreateDBInstance) → explicitDeny under domovyk's policy; `iam:CreateUser` allowed via the principal's pre-existing AccountFullAccessRole, not domovyk's policy.
- [ ] Run `price` against a real FOCUS billing export (synthetic testdata exists; verified working on synthetic data)
- [ ] Reference AWS `.tf` apply into a sandbox account (owner decision; spend is owner's decision — will not run)
- [ ] Portfolio drafts: owner review; lingonberry-ops has no data

## Next steps
- Live `discover` against AWS eu-north-1 works (14 resources / 7 services, all `completed=true`) after writing short-lived creds to `~/.aws/credentials` (mode 0600); `aws login` stores creds in its CLI-only `login` cache, which the Go SDK does not read.
1. `go run ./cmd/domovyk discover --region eu-north-1 --profile default --out inventory.json` (verified live)
2. `go run ./cmd/domovyk score --inventory inventory.json --out score.json` (verified live: 2 Lambda→rewrite, 12 IAM→unknown)
3. `go run ./cmd/domovyk generate --inventory inventory.json --out generated --target hetzner` (verified live: emits provider.tf with hetznercloud/hcloud, skips rewrite/unknown)
4. Review `generated/*.tf`, right-size, then `tofu apply` in a sandbox.
5. Run `domovyk price` with a real FOCUS 1.2 billing export.
6. Verify seed rows and IAM policy before claiming anything on a CV.

## Unverified items
- [ ] Run `price` against a real FOCUS billing export (synthetic testdata exists; verified working on synthetic data)
- [ ] Threat model is a draft, not reviewed (`docs/THREAT-MODEL.md`)
- [ ] Reference AWS plan has no cost figure (`tf/reference-aws/README.md`)
- [ ] Any price number: none written; none should be added without `checked_on`
- [ ] Portfolio drafts: owner review; lingonberry-ops has no data

Verified by the owner (cited as verified in `docs/RESEARCH.md`): EU Data Act switching charges end 12 Jan 2027; parallel multi-cloud egress allowed at cost.
