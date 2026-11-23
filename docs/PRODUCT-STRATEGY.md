# Product strategy

Status: proposed direction, not a validated product-market fit claim.

## Positioning

**Domovyk is a customer-controlled, auditable workflow for moving a narrowly scoped workload off AWS to an EU cloud—and proving the result is safe.**

It is not another cloud price calculator, portability score, generic Terraform generator, or migration consultancy marketplace. Existing inventory, IaC, transfer, database, and observability tools should be composed where possible.

## Who it serves first

- Small engineering teams operating AWS workloads that want an EU-hosted alternative.
- Teams without a dedicated cloud-migration department.
- Operators who require local execution, inspectable plans, explicit approvals, and a way to stop or resume safely.

Initial customer and workload assumptions remain hypotheses until interviews and a pilot confirm them.

## Product promise

For every supported migration, Domovyk must show:
1. What it found, what it did not inspect, and which permissions were missing.
2. The exact proposed changes, costs/risks where evidence exists, and the source of each claim.
3. Which actions are read-only, reversible, destructive, or irreversible.
4. What was executed, by whom, when, with what inputs, and with what result.
5. Whether target validation passed, whether cutover happened, and what rollback can and cannot restore.

## Initial wedge

Prove one safe, repeatable path before adding service breadth:

**One stateless EC2 workload → one Hetzner Cloud server**, using customer-owned credentials and infrastructure-as-code. First pilot excludes databases, EKS, Lambda, multi-account orchestration, complex networking, and persistent local state. Keep DNS cutover manual and explicitly approved.

This is a proposed vertical slice, not a claim that the current code supports it. Confirm provider APIs, workload shape, and customer demand with a pilot before committing to implementation.

## Differentiation

- Customer-run and local-first; no Domovyk control plane required.
- Workflow durability: journaled steps, resumability, idempotency, and checkpoints.
- Approval gates before infrastructure writes, data deletion, and traffic switching.
- Evidence-backed validation and explicit rollback limits.
- Composes mature tools instead of replacing OpenTofu, rclone, database replication, or monitoring.

MoveInfra already markets end-to-end cloud migration. Therefore “we migrate AWS to EU” is not a sufficient differentiator. Domovyk must demonstrate superior operator control, transparency, reproducibility, or a narrower workflow that customers actually prefer.

## Product principles

- **Fail closed:** unknown state, missing permission, stale evidence, or failed validation blocks progression.
- **No silent scope gaps:** every discovery run reports regions, services, permissions, pagination, and exclusions.
- **Plan before apply:** generate an immutable plan; require approval bound to its digest.
- **No hidden side effects:** read-only is the default; write capabilities are separate and opt-in.
- **Recoverability is designed in:** every mutating step defines retry, resume, compensation, and irreversibility.
- **Evidence over confidence:** never convert unverified mappings or prices into definitive recommendations.
- **Customer owns data and keys:** local artifacts, restrictive file permissions, documented retention.

## Go/no-go gates

Do not build the migration executor until:
- 5–10 target-buyer interviews identify a repeated, expensive problem.
- At least one pilot partner supplies a safe non-production workload or sanitized infrastructure state.
- The partner confirms the proposed first workload and accepts the manual approval model.
- A written runbook identifies rollback limits and an operator-owned recovery path.

Proceed only with the single workload and provider combination validated by the pilot. If no partner confirms the problem, revise positioning rather than expanding the feature list.

## Success measures

For the first pilot, measure:
- Inventory coverage and explicit unknowns (not a single opaque score).
- Percentage of planned steps completed without undocumented manual intervention.
- Recovery after interruption at every checkpoint.
- Validation results before and after cutover.
- Time to execute, operator interventions, and rollback outcome.
- Zero unapproved destructive or traffic-switching actions.

No marketing claim should be made from a single successful demo.