# DRAFT: career-agent (short)

Source folder: `~/Documents/Acore/Development/career-agent` (read only, not modified).

## What it is

Private job-hunt pipeline: discover roles, score fit, tailor CV and cover letter, approve, apply through a browser agent, track outcomes (`README.md`, first line). Python package `career-agent` 0.1.0 with a `career` CLI (`pyproject.toml`).

Also includes a stdlib-only freelance intake script, `scripts/freelance_watch.py`, with rule-based candidate qualification and a review queue (`README.md`, "Bounded freelance watch").

## Measured results

[TBD] No measured results found in the project's own files. Status in `README.md` says "Phase 1 in progress".

## How to run

```bash
pip install -e ".[dev]"
career init
python3 scripts/freelance_watch.py --fetch-hn --stats
python3 -m unittest discover -s tests -p test_freelance_watch.py
```

Profile data must be created from `profile/master_profile.example.json`. Keep it private.

## Honest limitations

- Private. Personal profile data and application state live in the repo. Do not publish the repo.
- Phase 1 only. Apply step requires the owner's signed-in browser and manual submit.
- Keyword rules in the freelance script are not fit scores; they need human review (`README.md`).
- Not a public product.
