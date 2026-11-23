# Reliability and safety requirements

Status: engineering baseline proposal. These are acceptance requirements, not current capabilities.

## Safety model

Domovyk progresses through explicit states:

`DISCOVERED → PLANNED → APPROVED → APPLYING → VALIDATING → READY_TO_CUTOVER → CUTOVER_APPROVED → CUTOVER → VERIFIED`

Any state may transition to `BLOCKED` or `FAILED`. Recovery is explicit; never silently skip a failed step. Cutover and destructive operations require separate approvals.

## Required invariants

1. **Read-only by default.** Discovery credentials must not grant write access. Execution credentials are a separate opt-in role and are not loaded during discovery.
2. **Approval binds to exact plan.** Approval records the plan digest, target, scope, operator, and timestamp. Any plan change invalidates approval.
3. **Durable journal.** Before and after each external mutation, persist step ID, inputs digest, start/end time, outcome, provider request/reference IDs, and recovery instructions.
4. **Idempotent execution.** A retry must detect whether the intended resource already exists and reconcile safely; never blindly repeat a non-idempotent action.
5. **Fail closed.** Missing permissions, incomplete pagination, unknown resource types, stale mappings, drift, and failed preconditions block apply.
6. **Secrets stay out of artifacts.** Never serialize credentials, tokens, private keys, secret values, or presigned URLs. Redact errors and command output.
7. **Restrictive local storage.** Inventory, plans, journals, and logs use owner-only permissions where supported. Do not overwrite existing files without an explicit safe policy.
8. **Rollback is honest.** Each step declares whether rollback is automatic, compensating, manual, or impossible. Data-loss windows and irreversible actions are shown before approval.

## Step contract

Every executable step must declare:
- Stable ID and version.
- Preconditions and required permissions.
- Read/write classification and affected resource identifiers.
- Inputs and outputs schemas.
- Idempotency strategy.
- Timeout, retry policy, and rate-limit handling.
- Checkpoint/journal behavior.
- Validation probe and expected result.
- Compensation/rollback procedure and limitations.

A step without a tested recovery story cannot run unattended.

## Discovery acceptance

A discovery run is trustworthy only when it records:
- Account identity and timestamp.
- Explicit region set and services attempted.
- Per-service pagination completion and resource counts.
- Access-denied and unsupported-operation results.
- Excluded resource classes and reasons.
- A coverage summary that distinguishes zero resources from failed inspection.

Do not label an inventory complete if any requested region/service was not fully inspected.

## Execution acceptance

Before enabling any write operation:
- Plan is deterministic for identical inputs and tool versions.
- Plan can be rendered and reviewed without credentials.
- Apply requires explicit confirmation of the plan digest.
- State transitions and journal survive process termination.
- Resume behavior is tested after interruption before, during, and after each mutation.
- Retry behavior is tested for timeouts where the provider may have completed the request.
- Validation checks target health and application-level behavior, not merely resource existence.
- Cutover is separate from provisioning and requires a second approval.
- Rollback is exercised in a disposable environment and reports any non-reversible data changes.

## Minimum test strategy

- Unit tests for parsers, mapping, state transitions, redaction, and plan hashing.
- Fixture tests for AWS API pagination, partial failures, throttling, and access denial.
- Contract tests for each provider adapter.
- Integration tests in disposable cloud accounts/projects with budget limits.
- Fault-injection tests for process kill, network timeout, duplicate requests, and provider-side partial success.
- End-to-end pilot rehearsal including rollback.

## Operational safety

- No cloud write calls in unit tests.
- No live-account test without explicit owner approval and a disposable scope.
- Every integration test has cleanup and budget controls.
- CI must run formatting, tests, vet, and build on every pull request.
- Security/threat model and IAM policy must be reviewed before real customer use.

## Current known gaps

The current repository is pre-alpha. README and progress notes say discovery was run once, but build/tests have not been verified. The IAM policy and threat model are drafts; mapping rows are unverified; credential flags are present in the CLI and need removal or a documented secure design. These gaps block any claim of production readiness.
