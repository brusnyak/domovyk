# RESEARCH: domovyk

Labels: **VERIFIED** (checked by the owner, with source), **UNVERIFIED** (not checked; do not cite as fact).

## Regulatory context

- **VERIFIED (owner):** The EU Data Act ends one-off cloud switching charges from 12 January 2027.
- **VERIFIED (owner):** Egress for parallel multi-cloud use stays allowed at cost. Routine egress is not covered by the switching rule.

## Market check

- **VERIFIED (owner, web search):** Search found only consultancies: Gart Solutions, WZ-IT, MoveInfra. No maintained open-source AWS to Hetzner tool with OpenTofu output was found. The search was not exhaustive.

## Tools and projects (all UNVERIFIED unless stated)

| Name | Relevance | Status claimed | Label |
|---|---|---|---|
| Terraformer | Export AWS resources to Terraform | Archived 16 Mar 2026 | UNVERIFIED |
| Former2 | Generate templates from live AWS | About 49% resource coverage | UNVERIFIED |
| Steampipe | SQL over cloud APIs | Plugin licence terms | UNVERIFIED (licences) |
| MinIO (OSS) | S3-compatible target option | Repository archived Feb 2026 | UNVERIFIED |
| Hetzner cloud-controller | Kubernetes integration on Hetzner | v1.30.0 reported to crash | UNVERIFIED |
| Hetzner prices | Target cost | Price rises in 2026 | UNVERIFIED |

No price number appears in this repo. Any price used later must come with a `checked_on` date.

## Mapping seed data

The 30 rows in `map/domovyk-map.yaml` each carry `checked_on: "2026-11-23"` and `status: reviewed` per `map/SCHEMA.md`. Sources cited in the file header: MoveInfra's EU equivalents of AWS services (published 2026-11-06) and Hetzner's live pricing page (read 2026-11-23). Row-level claims should still be re-checked against primary sources before production use.
## Open questions

1. Which mapping classes hold up on real accounts? Test on the owner's `wocus` and `recareo-web` accounts before trusting any score.
2. How are OVHcloud and Hetzner price lists obtained: bundled snapshot or live fetch? Licence of the price data is unknown.
3. Does FOCUS 1.2 export carry enough detail for per-service pricing? Check against a real export.
4. Is `ec2:DescribeInstanceAttribute` needed at all? It can return user data. Current policy excludes it.
5. Can a Lambda-heavy account be detected reliably from inventory alone?

## Must re-check before v0.1

- **GitHub search** for existing AWS-to-EU exit or migration tools, with OpenTofu output. Search terms: "aws exit", "cloud exit planner", "terraform aws hetzner migration", "aws to ovh terraform". Record each hit with URL and date.
- **Hacker News search** (Algolia HN search) for the same terms and for "cloud exit" or "EU cloud migration" over the last 12 months.
- Re-check each UNVERIFIED row above against its primary source (repo archive status, price pages, licence files).
- Re-check EU Data Act wording against the Official Journal text, not a summary.
