# DRAFT: Live Speech Translation (bp-main)

Source folder: `~/Developer/bp-main` (read only, not modified).

## What it is

Local speech translation and live captions for calls and presentations. Microphone to VAD to speech recognition to translation to voice, output to speakers or captions. Processing runs on the user's own machine after model installation (`README.md`, first section).

## Measured results (only from the project's own files)

- Live route EN to SK on M1 Pro, median delay from end of speech to first translated audio: Whisper base 0.82 s; turbo 1.01 s (`README.md`, "Saved measurements"; evidence `documentation/evidence/2026-11-07/`).
- German speech recognition, mean word error on FLEURS, 60 clips per language, Kaggle T4: turbo 4.8%, Parakeet v3 6.3% (`README.md`; table `documentation/evidence/2026-11-07/lang_matrix.json`).
- CS to EN full route on Kaggle T4, default profile: p50 3.48 s, p95 4.13 s (`README.md`; `route_receipt.json`).
- Streaming simulation: no setting reached under 1 s latency with under 10% word error (`README.md`; `documentation/evidence/2026-11-07/stream_sim.json`).

Caveat in the source: Parakeet ran on CPU in that run, so its speed is not a GPU figure.

## How to run

```sh
git clone https://github.com/brusnyak/bp.git
cd bp
python3.11 scripts/setup.py --profile mac --langs cs,sk,de
make run PRESET=mac
```

Open http://localhost:8000 in Chrome or Edge. Setup downloads models; allow disk space and internet for installation.

## Honest limitations

- Route quality varies: English to Slovak is experimental; Slovak to English is experimental.
- Room mode is experimental.
- Licences: Piper runtime is GPL-3.0-or-later; XTTS v2 and X-Voice weights are non-commercial (`README.md`, "Privacy and licences"). Not cleared for commercial use.
- Measurements are on specific hardware and clocks; they are not general claims.

## Portfolio framing

[TBD] one-line framing for the CV, after owner review.
