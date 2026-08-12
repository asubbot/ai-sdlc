# Validation Tool

Multi-purpose validation tool for SDLC pipeline.

## Subcommands

The CLI supports subcommand dispatch. If the first non-flag argument is a known subcommand, it is dispatched. Otherwise the tool defaults to `ac` (backwards compatible).

```
validate [subcommand] [EP-XXX] [--json]
```

For this canonical `ai-sdlc` repository, invoke the binary from repo root as `./tools/validate/validate ...`.  
Examples below use the canonical repository-root path `./tools/validate/validate`.

| Subcommand  | Description | Status |
|-------------|-------------|--------|
| `ac`        | AC coverage validation (default) | Implemented |
| `req`       | REQ ↔ AC traceability check (`req all` skips NEW/CANCEL epics) | Implemented |
| `pipeline`  | Pipeline state and gate validation | Implemented |
| `structure` | Artefact structure validation | Implemented |
| `ears`      | EARS requirements linting (`ears all` skips NEW/CANCEL epics) | Implemented |

The `--json` flag works in any position (before or after subcommand/epic).

### In-scope epics (`ears all` / `req all`)

Project-wide `ears all` and `req all` validate only epics whose `ep-scope.md` **Status** is not **NEW** or **CANCEL**/**CANCELED** (parsed from `| **Status** | … |` table row). Skipped epics are listed in the report but do not affect the exit code.

In-scope epics must have `ep-requirements.md` (ears) and both `ep-requirements.md` and `ep-acceptance-criteria.md` (req). REQ headings must use the canonical form `### REQ-EE.NNN — Summary` (see stage 4 skill).

## Pipeline State Validation

The `pipeline` subcommand validates one epic's stage ordering and HOTL gate evidence:

```bash
./tools/validate/validate pipeline EP-009
./tools/validate/validate pipeline EP-009 --json
```

It enforces:

- later stages cannot exist before required earlier artefacts;
- stage 8 cannot exist unless stage 7 has `gate: pass`, or the stage 7 review artefact records an operator decision for a non-pass gate;
- stage 11 cannot exist unless stage 10 code-review gate evidence exists and has `gate: pass`, or the code-review artefact records an operator decision for a non-pass gate;
- open **Blocker**, **Major**, **Medium**, or **Minor** counts make a review gate blocking. `Nit` and `Suggestion` are non-blocking.

Accepted operator decision evidence is the decision record pattern from `pipeline.spec.md`. The validator currently requires both lines below in the affected review artefact:

```markdown
Decision needed: <type>
Operator choice: <selected option>
```

This is intentionally minimal: the validator checks that decision evidence exists, while process review checks whether the decision is appropriate.

Only **`ENOENT`** (`fs.ErrNotExist`) is treated as **missing** for pipeline artefacts. Any other read failure (for example a directory at an artefact path, permissions, or I/O failure) is a **hard error** on that stage and reports the stage file plus the underlying error text.

For **`ep-implementation-plan.md`**, unchecked `- [ ]` tasks are counted from the **raw readable file content** even if YAML front matter is malformed. If stage 10 or stage 11 exists and the implementation plan cannot be read, checkbox evaluation does **not** silently skip: the pipeline validator raises a hard error instead.

**Exit codes:** `validate pipeline` sets exit **1** only when `Errors > 0`. **Warnings** (for example unchecked implementation-plan tasks when `ep-code-review.md` exists but `ep-audit-report.md` does not) do not change the exit code. Consumer CI that must fail on warnings needs a separate policy (for example parse JSON output); there is no `--strict` flag in v1.0.1.

## Artefact Structure Validation

The `structure` subcommand validates YAML front matter, required sections, and relative links for epic artefacts under `ai-sdlc-artefacts/epics/<epic-id>/`.

```bash
./tools/validate/validate structure EP-009
./tools/validate/validate structure EP-009 --json
```

For **`ep-context.md`**, required sections (heading text match, case-insensitive): **Purpose**, **Current Scope**, **Open Questions**, **Links**. Additional sections (Key Requirements, Acceptance Signals, Design Decisions, etc.) are recommended in skills but not enforced by CI.

Review artefacts (`ep-system-design-review.md`, `ep-code-review.md`) require **Current Gate Summary** and at least one **Review iteration** section.

## Enforcement model

Normative **process** rules (HOTL/HITL, stages, delegation) live in [specification/pipeline.spec.md](../../specification/pipeline.spec.md). The table below is the **authoritative CI matrix** for `./tools/validate/validate`; the spec §4.5 table mirrors it for agents.

| Requirement | Enforcement | Tool / evidence |
|-------------|-------------|-----------------|
| Stage ordering and required artefacts | Hard | `validate pipeline`, file presence, front matter |
| Review gate pass before downstream progression | Hard | Current Gate Summary, open severity counts, `validate pipeline` |
| AC↔test coverage before epic completion | Hard | `validate` / `validate ac` |
| Qualifying orphan AC trace comments | Hard | `validate` / `validate ac`, contextual `path:line` error (including traces attached to receiver-method `Test*` declarations) |
| Test-source parse/read/walk failures during AC scan | Hard | `validate` / `validate ac`, parse/read/walk error from scanned `*_test.go` inputs |
| Artefact structure (sections, links, front matter) | Hard | `validate structure` |
| EARS requirements format | Hard | `validate ears` / `validate ears all` |
| REQ↔AC traceability | Hard | `validate req` / `validate req all` |
| Artefact read failures other than ENOENT | Hard | `validate pipeline`, contextual stage read error |
| Unchecked implementation-plan tasks before audit | Hard | `validate pipeline` (when `ep-audit-report.md` exists; raw plan body still checked when readable) |
| Unreadable implementation plan for downstream checkbox evaluation | Hard | `validate pipeline` (when stage 10 or 11 exists) |
| Unchecked tasks when code review exists | Warning | `validate pipeline` (when `ep-code-review.md` exists and the plan is readable) |
| Mandatory delegation for stages 7 and 10 | Soft | Orchestrator run notes, subagent output signal |
| Required HITL decision quality | Soft | Decision record in chat or artefact |
| `ep-context.md` staleness judgment | Soft | Agent reads source artefacts when context is stale |
| Subagent discipline | Soft | Process compliance; not CI-verifiable |

**Stage 9:** Task execution produces the codebase, not a markdown artefact. Completion is inferred from `ep-implementation-plan.md` checkboxes (`- [x]` / `- [ ]` in `## Tasks`) when later stages exist.

## AC Validation

The `validate` tool automatically validates that all Acceptance Criteria (AC) from an epic's `ep-acceptance-criteria.md` are covered by tests.

### Purpose

Before completing an epic's audit, use this tool to ensure:
- ✅ Every in-scope AC-EE.NNN has traceability from a test comment (see [Test coverage declaration](#test-coverage-declaration)) or is explicitly marked **Obsolete** or **Deferred** in `ep-acceptance-criteria.md` (see [Excluding ACs from automation](#4-excluding-acs-from-automation-deferred-or-obsolete))
- ✅ AC codes in tests are found via the supported comment shapes (`covers` / `supporting`, EP-N AC-, label form, REQ+AC on the same line, etc.)
- ✅ No AC is silently missed without traceability or a documented exclusion (**Obsolete** / **Deferred** / manual-only)
- ✅ Every top-level `func Test…` under `tests/`, `internal/`, and `cmd/` has at least one **AC trace line** bound to it (see [Test functions must declare AC trace](#test-functions-must-declare-ac-trace-reverse-check))

### Usage

#### Build

```bash
cd tools/validate
go build -o validate .
cd ../..
```

#### Validate All Epics (Default)

```bash
./tools/validate/validate
```

#### Validate Single Epic

```bash
./tools/validate/validate EP-009
```

Output:
```
🔍 Validating AC coverage for all 9 epics...

📋 Epic Validation Summary

Epic       Trace%       Status
────────────────────────────────────
✓ EP-001        93%
✗ EP-004        82%
...

────────────────────────────────────

❌ OVERALL: in-scope 96/111 traced (86.5%), automated 96 (86.5%), manual-only 0 | deferred 1 | obsolete 1 | total ACs 113
   Project-wide: Test functions with t.Skip: 0

❌ AC not covered by tests (project-wide): 15

EP-009
  • AC-09.001
  ...

Tip: run `./tools/validate/validate EP-XXX` for per-AC detail and test refs.

❌ Test functions without AC trace comment (project-wide): N
  • internal/foo/bar_test.go::TestBaz
  ...

Action: Add a trace line (e.g. `// Covers AC-EE.NNN`) bound to each `Test*` per this document.
```

**Trace%** is traceability **in scope** (ACs that still require test traces): `(automated + manual-only) / in_scope`. **Deferred** and **Obsolete** ACs are **not** counted in the numerator; they reduce `in_scope` instead of inflating the percentage.

When a one-line criterion is parsed from `ep-acceptance-criteria.md` (not a markdown table row), it may appear after `—` on each bullet.

Use this for:
- Project health check
- Identifying which epics need work
- Seeing the full list of AC ids still missing `Covers AC-…` traceability
- Tracking overall AC coverage trend
- CI/CD dashboard reporting

#### Single Epic (Detailed View)

Output (table format; human mode prints the banner — JSON mode prints JSON only):
```
🔍 Validating AC coverage for EP-009...

📋 AC Coverage Report for EP-009

AC Code         Criterion                                          Coverage
───────────────────────────────────────────────────────────────────────────────────────────────
✓ AC-09.008                                                        5 tests
✓ AC-09.009                                                        3 tests
✗ AC-09.001                                                        NOT COVERED
↷ AC-09.005                                                        OBSOLETE
✎ AC-09.007                                                        MANUAL …
...

⚠️ RESULT: in-scope 15/16 traced (93.8%), automated 14 (87.5%), manual-only 1 | deferred 0 | obsolete 1 | total ACs 18
   Project-wide: Test functions with t.Skip: 3

❌ Missing coverage for:
  • AC-09.001
  • AC-09.006
  ...

Action: Add tests for missing ACs, or mark them **Obsolete** / **Deferred** in ep-acceptance-criteria.md (see below)

❌ Test functions without AC trace comment (project-wide): …
  • …
```

#### JSON output: Parse results programmatically

```bash
./tools/validate/validate --json EP-009
./tools/validate/validate --json
```

**Single epic** output:
```json
{
  "epic": "EP-009",
  "total_acs": 18,
  "deferred_acs": 0,
  "obsolete_acs": 1,
  "in_scope_acs": 17,
  "automated_covered_acs": 14,
  "manual_only_traced_acs": 1,
  "traceability_ratio": 0.8824,
  "automated_ratio": 0.8235,
  "test_funcs_with_skip": 3,
  "tests_missing_ac_trace": ["internal/foo/bar_test.go::TestBaz"],
  "gaps": [
    {"code": "AC-09.001", "criterion": "", "status": "not_covered"},
    {"code": "AC-09.005", "criterion": "", "status": "obsolete", "reason": "Obsolete in ep-acceptance-criteria.md"},
    ...
  ],
  "ac_to_tests": {
    "AC-09.008": [
      {"ref": "internal/tools/create_tool_test.go::TestFunc1", "manual": false}
    ]
  }
}
```

**All epics** (`./tools/validate/validate --json`) includes the same aggregate fields (`in_scope_acs`, `traceability_ratio`, `automated_ratio`, `test_funcs_with_skip`, …), plus `not_covered_acs` (flat list with `epic`, `code`, optional `criterion`) and `not_covered_count`, and **`tests_missing_ac_trace`**: a sorted list of `path/to/file_test.go::TestName` for top-level tests missing an AC trace (same scan roots as coverage).

For **`./tools/validate/validate EP-XXX --json`**, `tests_missing_ac_trace` is still the **full repository** list (not limited to ACs belonging to that epic), so CI and local runs can fix every stray `Test*` in one pass.

## Metrics

For each epic (and for the all-epics JSON aggregate):

| Field | Meaning |
|-------|---------|
| `in_scope_acs` | `total_acs - deferred_acs - obsolete_acs` — ACs that still require test traceability. |
| `deferred_acs` | ACs marked **Deferred** (or `MANUAL ONLY` / `**Status:** … Deferred …`) near the AC in `ep-acceptance-criteria.md`. |
| `obsolete_acs` | ACs marked **Obsolete** (or `**Status:** … Obsolete …`) near the AC — criteria superseded by product refactors. |
| `automated_covered_acs` | In-scope ACs with at least one **non-manual** test reference. |
| `manual_only_traced_acs` | In-scope ACs where **only** manual references exist (see below). |
| `traceability_ratio` | `(automated_covered_acs + manual_only_traced_acs) / in_scope_acs` — deferred and obsolete are **not** in the numerator. |
| `automated_ratio` | `automated_covered_acs / in_scope_acs`. |
| `test_funcs_with_skip` | Project-wide count of `Test*` functions whose body contains `t.Skip` (direct call on `t`); scanned under `tests/`, `internal/`, `cmd/`. |
| `tests_missing_ac_trace` | (JSON only, when non-empty) Sorted `rel/path_test.go::TestName` entries for top-level `Test*` functions without a bound AC trace line (see below). |

## Test functions must declare AC trace (reverse check)

In addition to **AC → tests**, the tool enforces **test → AC**: every top-level `func Test\w+` in scanned `*_test.go` files must have at least one qualifying trace line **bound** to that function by the shared Go AST index. A qualifying trace binds only when it is either:

1. part of the **doc comment** for a top-level `func Test…`; or
2. located **inside that same top-level `Test*` body**.

Only **actual parsed `//` comment lines** participate in binding. Text inside raw string literals that merely starts with `//`, or `/* ... */` block-comment lookalikes, is ignored.

**`TestMain` is excluded.** Method `Test*` declarations (`func (suite) TestX(...)`) are excluded from both valid binding and reverse-check enforcement. A qualifying trace attached to such a receiver method is treated as an **orphan trace** and hard-fails validation; it is **not** silently ignored. `Benchmark*`, `Example*`, and `Fuzz*` are not checked.

A line counts as an **AC trace** only if **both** hold:

1. `lineDeclaresACCoverage(line)` is true (same rules as the coverage scanner — e.g. `covers` / `supporting`, `// AC-EE.NNN:`, EP+AC lines, `REQ-` + AC on the same line).
2. The line contains at least one parseable **`AC-EE.NNN`** (`extractACsFromLine` is non-empty).

