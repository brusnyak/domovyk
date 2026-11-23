# Domovyk E2E Findings — AWS eu-north-1

Date: 2026-11-16
Account: `004406188742` (role `AccountFullAccessRole`)
Region: `eu-north-1` (Stockholm)

## 1. E2E command results

### `discover`
```
go run ./cmd/domovyk discover -region eu-north-1 -profile default -out /tmp/inventory2.json
```
**Failed at runtime** — the `default` profile holds expired short-lived session credentials (`ASIA…` keys). `aws sts get-caller-identity` returns `InvalidClientTokenId`. SSO is not configured (`~/.aws/config` references a `login_session` assumed-role ARN but no `sso_start_url`/`sso_region`), so the token cannot be refreshed unattended.

**However**, a successful discover ran earlier today (2026-11-16T02:43:20Z) against the same account/region and wrote real inventory to `/tmp/inventory.json`. That file is the source of truth below (it is real AWS API data, not "unknown"). `/tmp/inventory2.json` was not produced because the credentials expired before this session could re-run discover.

### `score`
```
go run ./cmd/domovyk score -inventory /tmp/inventory.json -out /tmp/score2.json
```
**Succeeded.** Scored all 10 resources. Result: `portable 0, rework 0, rewrite 0, unknown 10, total 10`.

All 10 resources scored `unknown` (source `unmapped`) because they are IAM roles and the mapping file `map/domovyk-map.yaml` has no IAM-role entries. This is an accurate mapping gap, not missing data.

### `price`
```
go run ./cmd/domovyk price -help
```
```
Usage of price:
  -focus string
    	FOCUS 1.2 billing export path (required)
```
The `price` subcommand exists (v0.2) but **requires `--focus <FOCUS 1.2 billing export path>`**, which is mandatory and was not provided. No FOCUS billing export exists in the repo or `/tmp`, so `price` could not be run and no AWS-vs-Hetzner/OVH cost comparison output was generated. To run it: obtain a FOCUS 1.2 CSV/Parquet billing export from AWS Billing and pass `--focus` to `domovyk price`.

---

## 2. Actual infrastructure findings (eu-north-1)

Source: `/tmp/inventory.json`, generated 2026-11-16T02:43:20Z, read via live AWS API earlier today.

**Service counts:** `iam_role: 10` — only IAM resources were discovered; no EC2, RDS, S3, VPC, Lambda, or other compute/resources were present in this account/region at discovery time.

**All 10 resources (IAM roles), state `active`:**

| # | Resource id | Service | Type | Region | State | ARN |
|---|---|---|---|---|---|---|
| 1 | AccountFullAccessRole | iam | role | eu-north-1 | active | arn:aws:iam::004406188742:role/managed/AccountFullAccessRole |
| 2 | AWSManagedAccountManagementAccessRole | iam | role | eu-north-1 | active | arn:aws:iam::004406188742:role/managed/AWSManagedAccountManagementAccessRole |
| 3 | AWSManagedAdvancedFeaturesActivationRole | iam | role | eu-north-1 | active | arn:aws:iam::004406188742:role/managed/AWSManagedAdvancedFeaturesActivationRole |
| 4 | AWSServiceRoleForAwsUserNotifications | iam | role | eu-north-1 | active | arn:aws:iam::004406188742:role/aws-service-role/notifications.amazonaws.com/AWSServiceRoleForAwsUserNotifications |
| 5 | AWSServiceRoleForBudgetsSpendLimitMember | iam | role | eu-north-1 | active | arn:aws:iam::004406188742:role/aws-service-role/spendlimit_member.budgets.amazonaws.com/AWSServiceRoleForBudgetsSpendLimitMember |
| 6 | AWSServiceRoleForIAMRoleManager | iam | role | eu-north-1 | active | arn:aws:iam::004406188742:role/aws-service-role/role-manager.iam.amazonaws.com/AWSServiceRoleForIAMRoleManager |
| 7 | AWSServiceRoleForNotificationsContactAssociations | iam | role | eu-north-1 | active | arn:aws:iam::004406188742:role/aws-service-role/contactassociations.notifications.amazonaws.com/AWSServiceRoleForNotificationsContactAssociations |
| 8 | AWSServiceRoleForOrganizations | iam | role | eu-north-1 | active | arn:aws:iam::004406188742:role/aws-service-role/organizations.amazonaws.com/AWSServiceRoleForOrganizations |
| 9 | AWSServiceRoleForSupport | iam | role | eu-north-1 | active | arn:aws:iam::004406188742:role/aws-service-role/support.amazonaws.com/AWSServiceRoleForSupport |
| 10 | AWSServiceRoleForTrustedAdvisor | iam | role | eu-north-1 | active | arn:aws:iam::004406188742:role/aws-service-role/trustedadvisor.amazonaws.com/AWSServiceRoleForTrustedAdvisor |

