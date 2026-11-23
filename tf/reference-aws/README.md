# Reference AWS setup

Purpose: a small, known AWS estate to test discover, score and (later) price against. It is the input for v0.1.

## Resources

- EC2: one small instance in a single VPC, public subnet, with a security group.
- RDS: one small PostgreSQL instance in a private subnet, not publicly reachable.
- S3: one bucket with versioning on and public access blocked. Holds only dummy data.
- ALB: one load balancer, one target group pointing at the EC2 instance.
- VPC: one VPC, two subnets, one internet gateway, route tables.

Each resource gets tags: `domovyk-ref = "true"`, `owner`, `expires`.

## Files

- `main.tf` — provider, variables, common tags, terraform required version
- `variables.tf` — input variables (region, instance_type, db_instance_class, owner, expires, db_password)
- `vpc.tf` — VPC, subnets, internet gateway, route tables, security groups
- `ec2.tf` — EC2 instance (Ubuntu Jammy)
- `rds.tf` — RDS PostgreSQL subnet group + instance
- `s3.tf` — S3 bucket with versioning and public access blocked
- `alb.tf` — ALB, target group, listener, attachment

## Rules

- Deploy only into a dedicated sandbox AWS account. Not a production account.
- Creating any resource costs money. Spend is the owner's decision; no cost figure is given here (UNVERIFIED).
- No credentials in the repo. Use the owner's own AWS profile at apply time.
- No Lambda in the reference setup. The point is a non-Lambda baseline.
- Tear down with `tofu destroy` when the v0.1 tests are done.

## Open

- Whether the state backend is local or remote (not decided; local is simplest).