So a comment like `// Covers integration` **without** an `AC-EE.NNN` code does **not** satisfy the reverse check, even though it contains the word “covers”.

Compatibility impact: comments merely **between functions**, after a test body, or attached to a non-top-level-test declaration no longer count. These orphan traces now fail validation instead of degrading into `::unknown`. Existing valid **doc-comment** traces remain valid, and valid **inline** traces remain valid and bind to their **enclosing** top-level test. Removing pseudo-traces from raw strings or block-comment lookalikes can also expose an AC that now needs a real trace comment.

## Test Coverage Declaration

### Automatic vs manual traceability

A test reference is **manual** if either:

- The traceability line contains the whole word `manual` (case-insensitive), e.g. `// manual Covers AC-01.004`, or
- The `Test*` function that owns the trace line contains a direct **`t.Skip(...)`** call in that top-level test body (parsed via Go AST). If both apply, manual wins for that reference.

If an AC has **only** manual references, it counts toward `manual_only_traced_acs` and `traceability_ratio`, but **not** toward `automated_ratio`. If an AC has **any** non-manual reference, it counts as **automated** for `automated_ratio`.

`t.Skip(...)` inside nested func literals or subtests does **not** mark the enclosing top-level `Test*` as skipped/manual. Repositories that previously relied on nested `t.Skip` may therefore see `test_funcs_with_skip` and manual/automated traceability metrics change.

