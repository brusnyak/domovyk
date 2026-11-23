# Domovyk competitive landscape and product decision

**Research date:** 2026-11-23  
**Status:** Desk research. Vendor capabilities are based on public documentation and marketing. No competitor migration was run or independently benchmarked.

## Executive decision

- **No-go:** build Domovyk as another generic AWS exit assessment, service mapper, cost calculator, or Terraform generator.
- **Conditional go:** investigate a customer-operated, open-source migration execution engine that composes existing tools, keeps credentials/data in the customer's environment, and provides durable approvals, checkpoints, evidence, validation and tested recovery.
- **Do not start broad implementation yet.** The competitive gap is not proven. MoveInfra is a direct competitor and claims end-to-end AWS-to-EU migration. Validate a narrow buyer/use case and test a real workload first.

The current README describes Domovyk as a pre-alpha AWS-to-EU planner; discovery/score are described as implemented but PROGRESS says build/tests have not run and generate/cutover do not exist. That is too early to claim a product advantage.

## 1. Direct competitors

### MoveInfra — closest overlap

- Publicly markets migrations from AWS, Azure, Google Cloud and DigitalOcean to UpCloud, Scaleway, Hetzner and OVHcloud.
- Describes a six-stage workflow: analyze, copy infrastructure, migrate data, update deployment pipeline, switch, retain rollback.
- Offers free bill analysis, a €1,900 exit plan (credited if migrating), and migration service starting at €4,900; pricing is published on its site and subject to its offer/conditions.
- Publishes open-source tools/data: bill analyzer, Kubernetes cost analyzer, cloud lock-in scanner, switching-fee calculator, EU Cloud Price Index, EU Cloud Map and an EU-cloud MCP server.
- Its AWS→Scaleway page claims EC2/RDS/S3/ALB/VPC/EKS/ElastiCache automation, with SQS/SNS/Lambda/DNS assisted and DynamoDB/Redshift not yet offered. Its route matrix states fewer automatic resources for Hetzner than for UpCloud/Scaleway/OVHcloud.

**Assessment:** direct competitor to the current plan and proposed broad execution workflow. Their migration capabilities are vendor claims; we have not tested their code or a live migration. Public site describes early access/waitlist, so operational maturity and customer results remain unverified.

Sources:
- https://moveinfra.eu/
- https://moveinfra.eu/tools/
- https://moveinfra.eu/migrate/
- https://moveinfra.eu/migrate/aws-to-scaleway/
- https://moveinfra.eu/pricing/
- https://github.com/moveinfra/eu-cloud-prices
- https://github.com/moveinfra/cloud-exit-cost-analyzer
- https://github.com/moveinfra/cloud-switching-kit
- https://github.com/moveinfra/cloud-lock-in-scanner

### EscapeCloud / cloudexit — assessment competitor

- Open-source, local-first CLI; AWS and Azure supported, GCP/OCI on roadmap.
- Inventory, cost inventory, rule-based exit risk, alternative technologies, HTML/PDF/JSON reports, Terraform/OpenTofu state input and CI action.
- No full migration execution: public repo positions it as an exit-readiness assessment engine, with hosted tiers for reports/scoring/governance.
- AGPL-3.0; commercial licence available through vendor.

**Assessment:** Domovyk discover/score/price overlap directly. Do not duplicate these unless Domovyk has a clear measured advantage or builds an integration.

Source: https://github.com/escapecloud/cloudexit

## 2. Migration platforms and cloud-native alternatives

