# Roadmap: Domovyk

This roadmap prioritizes a dependable narrow workflow over broad provider/service coverage. All future milestones are gated by tests and pilot evidence.

## Gate 0 — Establish a trustworthy baseline (now)

- [ ] Add CI for Go format, tests, vet, and build.
- [ ] Run baseline build and tests; record actual results and supported Go version.
- [ ] Remove direct credential flags; use the standard AWS credential chain only.
- [ ] Audit discovery for pagination, region scope, partial failures, timeouts, and context cancellation.
- [ ] Write inventories atomically with restrictive permissions; avoid unsafe overwrite.
- [ ] Add tests for serialization, permissions, partial discovery, and error paths.
- [ ] Review IAM policy against the exact API calls and AWS documentation.
- [ ] Replace placeholder license and add the mapping-data license before publishing.

**Exit criteria:** reproducible CI green; no known credential exposure; discovery reports incomplete coverage honestly; security review findings tracked.

## Gate 1 — Validate the wedge (before executor work)

- [ ] Conduct 5–10 interviews with teams that operate AWS and are considering EU hosting.
- [ ] Secure one pilot partner and a safe non-production workload or sanitized state.
- [ ] Confirm the first source/target route and workload constraints.
- [ ] Compare the proposed workflow against MoveInfra, cloudexit, and native tools using the pilot's actual requirements.
- [ ] Define measurable acceptance tests and rollback limits with the pilot operator.

**Exit criteria:** a partner commits to a safe rehearsal and confirms the problem is worth solving. Otherwise revise the product direction.

## Gate 2 — Trustworthy discovery and planning

- [ ] Multi-region discovery with explicit service/region coverage and pagination status.
- [ ] Typed inventory schema with schema version and provenance.
- [ ] Mapping evidence: source URL, checked date, confidence/status, known limitations.
- [ ] Explainable plan with unknowns and blockers; no opaque portability score.
- [ ] Deterministic plan output and content digest.
- [ ] Redaction and local artifact security tests.

**Exit criteria:** pilot operator can independently review and reproduce the plan.

## Gate 3 — First executable vertical slice

Proposed first path: one stateless EC2 workload to one Hetzner Cloud server, pending pilot confirmation.

- [ ] Generate target infrastructure with OpenTofu; do not implement a new provisioning engine.
- [ ] Build a dry-run plan and bind approval to its digest.
- [ ] Execute one bounded provisioning/deployment workflow with a durable journal.
- [ ] Support safe resume after interruption and idempotent retries.
- [ ] Validate application health and operator-defined smoke tests.
- [ ] Keep DNS/traffic cutover manual, separate, and explicitly approved.
- [ ] Document and rehearse rollback in a disposable environment.

**Exit criteria:** repeated end-to-end rehearsal succeeds, including injected failures and rollback, with no unapproved side effects.

## Gate 4 — Pilot and hardening

- [ ] Run with the pilot in a controlled non-production environment.
- [ ] Measure coverage, interventions, duration, failures, recovery, and rollback.
- [ ] Fix all critical/high severity issues before broadening scope.
- [ ] Publish limitations and a reproducible operator guide.
- [ ] Add a second workload only when the first is reliable.

## Deferred

- Multi-cloud source support, broad service catalog, automatic DNS cutover, databases, EKS, Lambda rewrites, hosted control plane, compliance claims, autonomous production cutover, and generic provider plugin marketplace.

## Release rule

A milestone is not complete because code exists. It is complete only when its acceptance tests pass, limitations are documented, and a clean operator can reproduce the workflow.
