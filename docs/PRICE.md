# Price comparison semantics

Run from the repository root:

```
go run ./cmd/domovyk price -focus export.csv -prices map/prices.yaml -out price.json
```

Supply `-fx USD=0.92` explicitly for USD billing. Rates multiply billed amounts
into EUR. All output amounts and deltas are EUR. Delta is AWS minus target.

Input must contain one calendar month; partial months are not extrapolated.
BilledCost is summed per ServiceName, including credits. EffectiveCost is
validated but is not the cost basis. Only Usage rows add consumed quantities.
Every supplied price of the mapped category with a matching unit produces an
alternative: consumption times price_eur. Hour aliases are normalized. No GB/GiB,
monthly/hourly, resource-size, or capacity conversions are inferred. These are
consumption comparisons, not a verified equivalent-capacity migration quote.

The existing seed map supplies target eligibility. Supported seed IDs are
aws-ec2-instance, aws-ebs-volume, aws-s3-bucket, aws-rds-postgresql. Their target
must be nonempty. Common FOCUS service aliases are supported. RDS aliases use
the PostgreSQL mapping as a category candidate, not proof of engine equivalence.
Other prose targets have no supported price category and are skipped. Mixed
EC2 usage containing volume, snapshot or transfer resource types is not priced
as compute. Incompatible or multiple consumption units produce no_price.

Price tables contain sourced rates in `rows`. See map/prices.NOTES.md for
the separately supplied snapshot and its limitations. Each row requires provider (hetzner or
ovhcloud), sku, category (compute, block_storage, object_storage, managed_db or
egress), unit, price_eur, an HTTP(S) source_url, checked_on (YYYY-MM-DD), and
status (verified or unverified). Status is supplied by the table author; the
CLI does not verify web sources. Unverified rates print as UNVERIFIED.
A null or missing price_eur is allowed only for unverified entries and is
never treated as zero; that entry cannot produce a numerical estimate.

For no_price, estimates and deltas are null, checked_on is `not_available`, and
source_url is the file URL of the inspected price table. These fields identify
missing evidence without inventing a verification date or provider price.
JSON is written atomically with mode 0600.
