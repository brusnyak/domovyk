# Architecture baseline

Status: proposed constraints for design and review, not implemented architecture. Keep the first implementation narrow; revisit decisions when pilot evidence requires it.

## Design goals

- Local-first execution with customer-owned cloud accounts and credentials.
- Explicit, reviewable plans; no hidden side effects.
- Durable progress across process restarts.
- Safe retries after ambiguous provider timeouts.
- Clear evidence for discovery coverage, validation, and recovery.
- Reuse mature provider tooling where possible.

## Component boundaries

1. **CLI:** command parsing, user interaction, explicit confirmation, and exit codes. It should not contain provider business logic.
2. **Discovery adapters:** read-only source-provider calls. Return typed observations plus coverage/errors/provenance; never silently treat partial results as complete.
3. **Inventory model:** versioned, deterministic representation of observed resources, with source account/region/service and observation time. Exclude secrets.
4. **Planner:** pure transformation from inventory plus explicit target constraints into a plan. No provider writes or hidden network access.
5. **Plan artifact:** versioned canonical representation, stable digest, tool/version metadata, assumptions, unknowns, blockers, and per-step risk/reversibility.
6. **Approval record:** binds operator approval to the exact plan digest, target, scope, and timestamp. Any changed digest requires fresh approval.
7. **Executor:** runs only approved, bounded steps through narrow adapters. Separate read and write capabilities.
8. **Journal/store:** durable local record of step intent and result. Persist intent before mutation and result after mutation; restrictive permissions and atomic updates.
9. **Validator:** target-side infrastructure and application checks. A resource existing is not sufficient evidence of application health.
10. **Cutover controller:** separate phase and approval boundary. Initial pilot keeps DNS/traffic changes manual.

## State model

Proposed states:

`DISCOVERED → PLANNED → APPROVED → APPLYING → VALIDATING → READY_TO_CUTOVER → CUTOVER_APPROVED → CUTOVER → VERIFIED`

Any stage may enter `BLOCKED` or `FAILED`. Recovery must be explicit and journal-backed. Never transition forward after an unmet precondition.

State transitions should be implemented as validated functions over persisted state, not scattered CLI conditionals. Invalid transitions return errors and leave the prior durable state intact.

## Execution step contract

Each step needs:
- Stable step ID and schema/version.
- Declared resource scope and read/write classification.
- Preconditions and required permissions.
- Canonical input digest and bounded output.
- Timeout and retry policy.
- Idempotency/reconciliation strategy.
- Journal checkpoints before and after external mutation.
- Postcondition/validation probe.
- Recovery classification: retry, resume, compensate, manual recovery, or impossible.
- Explicit data-loss and irreversibility notes.

A timeout is an ambiguous outcome, not proof that the provider did nothing. Reconcile provider state before retrying.

## Provider boundary

Use narrow interfaces around required capabilities, not a generic cloud abstraction. Keep source discovery and target provisioning separate. For the proposed pilot, investigate whether OpenTofu plus provider APIs covers the target workflow before writing custom provisioning code.

Do not build a provider plugin marketplace. Do not promise portability across arbitrary services. Provider-specific limitations should remain visible in the plan.

## Persistence and artifacts

Start with a versioned local journal suitable for a single operator and one migration. Define corruption handling, locking/concurrent-run behavior, backup/export, and schema migration before implementation. Never persist credentials or secret values. Avoid a database dependency until concurrent operators or query requirements justify it.

## Security boundaries

- Discovery identity and execution identity are separate.
- No raw credentials in CLI arguments, environment dumps, logs, plans, or journal.
- Redact provider errors and subprocess output.
- Validate artifact paths and permissions; do not follow unsafe symlinks when writing sensitive files.
- Require explicit confirmation for every write-capable run and a separate confirmation for cutover.
- No production mutation support in the initial vertical slice.

## Open decisions (do not guess)

- Exact pilot workload and target provider, pending customer validation.
- Local journal format and locking mechanism.
- Whether execution invokes OpenTofu as a subprocess or uses its automation API.
- Exact approval UX and non-interactive confirmation policy.
- Which validation probes are workload-specific.
- Backup/data transfer method and recovery point objectives, if persistent data enters scope.

Resolve these with a small ADR and testable acceptance criteria when evidence exists. Avoid premature generalization.
