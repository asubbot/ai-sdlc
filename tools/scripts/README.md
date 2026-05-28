# Maintainer scripts (canonical ai-sdlc)

Shell utilities for process maintainers. Not copied to consumer product repositories.

## bootstrap-smoke.sh

Automated regression for greenfield bootstrap: copies [templates/consumer/](../../templates/consumer/) into a temporary consumer layout, then runs product gates `make build`, `make check` (fmt, vet, vuln, lint, race+coverage on `tests/`), and `make validate` (AC all epics + pipeline/structure per epic). This exercises the consumer template; real product repositories enforce the same gates via `.github/workflows/ci.yml`. May require network on first run (govulncheck, golangci-lint via `go run`).

```bash
./tools/scripts/bootstrap-smoke.sh
```

### When to run

- After changes to `templates/consumer/` or [00-project-bootstrap.skill.md](../../specification/skills/00-project-bootstrap.skill.md)
- Before merging PRs that touch consumer templates (also runs in CI)

### Environment variables

| Variable | Purpose | Default |
|----------|---------|---------|
| `AI_SDLC_ROOT` | Canonical repository root | Parent of `tools/` (repo root) |
| `AI_SDLC_PIN` | Force pin written to `ai-sdlc.version` (tag or SHA) | See pin policy below |
| `KEEP_TMP` | Set to `1` to keep the temp product directory and print its path | Temp dir removed on exit |

### Pin policy

Unless `AI_SDLC_PIN` is set:

1. Read the example tag from `templates/consumer/ai-sdlc.version` (first non-comment line).
2. If `AI_SDLC_ROOT` HEAD equals that tag’s commit, write the **tag** (e.g. `v1.0.1`).
3. Otherwise write **HEAD** as a full commit SHA.

This mirrors bootstrap skill step 5 (tag when on a release, SHA for dev branches).

### Exit codes

- `0` — smoke passed
- `1` — missing inputs, build/validate failure, or Done-when assertion failure
