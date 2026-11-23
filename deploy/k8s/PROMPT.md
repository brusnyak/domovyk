# Paste-ready prompt: Colima + kind + Domovyk, then the rest of the queue

Paste everything under the line into a fresh Claude Code session started in `~/Developer/domovyk`.
Written 2026-11-17. State below was verified that day; re-check before trusting it.

---

You are finishing a bounded job for Yegor. Work in order, show real command output, label anything not run as `unverified:`. Caveman style, result first, max 25 lines per report. Never print secrets. No spend, no cloud resources, no push until step 7. Stop after two failures on one step and report the error instead of working around it.

## Known state (checked 2026-11-17)
- Mac: M1 Pro 16 GB. Docker Desktop is NOT installed; Colima is. `kind v0.33.0` and `kubectl 1.37.2` are installed via brew (logged in `~/memory-vault/pages/install-log.md`).
- Colima failed twice: the first VM hung on `docker build` and `colima stop -f`; I ran `colima delete -f` and `colima start --cpu 2 --memory 4` (worked, hello-world pulled), then the VM died during `docker build` while load average hit 8.0 (WindowServer, Orca renderer, ego lite also busy). `colima status` then said "error retrieving current runtime: empty value".
- Files already written, never built or applied (unverified): `~/Developer/domovyk/Dockerfile` (multi-stage golang:1.24 to distroless nonroot, copies `map/domovyk-map.yaml`), `~/Developer/domovyk/deploy/k8s/job.yaml` (ConfigMap with a 5-resource sample inventory, a Job running `domovyk score`, nonroot uid 65532, read-only rootfs, caps dropped, seccomp RuntimeDefault, no token, limits, `imagePullPolicy: Never`, image `domovyk:dev`, no secrets or AWS creds).
- `deploy/k8s/README.md` does not exist yet. Do not write outputs you did not produce.
- The repo has one commit (74d043a), no GitHub remote. `LICENSE` is a placeholder (needs the full Apache-2.0 text). `go.mod` module is `github.com/brusnyak/domovyk`.

## Step 0: calm the machine
1. `uptime`. If load average > 4, quit Orca and close ego-browser tabs you do not need, wait, re-check. Do not start Colima while load > 4.
2. `colima stop -f` with `timeout 60`. If it hangs: `pkill -f "limactl hostagent"` then `colima delete -f` (the VM holds nothing of value; images are rebuilt).

## Step 1: start Colima small and prove it works
```
timeout 170 colima start --cpu 2 --memory 4 --disk 20
export DOCKER_HOST=unix://$HOME/.colima/default/docker.sock
timeout 60 docker run --rm hello-world
```
Pass = "Hello from Docker!". If this fails twice, stop and report.