Trace lines are bound through the shared Go AST index. A qualifying line counts only when it falls within the doc-comment range of a top-level `func Test…` or within that function body. Inline traces bind to the enclosing test body. Comments stranded between functions, attached to other declarations, or attached to method `Test*` declarations are orphan traces and fail validation.

Malformed Go test files, orphan traces, unreadable test files, and walk failures are hard validation errors. The scan fails fast on the first structural/orphan error, so consumers should fix that reported issue and re-run. Where the parser can identify a source location, the error includes contextual `path:line`; validation exits with status **1**.

### Format: Comment before test function

Mark which ACs your test covers with a comment directly above (or before) the test function:

```go
// Covers AC-09.008: create_tool accepts required parameters
func TestCreateToolTool_Run_success(t *testing.T) {
    // ...test code...
}
```

### Supported formats

Single AC:
```go
// Covers AC-09.001
func TestSomething(t *testing.T) { }
```

Multiple ACs (comma-separated):
```go
// Covers AC-09.001, AC-09.002
func TestMultipleACs(t *testing.T) { }
```

Range of ACs (using en-dash or hyphen):
```go
// Covers AC-09.008–013
func TestRangeOfACs(t *testing.T) { }
```

Supporting test (non-primary coverage):
```go
// Supporting AC-09.001: helper test
func TestHelper(t *testing.T) { }
```

