# domovyk

Open-source planner for moving AWS workloads to EU providers (OVHcloud, Hetzner). A Go CLI that runs on your own machine with a read-only IAM role.

**Status: pre-alpha, nothing built.** This repo holds the skeleton, the spec and the research notes. No subcommand does real work yet.

## What it is meant to do

1. `discover`: read-only inventory of an AWS account into `inventory.json`.
2. `score`: classify each service as `portable`, `rework` or `rewrite` using the open mapping data (`map/`).
3. `price`: compare your AWS billing export (FOCUS 1.2) with OVHcloud and Hetzner price lists. Every price row carries a `checked_on` date.
4. `generate`: OpenTofu modules for the target. Roadmap only, not in scope yet.
5. `cut over` and `operate`: runbook with rollback and a monthly drift check. Roadmap only.

## Honest limits

- It does not move anything. v0.1 produces a sample exit plan, nothing more.
- It will say no to Lambda-heavy accounts. Those are `rewrite`, not a lift-and-shift.
- No compliance claims. This is not legal advice. The EU Data Act text is the reference.
- Prices are point-in-time. Read the `checked_on` date before trusting a number.
- The 30 mapping rows in `map/` are unverified seed data written from general knowledge.
- The IAM policy in `iam/` is unreviewed.

## Roadmap

In scope now:
- v0.0: skeleton, README, read-only IAM policy (unreviewed), 30 seed mapping rows (unverified), reference AWS setup plan.
- v0.1: discover, score, sample exit plan.
- v0.2: price from a FOCUS export.

Roadmap only: v0.3 generate, v0.4 measured cut-over, v0.5 drift job and plugin contract tests. See `docs/ROADMAP.md`.

## Docs

- `docs/SPEC.md`: architecture and six stages
- `docs/THREAT-MODEL.md`: draft threat model
- `docs/RESEARCH.md`: competitors and open questions, with verified and unverified labels
- `docs/ROADMAP.md`, `PROGRESS.md`

## Build (not run yet)

```sh
go build ./cmd/domovyk
```

Replace `OWNER` in `go.mod` before publishing.

## Licence

Code: Apache-2.0. `LICENSE` currently holds a placeholder; the full text must be added before release.
Mapping data in `map/`: intended CC-BY-4.0. The licence file for the data is not added yet.
