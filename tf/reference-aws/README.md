# Reference AWS setup (PLAN ONLY)

No `.tf` files exist yet. This is the plan for the v0.0 reference setup in OpenTofu.

Purpose: a small, known AWS estate to test discover, score and (later) price against. It is the input for v0.1.

## Resources (planned)

- EC2: one small instance in a single VPC, public subnet, with a security group.
- RDS: one small PostgreSQL instance, private subnet, not publicly reachable.
- S3: one bucket with versioning on and public access blocked. Holds only dummy data.
- ALB: one load balancer, one target group pointing at the EC2 instance.
- VPC: one VPC, two subnets, one internet gateway, route tables.

Each resource gets tags: `domovyk-ref = "true"`, `owner`, `expires`.

## Rules

- Deploy only into a dedicated sandbox AWS account. Not a production account.
- Creating any resource costs money. Spend is the owner's decision; no cost figure is given here (UNVERIFIED).
- No credentials in the repo. Use the owner's own AWS profile at apply time.
- No Lambda in the reference setup. The point is a non-Lambda baseline.
- Tear down with `tofu destroy` when the v0.1 tests are done.

## Open

- Exact instance sizes and region (EU region preferred, not decided).
- Whether the state backend is local or remote (not decided; local is simplest).
