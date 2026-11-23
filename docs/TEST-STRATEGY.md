# Test and verification strategy

Status: required test plan. This document does not claim that current tests exist or pass.

## Test layers

### 1. Unit tests — every PR

No live network or cloud access.

- Context cancellation and timeout propagation.
- Nil/empty/partial SDK responses.
- Pagination tokens, repeated tokens, and incomplete pagination.
- Access denied, throttling, transient errors, and unsupported APIs.
- Coverage accounting: zero resources versus failed or skipped inspection.
- Deterministic inventory/plan serialization and digest stability.
- Approval invalidation when plan digest changes.
- Valid and invalid state transitions.
- Journal recovery from truncated/corrupt records.
- Idempotency and ambiguous-timeout reconciliation logic.
- Secret redaction and artifact permission checks.

Use table-driven tests and small fixtures. Keep fixtures sanitized and representative.

### 2. SDK/API contract tests

Use generated SDK stubs or HTTP fixtures to verify request pagination, region scope, retries, and response handling. Assert exact API operations and required IAM actions. Do not require user credentials.

### 3. Provider adapter tests

Test adapter contracts against a local fake or provider-supported test/dry-run interface. Verify idempotency, timeouts, partial success, and error normalization. Do not treat a successful plan as proof of apply safety.

### 4. Disposable integration environment — gated

Only after unit and contract tests pass:
- Explicit owner approval for the account/project and spend ceiling.
- Dedicated disposable account/project and least-privilege credentials.
- Unique resource prefix and strict allowlist of resource types/regions.
- Budget alert and hard time limit.
- Automated cleanup plus a documented manual cleanup path.
- No production data, credentials, or customer state.
- Capture provider request IDs and sanitized logs.

If a provider cannot guarantee a safe disposable scope, do not run the test.

### 5. Failure-injection and recovery

For each mutating step, simulate:
- Process termination before request, during request, and after provider success but before journal completion.
- Network timeout with provider-side success.
- Duplicate invocation and retry.
- Rate limiting and transient provider errors.
- Invalid/stale approval and plan drift.
- Journal write failure and disk-full behavior.
- Validation failure after provisioning.
- Failed or partial compensation.

The expected result must be deterministic: block, reconcile, resume, or require operator intervention. Never blindly replay an ambiguous mutation.

### 6. End-to-end rehearsal

Use one approved non-production workload only after customer validation. Rehearse inventory → plan → review → approval → provision/deploy → validation → manual cutover approval → verification → rollback. Record timings, operator interventions, evidence, and limitations.

## CI minimum

Every PR should run:
- `gofmt` check.
- `go test ./...`.
- `go vet ./...`.
- `go build ./cmd/domovyk`.

Add race detection where practical. Pin or explicitly manage Go/toolchain versions. CI must not require cloud credentials or perform cloud writes.

## Evidence reporting

For each task, report exact commands and results. Distinguish:
- **PASS:** command ran and passed.
- **FAIL:** command ran and failed; include concise error.
- **NOT RUN:** not executed; explain why.
- **NOT APPLICABLE:** explain why.

Do not call code tested because tests were added. Do not call CI green until the actual run is green. Record test scope and known blind spots in the PR.