| Alternative | Strength | Structural limitation / implication |
|---|---|---|
| AWS Transform / AWS Transform MGN | Agentic discovery, wave planning, server rehost/replatform, cutover and modernization | Primarily migration into AWS; not an AWS exit product |
| Azure Migrate | Discovery/assessment and server replication/cutover from multiple sources | Destination is Azure |
| Google Migration Center / Migrate to VMs | VM migration from AWS/Azure/VMware and test/cutover | Destination is Google Cloud |
| RiverMeadow | Commercial workload mobility and modernization across public/private cloud and virtualization | Enterprise commercial platform; broad overlap with orchestration |
| RackWare RMM | Provider/hypervisor-neutral replication, assessment and recovery | Commercial, quote-led platform |
| OpenText Migrate / HPE Zerto | Cross-platform replication and low-downtime migration/DR | Commercial enterprise tools; often replicate workloads rather than translate every managed service |

Sources:
- https://docs.aws.amazon.com/transform/latest/userguide/
- https://learn.microsoft.com/en-us/azure/migrate/execute-server-migrations
- https://docs.cloud.google.com/migrate/virtual-machines/docs/5.0/discover/architecture
- https://www.rivermeadow.com/platform/migrate
- https://www.rivermeadow.com/clouds-platforms
- https://www.rackwareinc.com/
- https://www.ciopages.com/buyer-guides/cloud-migration-platform

## 3. Existing tools Domovyk should compose, not replace

| Migration layer | Existing tools / role |
|---|---|
| Infrastructure reconstruction | OpenTofu/Terraform plus destination provider plugins |
| AWS resource export | AWS APIs/CLI; Former2 for generated IaC; Terraform state/import workflows |
| Object storage | rclone; S3-compatible endpoints, checksum verification |
| PostgreSQL | native logical replication / pgcopydb or pg_dump/restore, depending on downtime and version |
| Kubernetes | Velero or Kasten for backup/restore; GitOps (Argo CD/Flux) for redeploy |
| Database CDC / heterogeneous data | AWS DMS, destination-native migration services, Debezium where appropriate |
| Inventory / cost / portability | cloudexit, MoveInfra open tools and provider billing exports |

Sources:
- https://rclone.org/
- https://pgcopydb.readthedocs.io/
- https://velero.io/
- https://opentofu.org/
- https://docs.aws.amazon.com/decision-guides/latest/decision-guides/migration-on-aws-how-to-choose.html
- https://moveinfra.eu/tools/

## 4. AWS → EU destination capability map

Planning-level map, not a tested compatibility guarantee. “Equivalent” means a candidate service exists; it does not mean semantics, APIs, IAM, performance, or failover are identical.

| AWS source | UpCloud | Scaleway | OVHcloud | Hetzner | Migration class |
|---|---|---|---|---|---|
| EC2 | Cloud servers | Instances | Public Cloud instances | Cloud servers | Rehost; low-medium |
| EBS | Block storage | Block volumes | Block storage | Volumes | Image/snapshot conversion; medium |
| S3 | Object storage | Object Storage | Object Storage | Object Storage | S3 API + endpoint/ACL/versioning checks; low-medium |
| RDS PostgreSQL/MySQL | Managed DB | Managed DB | Managed DB | Typically self-managed | Logical replication/restore; medium |
| EKS | Managed Kubernetes | Kapsule | Managed Kubernetes | Self-managed K8s | Recreate cluster/add-ons; medium-high |
| ALB/NLB | Load balancer | Load Balancer | Load Balancer | Load Balancer | Rebuild listeners, health checks, TLS; medium |
| ElastiCache Redis | Managed Redis | Managed Redis | Service/alternative varies | Self-managed/third-party | Data/cluster and client validation; medium |
| Lambda | Functions or containers | Functions or containers | Functions/containers | Containers | Runtime + trigger rewrite; high |
| SQS/SNS/EventBridge | Messaging alternatives | Queues/topics where supported | Queue/event alternatives | Self-managed broker | Delivery/retry/ordering semantics; high |
| DynamoDB | No direct drop-in | No direct drop-in | No direct drop-in | No direct drop-in | Data model/application redesign; high |
| IAM/KMS/Secrets | Provider IAM/secrets | Provider IAM/secrets | Provider IAM/secrets | Provider IAM/secrets | Rebuild trust, policy, key and secret model; high |
| CloudWatch | OTel + Prometheus/Grafana/Loki | Same | Same | Same | Rebuild telemetry/alerts/retention; medium |

