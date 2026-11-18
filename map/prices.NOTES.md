# Price collection notes — 2026-11-18

## Method
Prices collected from official pages only (hetzner.com, ovhcloud.com). All prices EUR, excl. VAT.
checked_on: 2026-11-18.

## Hetzner Cloud — hetzner.com
- CPX/CX/CAX shared vCPU: read from /cloud/cost-optimized/ (CX, CAX) and /cloud/regular-performance/ (CPX).
- CAX/CX types listed as "not available" on cost-optimized page at time of read; prices still displayed.
- Traffic overage: pages state "€ 1.00 per hour, excl. VAT" per additional TB (EU-central). Unit wording on the page is ambiguous ("per hour" per TB); Hetzner's list price is commonly €1.00/TB. Row kept at 1.00 EUR/TB with status verified from page text.
- Load Balancer: LB11 = €7.49/month. Not in the allowed `category` enum (compute|block_storage|object_storage|managed_db|egress); recorded under `category: compute`. LB is its own product; see page /cloud/load-balancer/.
- Object Storage: base €6.49/month incl. 1 TB storage + 1 TB egress; excess storage €0.0087/TB-hour, excess egress €1.00/TB. Standard storage row = 6.49 EUR/TB-month (bundled base). Egress row = 1.00 EUR/TB (excess).
- Volumes (block storage): official block-storage page and docs do NOT print a per-GB-month price (point to hetzner.com/cloud/block-storage, which has no price table). Row `volumes` left `status: unverified`, `price_eur: null`. Third-party sources cite ~€0.044/GB/month but that is NOT an official page and was NOT used.

## OVHcloud Public Cloud — ovhcloud.com (German page ovhcloud.com/de/public-cloud/prices/)
- OVHcloud official pages display prices in USD on /en/ and in EUR (incl. VAT, 19% MwSt) on /de/.
- EUR prices shown are INCL. VAT. excl. VAT = incl / 1.19. Monthly excl. VAT = hourly excl. VAT x 730 h.
  - b3-8: incl 0.060928 €/h → excl 0.0512 €/h → 37.38 €/month
  - b3-16: 74.68; b3-32: 149.36
  - c3-4: 33.36; c3-8: 66.65; c3-16: 133.23
- Block Storage Classic: 0.00007021 €/GB/h incl VAT → excl ≈ 0.0000590 → ×730 = 0.0431 €/GB/month.
- Object Storage Standard (Multi Zone, first 50 TiB): 0.00002281 €/GiB/h incl VAT → excl ×730/1024 (GiB→GB) ≈ 0.0140 €/GB-month.
- Object Storage egress: listed "Kostenlos" (free) → 0.00 EUR/TB.
- Managed PostgreSQL smallest plan: Essential DB1-4 = 0.088774 €/h/node incl VAT → excl ×730 = 54.46 €/month. (Gen 3 Discovery B3-4 is marginally cheaper at 0.10179 incl VAT → ~62.6 €/month; Essential DB1-4 used as smallest listed essential plan.)
- Load Balancer / traffic / gateway rows exist but are out of the requested coverage scope except where noted.

## Pages that could not be read / no price table found
- https://www.hetzner.com/cloud/block-storage/ — no per-GB price printed (volumes row unverified).
- https://www.ovhcloud.com/en/public-cloud/instances/prices/ — HTTP 404.
- https://www.ovhcloud.com/en/public-cloud/pricing/ — HTTP 404.
- https://docs.hetzner.com/cloud/pricing/ — HTTP 404.

## Not covered (out of scope or not listed)
- Hetzner CAX 8 vCPU in 2/4/8 vCPU request: CAX 8 vCPU = cax31 (21.49 €) included.
- No separate Hetzner "traffic overage" product page; overage price taken from server detail pages.