Mixed:
```go
// Covers AC-09.001, AC-09.003–005, AC-09.010
func TestMixed(t *testing.T) { }
```

### Epic-prefixed manual test files

Operator scenarios live in `ai-sdlc-artefacts/epics/EP-XXX/ep-manual-tests.md` (or `ep-manual-test-scenarios.md` for EP-001). To anchor those ACs in code **without** mixing with automated tests, use dedicated files under `tests/integration/`:

- `ep001_manual_test.go`, `ep004_manual_test.go`, `ep009_manual_test.go` (EP-002 and EP-006 use automated-only traces; no dedicated `ep002` / `ep006` manual files)

> **Note:** These test files live in the **consumer product repository** (e.g. `tests/integration/`), not in this canonical `ai-sdlc` repo.

Conventions: `//go:build integration`, `package integration_test`, `// manual Covers AC-…` on the trace line, `t.Skip("manual: …")` with a pointer to the epic manual doc (and optional anchor). `./tools/validate/validate` reads these files like any other `*_test.go` under `tests/`.

## Integration Points

### Before Git Commit (Optional)

Add to `settings.json` (hooks) to validate before every commit:
```json
{
  "hooks": {
    "before_git_commit": "epic=$(git rev-parse --abbrev-ref HEAD | grep -o 'EP-[0-9]*'); if [ -n \"$epic\" ]; then ./tools/validate/validate \"$epic\"; fi"
  }
}
```

