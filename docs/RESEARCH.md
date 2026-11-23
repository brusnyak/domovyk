# RESEARCH: domovyk

Labels: **VERIFIED** (checked against cited source), **UNVERIFIED** (not checked; do not cite as fact).

## Current market finding — 2026-11-23

The previous market statement below is obsolete: direct competitors now include **MoveInfra** and **EscapeCloud/cloudexit**. MoveInfra publicly offers open-source cloud-exit tools and claims AWS/GCP/Azure/DigitalOcean to EU-provider migration. EscapeCloud/cloudexit is an open-source AWS/Azure exit-readiness CLI.

See the full comparison, AWS-to-EU service matrix, existing-tool alternatives, product decision and go/no-go gate in [COMPETITIVE-LANDSCAPE-2026-11.md](COMPETITIVE-LANDSCAPE-2026-11.md).

**Decision:** do not build a generic assessment, price calculator, service map or IaC generator. Conditional exploration only for a customer-run, open-source execution workflow that composes existing tools and proves reliable checkpoints, validation and recovery. MoveInfra is a direct competitor; differentiation and demand remain unvalidated.

## Regulatory context

- **VERIFIED:** Article 29 of EU Regulation 2023/2854 prohibits cloud switching charges from 12 January 2027. Until then, charges may not exceed directly linked switching costs.
- **VERIFIED:** AWS publicly offers a data-transfer-out programme for customers moving off AWS, subject to eligibility and process. Egress savings are not a defensible standalone moat.
- **UNVERIFIED:** Domovyk's suitability for DORA exit evidence; requires a precise controls/evidence mapping and legal review.

Sources:
- https://eur-lex.europa.eu/eli/reg/2023/2854/oj/eng
- https://digital-strategy.ec.europa.eu/en/factpages/data-act-explained
- https://aws.amazon.com/blogs/aws/free-data-transfer-out-to-internet-when-moving-out-of-aws/

## Mapping seed data

The 30 rows in map/domovyk-map.yaml are unverified seed data. Do not use them to make customer-facing portability claims until each row has a primary source, destination scope, confidence, and checked_on date.

## Open engineering questions

1. Which mappings hold up against a representative real account and Terraform/OpenTofu state?
2. Does FOCUS 1.2 provide enough detail for accurate service-level cost comparisons?
3. How should price data be sourced, refreshed, licensed and marked stale?
4. Can AWS inventory safely avoid sensitive fields and secret/user-data exposure?
5. What workload and destination does a prospective pilot actually require?

## Next gate

Before broad implementation, complete customer validation and select one non-critical migration route. Required: 5–10 buyer interviews, one partner willing to test with sanitized state/billing or a safe workload, a documented competitor gap, and explicit acceptance criteria for validation and rollback.