Created dates range from 2026-09-18 to 2026-11-16. Two are customer-managed (`/managed/` path: `AccountFullAccessRole`, `AWSManagedAccountManagementAccessRole`, `AWSManagedAdvancedFeaturesActivationRole`); the rest are AWS service-linked roles.

**Score summary (`/tmp/score2.json`):**
```json
{
  "summary": {
    "portable": 0,
    "rework": 0,
    "rewrite": 0,
    "unknown": 10,
    "total": 10
  }
}
```
Every resource: class `unknown`, source `unmapped`.

---

## 3. Infrastructure alternatives analysis

The account holds **zero compute, storage, or network resources** in eu-north-1 — only IAM identity roles. There is no running workload to migrate today, so a move off AWS is primarily an identity/control-plane question plus any future workload placement. Alternatives below assume standing up equivalent capacity on another provider or on-prem/local.

### Option A — Stay on AWS (eu-north-1)
- **Pros:** Zero migration effort; existing IAM roles and credentials already valid; Stockholm region (eu-north-1) offers full AWS service coverage; GDPR-compliant EU data residency.
- **Cons:** AWS is typically the most expensive of the options; vendor lock-in; egress charges.
- **Fit:** Baseline for cost comparison. Since there are no billable resources currently, monthly AWS spend is effectively the IAM/control-plane minimum (~$0–few USD/month for API calls, CloudTrail, Config if enabled).

### Option B — Hetzner (Nuremberg / Helsinki, EU)
Hetzner offers dedicated servers and cloud instances in Nuremberg (Germany) and Helsinki (Finland), both EU/GDPR.
- **Typical pricing (approx., monthly):**
  - Cloud CX11 (1 vCPU, 2 GB RAM, 20 GB NVMe): ~€4–5
  - Cloud CX22 (2 vCPU, 4 GB RAM, 40 GB NVMe): ~€12–14
  - Dedicated AX41 (4 vCPU, 32 GB RAM, 2×240 GB NVMe): ~€55–65
  - Dedicated AX51 (6 vCPU, 64 GB RAM, 2×960 GB NVMe): ~€95–110
  - Included traffic: ~20–200 TB/month depending on product; egress generally free/unmetered within reason.