Or validate manually before committing:
```bash
./tools/validate/validate EP-009
git commit -m "feat(EP-009): implement create_tool..."
```

### In CI/CD

Example GitHub Actions:
```yaml
- name: Validate AC coverage
  run: |
    ./tools/validate/validate --json EP-009 > /tmp/report.json
    coverage=$(jq '.traceability_ratio' /tmp/report.json)
    if (( $(echo "$coverage < 1.0" | bc -l) )); then
      echo "❌ Not all ACs covered"
      exit 1
    fi
```

### In SDLC Pipeline

See [Stage 9 (Task Execution)](../../specification/skills/09-task-execution.skill.md), [Stage 10 (Code review)](../../specification/skills/10-code-review.skill.md), and [Stage 11 (Audit)](../../specification/skills/11-audit.skill.md).

## Architecture

**Location:** `ai-sdlc/tools/validate/` (multi-purpose validation tool)

| File | Role |
|------|------|
| `main.go` | CLI dispatch, epic scan, subcommand routing |
| `ac_parse.go`, `ac_coverage.go`, `ac_report.go` | AC parsing, coverage scan, human/JSON reports |
| `pipeline_state.go` | `validate pipeline` |
| `artefact_structure.go` | `validate structure` |
| `req_ac_trace.go`, `ears_lint.go` | `validate req`, `validate ears` |
| `ast_skip.go`, `test_ac_trace.go`, `policy_nolint_gocyclo.go`, `output.go` | test parsing, reverse AC trace, gocyclo policy, stdout helpers |

**Tests:** `*_test.go` across the package; epic fixtures under `testdata/EP-099` (healthy) and `testdata/EP-098` (broken). Run `go test ./...` from `tools/validate/`.

## Building

