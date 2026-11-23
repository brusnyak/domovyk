# PROGRESS: Domovyk

Status: pre-alpha; reliability baseline is being established. No production migration capability.

## Completed in this branch

- [x] Repositioned project around customer-controlled, auditable migration workflows.
- [x] Defined a narrow proposed first workload: stateless EC2 to Hetzner Cloud, pending pilot confirmation.
- [x] Added product strategy and explicit validation gates.
- [x] Added reliability/safety requirements for approvals, journaling, idempotency, recovery, and evidence.
- [x] Replaced feature-count roadmap with gated reliability milestones.
- [x] Added GitHub Actions CI for gofmt, tests, vet, and build.

## Not verified yet

- [ ] CI has not run on this branch; build/tests are not claimed to pass.
- [ ] Discovery pagination is incomplete or unverified; service coverage and region coverage are not complete.
- [ ] Discovery contains unchecked pointer dereferences and uses context.Background.
- [ ] CLI exposes direct access-key/secret-key/session flags; remove these and use the standard credential chain only.
- [ ] IAM policy is unreviewed; do not use it for production/customer accounts.
- [ ] Mapping rows remain unverified.
- [ ] Threat model remains a draft.
- [ ] License texts and mapping-data license need verification.
- [ ] No executor, plan approval, journal, resume, validation, cutover, or rollback implementation exists.

## Current development order

1. Merge/review the product and CI foundation only after CI and docs review.
2. Establish build/test baseline and fix discovery safety: pagination, nil handling, cancellation, partial-coverage reporting, atomic output, and restrictive permissions.
3. Review IAM policy against actual SDK calls and least-privilege AWS documentation.
4. Validate the first workload with buyer interviews and a safe pilot.
5. Implement execution only after a pilot confirms the workload and acceptance criteria.

## Release rule

A checkbox means evidence exists. Code presence alone is not completion. Record commands, CI results, test scope, and limitations. Never run live cloud mutations without explicit owner approval and a disposable target.
