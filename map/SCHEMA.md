# map/ schema

File: `domovyk-map.yaml`. Licence intended: CC-BY-4.0.

Top-level key: `rows` (list). Each row:

| Field | Type | Meaning |
|---|---|---|
| `id` | string | Stable ID, `aws-<service>-<slug>` |
| `aws_service` | string | AWS service name as written in AWS docs |
| `class` | enum | `portable` (lift and shift to a VM or container), `rework` (same function, different service or config), `rewrite` (no direct equivalent; redesign) |
| `typical_eu_target` | string | Common EU target. Not a recommendation. |
| `caveat` | string | What breaks or costs extra |
| `checked_on` | date or empty | Date a human last checked the row against primary sources. Empty = never checked. |
| `status` | enum | `unverified`, `reviewed`, `stale` |

Rules:
- A row with empty `checked_on` must have `status: unverified`.
- Changing `class` needs a pull request with a source link in the PR body.
- No price in this file. Prices live in the price step, with their own `checked_on`.
