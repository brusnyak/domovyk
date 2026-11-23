# Agent instructions

These rules apply to every human and automated contributor working in this repository.

## 1. Product boundary

Domovyk is pre-alpha. Its proposed direction is a customer-controlled, auditable workflow for moving a narrowly scoped workload off AWS to an EU cloud.

Read these before substantial changes:
- `README.md`
- `docs/PRODUCT-STRATEGY.md`
- `docs/RELIABILITY.md`
- `docs/ROADMAP.md`
- `PROGRESS.md`

The initial route under consideration is one stateless EC2 workload to one Hetzner Cloud server. This is a hypothesis pending buyer validation and a pilot, not a supported capability.

Do not expand scope to databases, EKS, Lambda rewrites, multi-account orchestration, automatic DNS cutover, or additional providers without an explicit product decision backed by pilot evidence. Prefer composing mature tools over reimplementing them.

## 2. Non-negotiable safety rules

- Default to read-only. Discovery must not perform cloud writes.
- Never run a real cloud mutation, production cutover, deletion, or destructive cleanup.
- Never request, print, commit, log, or place credentials, tokens, private keys, secret values, or customer data in fixtures.
- Never add CLI flags that accept raw secrets. Use supported SDK credential providers/profile mechanisms.
- Never weaken IAM permissions to make tests pass. Document the exact API operation and least-privilege permission needed.
- Unknown, incomplete, stale, or ambiguous state must block progression; do not guess.
- Treat plans as immutable approval artifacts. Any plan change invalidates approval.
- Separate infrastructure provisioning from traffic cutover. Cutover requires a distinct explicit approval.
- Do not claim rollback is possible unless the specific operation has a tested recovery/compensation path.
- Do not make network calls or cloud writes in unit tests. Integration tests require a disposable environment, explicit owner approval, budget guardrails, and cleanup.
- Do not overwrite user files or existing artifacts without an explicit safe policy.

If a task appears to require violating a rule, stop and report the blocker instead of improvising.

## 3. Engineering rules

- Keep changes small, focused, and reviewable. One logical concern per pull request.
- Follow existing Go package structure and idioms. Avoid new dependencies unless the benefit and maintenance cost are justified.
- Thread caller-provided `context.Context` through I/O and SDK calls. Do not introduce `context.Background()` inside library code.
- Handle nil pointers and partial API responses defensively.
- Paginate every list/describe API that supports pagination; expose whether pagination completed.
- Distinguish an empty result from an API failure, access denial, unsupported operation, or incomplete scan.
- Return actionable errors with operation context; redact sensitive provider output.
- Use deterministic serialization where output is hashed, compared, or approved.
- Use atomic file writes and restrictive permissions for local artifacts where supported.
- Keep public interfaces minimal. Avoid speculative frameworks and generic plugin abstractions before a second proven use case.

## 4. Required workflow

1. Read the relevant docs and inspect the existing implementation before editing.
2. State the narrow task and acceptance criteria in the issue/PR.
3. Check for an existing issue or implementation; do not duplicate work.
4. Implement the smallest safe change.
5. Add or update tests for success, failure, partial response, and relevant security boundaries.
6. Run formatting, tests, vet, and build where possible.
7. Inspect the final diff for secrets, accidental scope expansion, and unrelated changes.
8. Update documentation and `PROGRESS.md` only when evidence justifies the status change.
9. Open a focused PR. Do not merge unless explicitly authorized.

## 5. Definition of done

A task is done only when:
- Acceptance criteria are met and backed by tests or other reproducible evidence.
- Relevant tests, `go vet`, and build have been run; results and exact commands are reported.
- Failure paths and safety boundaries are covered.
- Docs describe actual behavior, not intended future behavior.
- No secrets or customer-sensitive data are present.
- Limitations and unverified behavior are stated plainly.

If a command cannot be run, mark it **NOT RUN** and explain why. Never infer success from code inspection or a CI workflow file.

## 6. Agent coordination

- Assign one owner per file or tightly scoped task to avoid conflicting edits.
- Do not let multiple agents modify the same files or branch concurrently.
- Agents must report: files changed, design decisions, commands run, test results, remaining risks, and commit/PR reference.
- The integrator reviews every diff and resolves conflicts; agent output is not automatically trusted.
- Do not create parallel implementations of the same feature.
- Keep task dependencies explicit. Discovery reliability precedes planning; planning and approval precede execution; execution rehearsal precedes any pilot.

## 7. Current priority

Work from `docs/AGENT-HANDOFF.md`. Reliability issue #4 is the first implementation task. Do not begin an executor until Gate 1 in `docs/ROADMAP.md` is satisfied.

## 8. Reporting format

For every completed task, report:
- **Status:** verified / unverified / blocked
- **Changes:** paths and behavior
- **Validation:** exact commands and outcomes
- **Risks:** remaining gaps
- **Next:** one concrete dependent task
