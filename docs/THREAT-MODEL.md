# Threat model: domovyk (DRAFT)

Status: draft, not reviewed by anyone outside the project.

## Assets

- AWS account metadata: account IDs, resource names, tags, network layout, IAM role and user names.
- Billing export (FOCUS 1.2) supplied by the customer.
- Credentials used by the tool (the operator's own AWS profile or role).
- Output files: `inventory.json`, score reports, price tables, exit plans.

## Trust boundaries

- The operator's machine: the only place the tool runs.
- AWS APIs: the tool calls read-only describe and list endpoints with the operator's credentials.
- Price list sources (OVHcloud, Hetzner): how they are read is an open decision (bundled snapshot or live fetch).
- Mapping data repository (`domovyk-map`): reviewed by pull request.

## Controls

1. Read-only IAM. The policy in `iam/domovyk-readonly-policy.json` allows only describe, list and a few bucket-metadata reads. It is UNREVIEWED. No write actions. No object reads.
2. Local-only. No Domovyk server, no telemetry, no account. Customer data is not sent to any project-operated endpoint.
3. Credentials come from the standard AWS credential chain. The tool does not store, print or write them.
4. Outputs are written locally. `inventory.json` is written with mode 0600.
5. Pinned dependencies. Go modules are pinned by `go.sum` once added.
6. Price rows carry `checked_on`, so stale numbers are visible.

## What leaves the machine

- Customer data sent to a project-operated endpoint: none. There is no such endpoint.
- Traffic to AWS APIs: yes, by design, under the operator's credentials. This is the discovery itself.
- Traffic to price list sources, if they are fetched live: yes. These requests carry no customer data, but the provider sees the operator's IP address.
- Mapping data: fetched or copied by the operator. Not customer data.

## Threats and mitigations

| Threat | Mitigation | Status |
|---|---|---|
| Over-broad IAM policy lets the tool read secrets | Explicit action list, no `GetObject`, no `GetPasswordData`, no `lambda:GetFunction` (returns a code download URL), no `ec2:DescribeInstanceAttribute` (can return user data) | UNREVIEWED |
| Operator runs the tool with admin credentials | Document the read-only policy; a preflight check is roadmap | Open |
| `inventory.json` leaks account structure | File mode 0600; warn in the README; do not commit | Planned |
| Malicious or wrong mapping row changes a verdict | Pull-request review; `checked_on`; `status` field in rows | Planned |
| Stale prices lead to a wrong cost decision | `checked_on` on every row; the tool reports the oldest date | Planned |
| Supply chain (Go modules) | `go.sum`, pinned versions | Planned |
| Over-claiming compliance | Documentation rule: no compliance claims; legal text is the reference | Policy |

## Out of scope

- Insider threats at the operator's own machine.
- Security of the target cloud after cut-over.
- Legal sufficiency under the EU Data Act or GDPR.