Provider capability sources:
- https://www.upcloud.com/products/
- https://www.scaleway.com/en/
- https://www.ovhcloud.com/en/public-cloud/
- https://www.hetzner.com/cloud/
- https://developers.eu.ovhcloud.com/terraform/guides-start/
- https://github.com/hetznercloud/terraform-provider-hcloud
- https://moveinfra.eu/migrate/aws-to-scaleway/

The first useful migration target is not “all AWS services.” It is a defined reference architecture with a tested matrix: VM/container compute, PostgreSQL, S3-compatible storage, networking, DNS and observability. Serverless, event-driven, NoSQL and data warehouse workloads must be explicit re-architecture projects or exclusions.

## 5. Market and regulatory implications

- The EU Data Act removes switching charges from **12 January 2027** (Article 29). This reduces one friction but does not remove migration labour, parallel-run costs, application rewrites, operational risk or target validation.
- AWS already offers a free data-transfer-out programme for customers moving off AWS, with eligibility/process conditions. “Save egress fees” is not a durable product moat.
- DORA exit-strategy requirements are relevant to regulated financial entities, but Domovyk must not claim compliance without a validated evidence/workflow mapping and legal review.

Sources:
- https://eur-lex.europa.eu/eli/reg/2023/2854/oj/eng
- https://digital-strategy.ec.europa.eu/en/factpages/data-act-explained
- https://aws.amazon.com/blogs/aws/free-data-transfer-out-to-internet-when-moving-out-of-aws/
- https://www.bankingsupervision.europa.eu/ecb/pub/pdf/ssm.supervisory_guides202507.en.pdf

## 6. Product recommendation

### Stop building these as standalone differentiators
- Generic inventory and exit-readiness scoring
- Generic cost comparison / price index
- Generic AWS service-to-EU mapping
- Generic Terraform/OpenTofu generation
- A report-only product

These already exist in cloudexit, MoveInfra, native providers and consultancy offerings.

### Investigate this narrow wedge

**Customer-run, open-source, code-first migration workflow for repeatable AWS → EU moves.**

Possible differentiation (not proven):
1. Runs in the customer's own machine/CI runner; no hosted control plane or customer inventory upload.
2. Uses existing exporters, OpenTofu providers, rclone, database replication and Kubernetes tools rather than reimplementing them.
3. Durable state machine, idempotent steps, approvals, checkpoint/resume and immutable audit/evidence.
4. Explicit validation gates and tested rollback/reverse-sync, not merely “DNS can be changed back.”
5. Provider adapter and recipe contract; first route selected from a real pilot, not a broad provider matrix.

This still competes with MoveInfra's claimed automation. The possible distinction is self-service, self-hosted, transparent execution software rather than a managed migration service—but customer demand must be demonstrated.

## 7. Go/no-go gate before engineering

Proceed to a thin vertical slice only after:
- 5–10 cloud consultants/MSPs or engineering leads review the proposed workflow;
- at least one partner supplies a sanitized Terraform state/billing export or non-critical test workload;
- they identify a repeated pain that MoveInfra/cloudexit/native tools do not solve adequately;
- the selected route and workload are bounded, and a safe test environment is available.

**Pilot success criteria:** reproducible plan; no source-side writes by default; target IaC plan review; data checksum/DB validation; documented downtime; recoverable workflow after interruption; approved cutover; demonstrated rollback/recovery limits; complete audit trail.

## Sources and caveats

Sources are primarily official vendor documentation, project repositories and vendor product pages. Vendor claims (including MoveInfra automation, pricing, time-to-migrate and rollback claims) are not independent benchmarks. No product installation, code audit, pricing quote, customer interview or live migration was performed in this research.
