# ai-sdlc — canonical process repository

This repository is the canonical source of truth for the shared **agentic SDLC process** used by multiple product repositories.

`ai-sdlc` describes the process only; project execution outputs are stored in each product repository under `ai-sdlc-artefacts/`.

---

## Repository layout

| Path | Role |
|------|------|
| **[specification/pipeline.spec.md](specification/pipeline.spec.md)** | **Normative process:** stage order, inputs/outputs, stage→skill mapping, **HOTL by default** (HITL at decision points), delegated execution, artefact naming, traceability, **agent execution expectations** (single process, AC validation, etc.). |
| **[specification/skills/](specification/skills/)** | **Per-stage agent instructions** (`00-` bootstrap, `01-` … `11-` plus optional cross-cutting skills). Each skill defines workflow and artefact structure for its stage. |
| **[specification/skills/README.md](specification/skills/README.md)** | Index and **common behaviour** across skills. |
| **[templates/consumer/](templates/consumer/)** | **Greenfield templates** copied to new product repositories during bootstrap. |
| **[reference/architecture-patterns/](reference/architecture-patterns/)** | **Advisory** architecture-pattern index for stage 6 (thin cards; upstream URLs are the source of truth). Versioned with the process via `ai-sdlc.version`; not part of the normative pipeline. |
| **[proposals.md](proposals.md)** | Open process proposals plus an index of adopted ones; a staging area, not the normative process. |
| **[tools/validate/](tools/validate/)** | AC↔test coverage checker (`./tools/validate/validate` from repo root); see [VALIDATION.md](tools/validate/VALIDATION.md) and [README](tools/validate/README.md). |
| **[tools/scripts/](tools/scripts/)** | Maintainer scripts (e.g. `bootstrap-smoke.sh` for greenfield template regression); see [tools/scripts/README.md](tools/scripts/README.md). |

---

## Starting a new project

Use this flow when the **product repository** is greenfield and you adopt ai-sdlc via a nested process clone.

### 1. Create the product folder and open it in the IDE

The workspace root must be the **product repository root** (e.g. `my-app/`), not `my-app/ai-sdlc/`.

### 2. Clone the process

```bash
cd my-app
git clone https://github.com/asubbot/ai-sdlc.git ai-sdlc
```

Alternatively: initialize the product git repo first, then clone `ai-sdlc/`.

### 3. Ask the agent to bootstrap

Example prompt:

> New consumer project. Run project bootstrap, then stage 1 scope analysis — ask me what we are building.

The agent runs [00-project-bootstrap.skill.md](specification/skills/00-project-bootstrap.skill.md) and materializes files from [templates/consumer/](templates/consumer/).

Maintainers: after changing consumer templates, run `./tools/scripts/bootstrap-smoke.sh` ([tools/scripts/README.md](tools/scripts/README.md)).

### 4. Continue the pipeline (normative order)

1. **Bootstrap** — skeleton (`AGENTS.md`, `Makefile`, `ai-sdlc.version`, `ai-sdlc-artefacts/`, CI).
2. **Stages 1–2** — refine [scope.md](specification/skills/01-scope-analysis.skill.md) and [strategy.md](specification/skills/02-strategy-analysis.skill.md).
3. **EP-000** — adoption epic (stages 3→11).
4. **EP-001+** — product epics.

### Target directory layout (summary)

**After bootstrap (minimum):**

```text
my-app/
├── AGENTS.md
├── README.md
├── Makefile
├── ai-sdlc.version
├── .gitignore              # includes ai-sdlc/
├── bin/validate
├── .github/workflows/ci.yml
├── .golangci.yml
├── scripts/check-module-boundaries.sh
├── ai-sdlc/                 # local clone (gitignored; not in product git)
└── ai-sdlc-artefacts/
    ├── scope.md            # stub, refine in stage 1
    ├── strategy.md         # stub, refine in stage 2
    └── epics/EP-000/       # adoption epic stubs
```

`ai-sdlc/` is listed in `.gitignore` and is not committed to the product repository. CI checks out the process at the pin in `ai-sdlc.version`. See [pipeline.spec.md](specification/pipeline.spec.md) (Consumer onboarding) and [templates/consumer/README.md](templates/consumer/README.md).

### Fresh clone of the product repo only

If `ai-sdlc/` is missing after `git clone` of the product:

```bash
git clone https://github.com/asubbot/ai-sdlc.git ai-sdlc
git -C ai-sdlc checkout "$(tr -d '[:space:]' < ai-sdlc.version)"
make build
```

---

## Consumption model (workspace-only)

Projects consume this repository as a separate workspace root or as a **nested clone** under the product root (`ai-sdlc/`). To keep reproducibility:

1. Each consumer project stores an `ai-sdlc.version` file with a pinned tag or commit SHA.
2. CI in consumer projects verifies that the pinned revision exists in this canonical repository.
3. Process changes are made here via PR; product repositories update only their pinned version.

### Pin file example (`ai-sdlc.version`)

```text
v1.0.1
```

Or pin an exact commit:

```text
f0264d2063e2d881250e0db0542f0de3dfd9c413
```

### Consumer CI example (GitHub Actions)

Verify the pin exists in the canonical repository (set `AI_SDLC_REPO` to your fork or upstream, e.g. `ORG/ai-sdlc`):

```yaml
- name: Verify ai-sdlc pin
  env:
    AI_SDLC_REPO: ORG/ai-sdlc
  run: |
    PIN="$(tr -d '[:space:]' < ai-sdlc.version)"
    if git ls-remote "https://github.com/${AI_SDLC_REPO}.git" "refs/tags/${PIN}" | grep -q .; then
      echo "Tag ${PIN} found"
    elif git ls-remote "https://github.com/${AI_SDLC_REPO}.git" "${PIN}" | grep -q .; then
      echo "Commit ${PIN} found"
    else
      echo "Pin ${PIN} not found in ${AI_SDLC_REPO}" >&2
      exit 1
    fi
```

Product CI template: [templates/consumer/.github/workflows/ci.yml](templates/consumer/.github/workflows/ci.yml) (pin verify, checkout `ai-sdlc` at pin, `make build`, `make check`, `make validate`). Process template regression runs in this repository via [tools/scripts/bootstrap-smoke.sh](tools/scripts/bootstrap-smoke.sh).