## Step 2: build the image in the background (it is slow, 2 CPUs)
Run with `run_in_background` and a long timeout, never in the foreground:
```
cd ~/Developer/domovyk && docker build -t domovyk:dev .
docker images domovyk:dev
```
Check `docker run --rm domovyk:dev --help` (or the binary's usage) works under the non-root user. If the Dockerfile is wrong, fix only the Dockerfile.

## Step 3: kind cluster
```
kind create cluster --name domovyk
kind load docker-image domovyk:dev --name domovyk
kubectl config current-context
```
Context must be `kind-domovyk`. Never apply to any other context.

## Step 4: deploy and capture real output
```
kubectl apply -f deploy/k8s/job.yaml
kubectl wait --for=condition=complete job --all --timeout=120s
kubectl get pods -o wide
kubectl logs job/<job-name>
kubectl get job -o yaml | grep -E "runAsNonRoot|readOnlyRootFilesystem"
```
Pass = job Complete and the log shows the scored sample inventory (portable / rework / rewrite counts, the unmapped glue resource as unknown).

## Step 5: README with proof
Write `deploy/k8s/README.md`: exact commands, the pasted real outputs, and a plain "What this does not prove" section (single-node kind, no ingress, no autoscaling, no production, no cloud cluster, sample data only).

## Step 6: tear down
```
kind delete cluster --name domovyk
colima stop
```
Report: what passed, what failed, the load average during the build.

## Step 7: publish Domovyk (Yegor already approved this push, conditional on step 4 passing)
1. Replace the `LICENSE` placeholder with the full Apache-2.0 text. Check `go build ./cmd/domovyk` and `go vet ./...` and `go test ./...`.
2. README already says "pre-alpha, nothing built": keep that honesty; add one line linking `deploy/k8s/README.md` as the container proof.
3. `git add -A && git commit -m "Add Dockerfile, kind deploy proof, Apache-2.0 license"`.
4. `gh repo create brusnyak/domovyk --public --source . --remote origin --push`.
5. `gh repo view brusnyak/domovyk --json visibility,url` to confirm it exists. The CV links this URL; a 404 is a defect.

## Step 8: CV line (only if step 4 passed)
Edit the single master `~/cv/cv-print.html` (backup first: `cp cv-print.html cv-print.html.bak-$(date +%Y%m%d)`), add to Skills Tools: "Kubernetes (kind, local)". Do not claim cloud clusters, years, or production use. Then `cp cv-print.html index.html`. Do not create any other CV file.

---

## The rest of the queue (do after the above, one at a time, gate each with pytest)
Repo: `~/Documents/Acore/Development/career-agent`. Baseline: `.venv/bin/python -m pytest -q` = 113 passed.

A. **Scoring gate (US-only roles scored 87-88).** A worker patch exists at `~/.agent-router/data/runs/92f43a22/patch.diff`, built from a stale snapshot: it adds a second `pre_screen_reject` and collided with the existing one (47 failures when applied). Do NOT `patch` it blindly. Merge by hand: keep the existing `pre_screen_reject(*, title, description, source)` signature, add the worker's geography/years logic as a helper (e.g. `geo_years_reject`) and call it from the existing function; keep tests green. Rules: reject US-only / "authorized to work in the US" / San Francisco / New York without remote-Europe wording; bounded bonus for Slovakia, Bratislava, Vienna, EU-remote; required years > 3 reduces or rejects.

B. **Slovak discovery** (13k listings, only 7 in Bratislava/Slovakia/Vienna; Profesia stored the nav link "Všetky ponuky" as a job; Kariera returns 0). Patch: `~/.agent-router/data/runs/bc0f55f0/patch.diff` (adds `discover/slovak.py`, fixtures, `tests/test_slovak_discovery.py`). It applies cleanly but 6 tests fail in the integrated repo (the "blocked/unparseable" cases). Fix the tests or the code, then re-run.

C. **Copy-engine test for auto-apply.** Build `scripts/copy_eval.py`: take 20 listings with fit_score >= 75 from `data/career_agent.db`, generate each cover letter with `pipeline._write_cover_letter`, then gate each with (1) `guard.gate_text` (fact gate vs master CV `~/cv/cv-print.html`), (2) `python3 ~/memory-vault/scripts/copy_lint.py --channel application` (window 80-220 words; the prompt currently says "under 250", align them), (3) `~/memory-vault/scripts/slopcheck.py`, (4) structure: one opening line naming company and role, three bullets, one closing line. Write results to `data/copy_eval/<date>.json` and a pass-rate line. PASS = 20/20 deterministic gates, then run the `recipient-panel` agent on 5 of them. Drafts only; nothing is sent. Log the result in the vault (one line in `pages/jobs-system/`).

D. **Security gate before any unattended run.** `python3 ~/.claude/scripts/config_selftest.py` was 21/25. Failing lines: CLAUDE.md size > 7500 bytes, "mentor+roast EVERY reply" line, "smart-lazy pre-flight" line, "PreToolUse Agent hook present". CLAUDE.md was edited by the owner on 2026-11-17 (rules moved out), so re-run first; fix only what still fails and only with the owner's yes (it is global config). Then `python3 ~/.claude/scripts/security_check.py` must be 9/9. Auto-apply stays off until both pass AND step C passes AND the owner records the go.

E. **Oracle private access (Tailscale).** `ssh oracle-vm` works (public IP, key-only, `PasswordAuthentication no`). The Mac's tailscaled runs userspace: socket `~/.tailscale-run/tailscaled.sock`, so plain ssh to `100.96.183.88` times out. Fix: (1) fetch the Oracle host key fingerprint over the working alias (`ssh oracle-vm 'ssh-keygen -lf /etc/ssh/ssh_host_ed25519_key.pub'`), (2) `ssh-keyscan` through `tailscale --socket=$HOME/.tailscale-run/tailscaled.sock nc 100.96.183.88 22`, compare fingerprints, add to known_hosts only if identical, (3) add an `oracle-ts` Host with `ProxyCommand tailscale --socket=... nc %h %p` (an auth-config edit: get the owner's yes), (4) test `ssh oracle-ts hostname`. Only after that works twice: on Oracle add `sudo ufw allow in on tailscale0 to any port 22`, test again, then `sudo ufw delete allow 22/tcp`, test `oracle-ts` again, and check the Oracle cloud network security list for port 22. Never close public 22 before `oracle-ts` is proven (the `oracle-vm` alias uses the public IP and the router's `@host:oracle-vm` depends on it: repoint it or it breaks).

F. **Oracle research lane.** `~/career-ops-live` exists on Oracle; install/run `career-agent` there (Python venv, `career run --no-submit` on a cron, Slovak/EU boards, results to the DB). Oracle has 11.9 GB RAM, ~9 GB free, 150 GB disk and already runs trading services: do not start anything heavy there. Submits stay on the Mac via ego-browser; LinkedIn and Indeed stay packet-only.

Report once at the end with `status:` and one `suggest:`.
