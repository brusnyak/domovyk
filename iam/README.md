# iam/ (UNREVIEWED)

`domovyk-readonly-policy.json` is a candidate least-privilege read-only policy. **UNREVIEWED**: not checked against AWS docs, not run through the IAM policy simulator, not reviewed by anyone else.

Deliberate exclusions: no `s3:GetObject` or `s3:ListBucket` (no object reads or key listings); no `ec2:GetPasswordData`, `ec2:GetConsoleOutput`, `ec2:DescribeInstanceAttribute` (can expose user data); no `lambda:GetFunction` (returns a code download URL); no `iam:Get*`.

Also unknown: whether `eks:DescribeCluster` returns anything sensitive. Check before v0.1.
