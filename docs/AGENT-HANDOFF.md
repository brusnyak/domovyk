# Agent handoff: next implementation cycle

This is the task queue for the current pre-alpha stage. It is intentionally ordered. Do not parallelize tasks that edit the same files or depend on unfinished interfaces.

## Mission

Make AWS discovery trustworthy before designing or implementing migration execution. Domovyk currently has no production migration executor. The proposed EC2-to-Hetzner path remains subject to customer validation.

## Global rules

Read `/AGENTS.md`, `docs/PRODUCT-STRATEGY.md`, `docs/RELIABILITY.md`, `docs/ARCHITECTURE-BASELINE.md`, and `docs/TEST-STRATEGY.md` first. Follow the repository safety rules. No real cloud writes, production cutovers, secrets, or customer data.

## Ordered tasks

### Task 0 — Establish baseline (integrator)

- Inspect CI workflow and Go version.
- Run formatting, tests, vet, and build locally or report why unavailable.
- Confirm CI actually ran on the PR branch.
- Record exact commands/results; do not fix unrelated issues in this task.

**Done when:** baseline is reproducible and all failures are recorded.

### Task 1 — Harden AWS discovery (existing GitHub issue #4)

Owner: one Go agent. Do not create a duplicate issue.

- Remove raw access-key, secret-key, and session-token CLI flags; use AWS SDK standard credential chain and profile support.
- Thread caller context through discovery and SDK calls; add bounded timeouts.
- Complete pagination for every list/describe operation that supports it.
- Make pointer handling nil-safe.
- Represent per-region/per-service attempted, completed, denied, unsupported, failed, and skipped states.
- Distinguish confirmed empty results from incomplete inspection.
- Write output atomically with restrictive permissions and explicit overwrite behavior.
- Add tests for every affected API path and error mode.
- Review IAM policy against exact SDK calls; list missing/unused permissions.

**Done when:** tests cover pagination, cancellation, nil fields, access denial, partial failures, and artifact safety; output never claims full coverage when inspection was incomplete.

### Task 2 — Define inventory coverage contract

Depends on Task 1 design findings.

- Add a versioned inventory schema with account identity, observation timestamp, region/service coverage, resource provenance, and schema version.
- Keep secrets out of inventory.
- Add deterministic serialization and compatibility tests.
- Document unsupported resource classes and coverage semantics.

**Done when:** fixtures prove zero resources and failed inspection are distinguishable; schema is documented and stable.

### Task 3 — Add fixture-driven discovery regression suite

Can proceed after Task 1 establishes response/coverage interfaces.

- Add sanitized fixtures for paginated responses, empty pages, missing fields, throttling, access denial, and partial region failure.
- Add regression tests for previously identified panic and pagination risks.
- Ensure tests use fake clients/HTTP fixtures and no live credentials.

**Done when:** the full suite is deterministic and runs in CI without cloud access.

### Task 4 — Validate customer wedge (product owner, not coding agent)

- Conduct 5–10 interviews with AWS operators considering EU hosting.
- Secure one safe non-production pilot or sanitized workload state.
- Confirm source, target, workload constraints, data sensitivity, downtime tolerance, and operator approval expectations.
- Decide go/no-go on the proposed stateless EC2 → Hetzner slice.

**Done when:** interview evidence and a pilot commitment are documented. If not validated, revise the product direction before executor work.

### Task 5 — Plan artifact and digest (blocked on Gate 1)

Only start after Task 4 passes.

- Define canonical versioned plan schema.
- Include assumptions, unknowns, blockers, resource scope, risk, reversibility, and provenance.
- Implement stable digest and tests for deterministic output.
- Ensure approval binds to exact digest and becomes invalid after any plan change.

**Done when:** independent runs on identical inputs produce identical digest; changed inputs invalidate approval.

### Task 6 — Local journal and state machine (blocked on plan contract)

- Define validated state transitions and durable journal schema.
- Persist intent before and result after each external mutation.
- Add locking/concurrent-run policy and corruption recovery.
- Test restart, duplicate invocation, ambiguous timeout, and disk failure.

**Done when:** fault-injection tests prove safe block/reconcile/resume behavior without blind retries.

### Task 7 — One bounded execution step (blocked on pilot + Tasks 5–6)

- Choose exactly one pilot-approved operation.
- Define preconditions, permissions, timeout, idempotency, validation, and recovery.
- Use OpenTofu/provider tooling where appropriate; do not create a generic provisioning engine.
- Keep all writes confined to a disposable non-production environment with explicit approval.

**Done when:** repeatable rehearsal passes, including injected failures and recovery.

### Task 8 — Validation and manual cutover rehearsal (blocked on Task 7)

- Add workload-specific application health checks.
- Keep DNS/traffic change manual and separate.
- Require a second approval bound to the validated deployment.
- Rehearse documented rollback and report non-reversible effects.

**Done when:** an operator can reproduce the run and understand exactly what rollback can and cannot restore.

## Parallelization guidance

Safe parallel work:
- Task 0 baseline review can run alongside customer interviews.
- Test fixture design can begin after the discovery interfaces are agreed.
- Product interviews do not require code changes.

Do not parallelize:
- Multiple agents editing discovery implementation.
- Inventory schema and discovery coverage implementation before agreeing on coverage semantics.
- Plan digest and approval implementation before canonical plan schema is reviewed.
- Executor work before Gate 1 and journal/state-machine design.

## Agent completion report

Every agent returns:
1. Task number and status (verified / unverified / blocked).
2. Files changed and concise behavior summary.
3. Exact commands run and outputs/status.
4. Tests added and what they cover.
5. Security/reliability risks remaining.
6. Commit SHA and PR URL.
7. One recommended next task.

The integrator reviews the diff and test evidence. Never merge agent branches automatically.