- **Pros:** Significantly cheaper than AWS equivalent compute; EU data residency; bare-metal option for predictable performance; simple pricing, no egress surprise.
- **Cons:** Fewer managed services (no managed RDS/Lambda/IAM equivalent); self-managed OS, patching, backups; no IAM-style fine-grained roles (you'd manage users/SSH/SSO yourself); smaller global footprint (EU-only here).
- **Fit:** Strong for lift-and-shift of a VM-based workload; weak if the workload relies on AWS managed services.

### Option C — OVHcloud (Paris / Strasbourg / Frankfurt, EU)
OVHcloud offers public cloud instances, bare metal, and private cloud in Paris (FR), Strasbourg (FR), Frankfurt (DE), etc.
- **Typical pricing (approx., monthly):**
  - Public cloud b2-7 (2 vCPU, 7 GB RAM, 80 GB): ~€20–25
  - Public cloud b2-15 (4 vCPU, 15 GB RAM, 160 GB): ~€45–55
  - Bare metal Advance-1 (4 vCPU, 32 GB RAM, 2×240 GB): ~€55–70
  - Bare metal Rise (8+ vCPU, 64+ GB): ~€110–160
  - Included traffic: unmetered on many plans; public cloud has outbound allowances.
- **Pros:** EU data residency; broad product range (public cloud, bare metal, hosted private cloud); competitive with Hetzner; managed Kubernetes and databases available.
- **Cons:** Still self-managed for IAM-equivalent; pricing tiers can be complex; support tiers add cost.
- **Fit:** Good middle ground if you want some managed services (K8s, DB) but cheaper than AWS.

### Option D — On-prem / local server
A single server colocated or in-office (e.g., in Stockholm or nearby EU).
- **Typical hardware (one-time):** 2U rack server (2×16-core CPU, 128 GB RAM, 4×4 TB NVMe): ~€6,000–12,000 upfront.
- **Recurring:** colocation rack space + power + bandwidth in Stockholm: ~€150–400/month; or office broadband for a small server (no SLA, no redundancy).
- **Pros:** Full data sovereignty; no recurring cloud compute cost after amortization; no egress charges within LAN.
- **Cons:** High upfront capex; you manage hardware, power, cooling, networking, security, backups, and IAM from scratch; no elastic scale; single point of failure without redundant gear.
- **Fit:** Only justified if the workload is steady, predictable, and sovereignty requirements rule out third-party clouds. Overkill for an account that currently has zero running resources.

### Cost comparison reasoning (per month, illustrative)

Assuming a modest workload equivalent to ~2 vCPU / 4 GB RAM:

| Provider | Product | Approx. €/mo | Notes |
|---|---|---|---|
| AWS eu-north-1 | t3.small (2 vCPU, 2 GB) + EBS | ~€20–30 | Plus egress, managed-service premiums |
| Hetzner | CX22 cloud | ~€12–14 | ~40–50% cheaper than AWS; traffic included |
| OVHcloud | b2-15 public cloud | ~€45–55 | Larger instance; K8s/DB add-ons extra |
| Local server | Amortized hardware + colo | ~€200–400 | Upfront capex €6k–12k; breaks even after ~1–2 yrs |

Hetzner/OVHcloud are clearly cheaper than AWS for equivalent raw compute and include traffic, but neither replaces AWS IAM/service-linked roles — you'd rebuild identity and access management from scratch (users, SSH keys, SSO/OIDC, RBAC). Given the current inventory contains only IAM roles and no workloads, there is no running cost to save yet; the alternatives matter only when you launch a workload.

### Recommendation
Since no compute exists in eu-north-1 today, defer infrastructure migration until a workload is provisioned. When you do:
- For a simple VM/container workload → Hetzner (cheapest, EU, traffic included).
- If you need managed K8s/databases alongside → OVHcloud.
- If you must avoid any third-party cloud → local server, accepting capex and ops burden.
- Re-run `domovyk price` with a real FOCUS billing export once AWS usage generates billing data to get an actual cost comparison instead of estimates.

---

## 4. Files produced / referenced

- `/tmp/inventory.json` — real AWS inventory (eu-north-1, 10 IAM roles), source of findings above.
- `/tmp/score2.json` — score output from this session (10 resources, all unknown/unmapped).
- `/tmp/inventory2.json` — not produced (discover failed, expired credentials).
- `map/domovyk-map.yaml` — mapping used by score; lacks IAM-role entries (source of all `unknown` classes).
- `cmd/domovyk/main.go` — CLI: `discover`, `score`, `price` (`price` needs `--focus` FOCUS export).
