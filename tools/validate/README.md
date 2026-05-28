# validate — SDLC Validation Tool

Multi-purpose validation tool for the `ai-sdlc` pipeline.

## CLI

```
validate [subcommand] [EP-XXX] [--json]
```

### Subcommands

| Subcommand  | Description | Status |
|-------------|-------------|--------|
| `ac`        | AC coverage validation (default) | ✅ Implemented |
| `req`       | REQ ↔ AC traceability check | ✅ Implemented |
| `pipeline`  | Pipeline state and gate validation | ✅ Implemented |
| `structure` | Artefact structure validation | ✅ Implemented |
| `ears`      | EARS requirements linting | ✅ Implemented |

If no subcommand is given, or the first argument is not a known subcommand (e.g. `EP-XXX`), the tool defaults to `ac` (backwards compatible).

### Examples

```bash
./tools/validate/validate                    # AC coverage for all epics (from repo root)
./tools/validate/validate EP-009             # AC coverage for single epic
./tools/validate/validate ac EP-009          # Same as above (explicit)
./tools/validate/validate req EP-009         # REQ-AC traceability
./tools/validate/validate ears EP-009        # EARS linter
./tools/validate/validate --json             # JSON output (any subcommand)
./tools/validate/validate req EP-009 --json  # JSON for specific subcommand
```

## Current Validators

### Policy: no gocyclo suppressions (AGENTS.md)

Scans all `*.go` files under `tests/`, `internal/`, and `cmd/` for the forbidden golangci-lint cyclomatic-complexity suppression (substring `nolint:gocyclo` **outside** double-quoted string literals, so tests that assert on source text are not false positives). Violations are listed as `path/to/file.go:LINE`. On failure, human mode prints a dedicated block; JSON includes `nolint_gocyclo_violations` and sets `has_gaps` to true (single-epic JSON includes the same field on the epic report object).

### AC (Acceptance Criteria) Validation — `ac` subcommand

Validates that all Acceptance Criteria from an epic's `ep-acceptance-criteria.md` are covered by tests (with separate metrics for **automated** vs **manual-only** traceability; deferred ACs do not inflate the traceability percentage). It also checks the **reverse**: every top-level `Test*` under `tests/`, `internal/`, and `cmd/` must have at least one trace line that both matches the coverage declaration rules **and** contains a real `AC-EE.NNN` code bound to that test (see [VALIDATION.md](./VALIDATION.md#test-functions-must-declare-ac-trace-reverse-check)).

```bash
cd tools/validate
go build -o validate .
cd ../..
./tools/validate/validate              # all epics (default)
./tools/validate/validate EP-009       # single epic
./tools/validate/validate ac EP-009    # explicit subcommand
./tools/validate/validate --json       # JSON output
```

Output (human mode) includes **Trace%** per epic (in-scope traceability), an **OVERALL** line with `in-scope … traced`, **automated** / **manual-only** counts, **deferred**, **Project-wide: Test functions with t.Skip** (count of `Test*` bodies with `t.Skip` in scanned trees), and — on failure — a list of **`path/to/file_test.go::TestName`** entries for `Test*` functions missing an AC trace.

```
🔍 Validating AC coverage for all 9 epics...

📋 Epic Validation Summary

Epic       Trace%       Status
────────────────────────────────────
✓ EP-001        93%
✗ EP-009        61%
────────────────────────────────────

❌ OVERALL: in-scope 96/111 traced (86.5%), automated 96 (86.5%), manual-only 0 | deferred 2 | total ACs 113
   Project-wide: Test functions with t.Skip: 0

❌ AC not covered by tests (project-wide): 15
...
```

**Coverage declaration (traceability):**

The scanner looks for AC codes on lines that match project conventions — not only `// Covers AC-…` (see [VALIDATION.md](./VALIDATION.md#test-coverage-declaration)): case-insensitive `covers` / `supporting`, `// EP-NNN AC-EE.NNN`, `// AC-EE.NNN:`, lines with `REQ-` and AC codes, etc.

Use the whole word **`manual`** on the line to mark **manual-only** traceability, or put **`t.Skip`** in the `Test*` function body (then all AC refs from that function are treated as manual). Epic operator scenarios may be grouped in **`tests/integration/epXXX_manual_test.go`** (see [VALIDATION.md](./VALIDATION.md#epic-prefixed-manual-test-files)).

Example:

```go
// Covers AC-09.008: create_tool accepts parameters
func TestCreateToolTool_Run(t *testing.T) {
    // ...
}
```

Also supported: ranges (`// Covers AC-09.008–013`), comma-separated ACs, and `Supporting AC-…`.

See [VALIDATION.md](./VALIDATION.md) for full documentation (metrics JSON schema, `t.Skip` behavior, comment-to-function resolution).

## Exit Codes

- **0** — All validations passed (AC coverage, per-`Test*` AC trace, and policy scan) ✅
- **1** — Validation failed (missing AC coverage, `Test*` without bound AC trace, and/or forbidden gocyclo suppression in product trees) ❌

JSON (`--json`): failures set `"has_gaps": true` when any in-scope AC is untraced, when `tests_missing_ac_trace` is non-empty, or when `nolint_gocyclo_violations` is non-empty. The `tests_missing_ac_trace` and `nolint_gocyclo_violations` arrays are always **project-wide**, including when you run `./tools/validate/validate EP-009 --json`.

## Future Validators

- [ ] Design specification consistency
- [ ] Code coverage thresholds
- [ ] Dependency graph validation