```bash
# Build
cd tools/validate && go build -o validate . && cd ../..

# Run tests
cd tools/validate && go test ./...
```

## Exit Codes

- **0** — Requested validation passed: for `ac`, every in-scope AC is traced, every scanned `Test*` has an AC trace line, and the AGENTS.md gocyclo-suppression policy scan is clean; for `pipeline`, stage ordering and gate evidence are valid ✅
- **1** — Failure: missing AC traceability, a `Test*` without a bound AC trace comment, an orphan AC trace comment, malformed or unreadable scanned Go input, a gocyclo policy violation, an artefact structure issue, a pipeline ordering/gate violation, an artefact read failure beyond ENOENT, or missing operator decision evidence ❌

## Common Workflows

### 1. Project Health Check

```bash
# Quick overview of all epics
./tools/validate/validate

# Output shows which epics need work
```

### 2. During Development (Task Execution)

```bash
# Add test, mark it with "// Covers AC-09.001"
# Build and check coverage for specific epic
cd tools/validate && go build -o validate . && cd ../..
./tools/validate/validate EP-009

# If incomplete: add more tests or mark ACs Obsolete/Deferred in ep-acceptance-criteria.md
# If complete: ready for code review and audit (stages 10–11)
```

### 3. Before Audit

```bash
# Last check before epic completion
./tools/validate/validate EP-009

# If gaps exist: list them for manual review/deferral
# If all covered: proceed to stages 10–11 (code review, then audit)
```

### 4. Excluding ACs from automation (deferred or obsolete)

Use this when an AC must **not** fail validation for lack of `// Covers AC-…` tests.

- **Deferred:** work postponed or validated outside unit tests (e.g. bootstrap gate, manual-only). Document near the AC using `**Deferred:**`, `MANUAL ONLY`, `DEFERRED`, or `**Status:** … Deferred …` on a line that also references the same `AC-EE.NNN` (or a nearby `**Status:**` line, same heuristic as before).
- **Obsolete:** the criterion no longer applies after a **vision or refactor change** (superseded behaviour). Document near the AC using `**Obsolete:**`, `OBSOLETE`, or `**Status:** … Obsolete …` together with that `AC-EE.NNN` within a few lines (index table row or criterion heading).

The tool then marks the AC as **↷ DEFERRED** or **↷ OBSOLETE** in the report and does **not** require a `Covers AC-…` line for that AC.

**Important:** Do not mention other epics’ `AC-EE.NNN` codes inside an epic’s `ep-acceptance-criteria.md` except as real ACs for that file — the parser extracts every `AC-\d{2}\.\d{3}` substring.

Optional: add a normal test comment if you still want traceability for partial automation — the validator does **not** treat `// Obsolete AC-…` in Go as a markdown exclusion (use `ep-acceptance-criteria.md` markers above).

## Troubleshooting

### No coverage found even though tests exist

1. Check test file is under `tests/`, `internal/`, or `cmd/` (only `*_test.go` files are scanned).
2. Use a traceability comment the tool recognizes, for example:
   - `// Covers AC-XX.YYY` or `// Supporting AC-…` (case-insensitive `covers` / `supporting` is OK)
   - `// EP-008 AC-08.001 / REQ-08.001: …`
   - `// AC-04.025: …` (label form)
   - `// … (AC-06.005, … AC-06.010 / REQ-06.013)` when `REQ-` appears on the same line
3. Ensure AC codes are comma/range-separated as documented; ranges like `AC-09.008–013` are supported.
4. Run debug: `grep -rE 'Covers AC-|EP-[0-9]+ AC-|AC-[0-9]{2}\\.[0-9]{3}:' tests/ internal/ cmd/`

### AC code not found in markdown

The parser matches `AC-EE.NNN` in many line shapes, including `**AC-09.001**`, `### AC-09.001`, and `[AC-09.001](...)`. Prefer the same bold form as in your epic template for consistency.

### JSON output is invalid or mixed with log lines

In JSON mode, **stdout** is only the JSON document (no banner lines). Use `--json` before or after the epic id:
```bash
./tools/validate/validate --json EP-009
./tools/validate/validate EP-009 --json
```
Diagnostics and usage errors go to **stderr**.

