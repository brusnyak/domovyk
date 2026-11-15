# DRAFT: agent-router

Source folder: `~/Developer/agent-router` (read only, not modified).

## What it is

Measures whether the skills and agents installed for Claude Code are used, and tests whether surfacing the right one at prompt time changes that. Shadow hook logs what would have been surfaced; nothing is injected into sessions yet (`README.md`, "Why" and "How it works"). Also a decision layer that picks worker and command per task (`README.md`, "Decision layer").

## Measured results (only from the project's own files)

- 953 real prompts: 3% led to a Skill-tool call (`README.md`, "Why"; `docs/SKILL-ACTIVATION.md`).
- Silver labels: 517 of 1,011 prompts, 2 of 3 LLM judges agree (`README.md`, "Measured so far"; `eval/`). Labels are not human-verified.
- Retrieval top-1 / top-3, dense NIM embeddings: 25.0% / 46.0% (`HANDOFF.md`).
- Always-guess floor, top-1: 14.1% (`HANDOFF.md`).
- Compiled vs raw prompt, top-3: +1.2 pp against a 5 pp bar, failed (`README.md`, "Status and roadmap"; `eval/compile_test.py`).
- A/A replay, treat vs holdout follow rate: 1.7% vs 1.6%, p=0.89 (`README.md`).
- Hook cold start: 30 ms p50, 33 ms p95 (`README.md`, "Safety and invariants").
- Test count: the README badge says 192 passing (`README.md`, line 5). Not re-run in this check.

## How to run

```sh
python3.13 -m venv .venv
.venv/bin/pip install -r requirements.txt
.venv/bin/python -m backbone ask --dry "fix the failing test in the router"
.venv/bin/python -m backbone router-report
```

NIM keys are read from the environment only. Shadow hook install is `python3 scripts/shadow_hook.py install` and writes a settings backup first (`README.md`, "Quickstart").

## Honest limitations

- Private WIP. Routing claim is not made: dense top-1 (17.7%) is below the always-"research" guess (25.8%) on the frozen test split (`README.md`, line 40).
- Prompt compiler failed its bar.
- Labels are LLM-judged, not human-verified. Owner spot-check not done (`README.md`, "Status and roadmap").
- Live injection is blocked on the above.

## Portfolio framing

[TBD] one-line framing after owner review. Use the negative result as evidence of measurement discipline only if the owner agrees.
