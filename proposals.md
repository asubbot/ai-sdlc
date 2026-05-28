# ai-sdlc Pipeline Proposals

This document captures proposed improvements to the `ai-sdlc` pipeline. The main goal is to reduce token usage while preserving traceability, review quality, and human-in-the-loop control.

The proposal has three optimization layers:

1. YAML front matter: cheap file metadata and routing state.
2. Current Gate Summary: current review-gate state.
3. `ep-context.md`: compact semantic context for the epic.

## 1. Add an Epic Context Artefact

Introduce `ep-context.md` as a short epic-level artefact:

```text
ai-sdlc-artefacts/epics/<epic-id>/ep-context.md
```

Purpose: provide the next pipeline stage with the current working context without requiring the agent to read every upstream artefact in full.

`ep-context.md` is not a source of truth. If it conflicts with a full artefact, the full artefact wins.

Recommended structure:

```markdown
# Epic Context — EP-XXX

## Purpose
## Current Scope
## Key Requirements
## Acceptance Signals
## Design Decisions
## Interfaces / Contracts
## Current Gate Summary
## Open Questions
## Links
```

Keep the file short, ideally 100-150 lines. Details that are not needed by the next stage should stay in the full artefacts.

Staleness safeguard: if `ep-context.md` has YAML `updated_at` older than any source artefact it summarizes, agents must treat it as stale and open the changed source artefacts before relying on the context. The `Open Questions` section should contain unresolved questions from previous stages plus questions intentionally handed off to the next stage.

## 2. Read Context Before Full Artefacts

Update stage skills to read `ep-context.md` first when it exists.

Full upstream artefacts should be opened only when needed for:

- traceability checks;
- resolving ambiguity or contradiction;
- filling missing context;
- reviewing material changes;
- validating links or gate status.

This keeps full documents as the source of truth while making the normal path cheaper.

This rule is strongest for authoring and consuming stages, such as stages 6, 8, 10 when epic context is optional, and 11. Verification stages must not use `ep-context.md` as a replacement for source artefacts:

- Stage 7 system design review must still read `ep-scope.md`, `ep-requirements.md`, `ep-acceptance-criteria.md`, and `ep-system-design.md` for traceability and quality checks.
- Stage 10 code review should remain diff-first. When reviewing against an epic, it should use `ep-context.md` to identify likely relevant REQ/AC, then open the corresponding source sections as needed.

## 3. Update `ep-context.md` After Key Stages

`ep-context.md` should be refreshed after stages that materially change epic knowledge:

| Stage | Update |
|---|---|
| 3. Epic planning | Create initial purpose, scope, links |
| 4. Requirements | Add key requirements |
| 5. Acceptance criteria | Add acceptance signals |
| 6. System design | Add design decisions and contracts |
| 7. System design review | Update design gate summary and open findings |
| 8. Implementation planning | Add plan-level execution notes if useful |
| 9. Task execution | Update only if implementation causes material design or contract changes |
| 10. Code review | Update code review gate summary when epic-scoped |
| 11. Audit | Use the context as the fast entry point |

The update should be concise. Do not copy full sections from upstream artefacts.

Responsibility: the agent or orchestrator that owns the stage updates `ep-context.md`. Delegated readonly reviewers for stages 7 and 10 should not edit it directly; the orchestrator applies the resulting gate summary after the review is accepted or saved.

## 4. Add Current Gate Summaries to Review Files

Review files currently retain all iterations, which is useful for history but expensive for downstream stages.

Add a short current-state block at the top of review artefacts:

```markdown
## Current Gate Summary

Gate: Pass / Fail / Cap
Latest iteration: N
Last updated: YYYY-MM-DD
Open counts: Blocker X | Major X | Medium X | Minor X
Non-blocking counts: Nit X | Suggestion X
Open findings:
- F-001 Major: One-line finding title.
Next action: Proceed to stage X / Return to stage Y / Operator decision required.
```

The summary is a derived current-state view. The latest full `## Review iteration N` remains the source of truth; if the two disagree, update the summary from the latest iteration.

Update rules:

1. Add or update the latest `## Review iteration N`.
2. Recalculate open counts from that iteration.
3. Refresh `Current Gate Summary` at the top of the file.
4. Keep all previous review iterations for history.

Stage 7 and stage 10 skills must update the latest review iteration and `Current Gate Summary` atomically on save. `Open findings` should list finding IDs plus severity and a one-line title, so downstream agents can decide whether they need to read the full finding.

Downstream stages should read this block first. Older iterations should be read only when the summary reports open findings, references repeated issues, or a traceability dispute requires the full review history.

Apply this to:

- `ep-system-design-review.md`;
- `ep-code-review.md` when saved.

## 5. Add YAML Front Matter to Pipeline Artefacts

Add a small YAML front matter block to pipeline artefacts that need fast status checks. This gives agents and future validation tools a cheap machine-readable entry point before reading the Markdown body.

Example:

```yaml
---
artefact: ep-system-design-review
epic_id: EP-001
status: draft
source_of_truth: true
gate: fail
latest_iteration: 3
open_counts:
  blocker: 0
  major: 1
  medium: 0
  minor: 2
non_blocking_counts:
  nit: 0
  suggestion: 0
next_action: return_to_stage_6
updated_at: 2026-04-28
---
```

Recommended use:

- all epic artefacts should identify `artefact`, `epic_id`, `status`, `source_of_truth`, and `updated_at`;
- review artefacts should also expose `gate`, `latest_iteration`, `open_counts`, optional `non_blocking_counts`, and `next_action`;
- `ep-context.md` should set `source_of_truth: false`;
- full source artefacts such as `ep-requirements.md`, `ep-acceptance-criteria.md`, and `ep-system-design.md` should set `source_of_truth: true`.
- project-level artefacts such as `scope.md` and `strategy.md` should identify `artefact`, `status`, `source_of_truth`, and `updated_at`.

Recommended `status` values: `draft`, `approved`, `superseded`, `archived`.

YAML front matter must be updated in the same operation as the body content it reflects. If YAML front matter is absent, agents must fall back to reading the body content directly.

YAML front matter does not replace `ep-context.md` or `Current Gate Summary`. It is the cheapest metadata layer:

1. YAML front matter answers "what is this file and what is its current status?"
2. `Current Gate Summary` answers "did the review gate pass and what remains open?"
3. `ep-context.md` answers "what is the current working context of the epic?"

Agents should read YAML front matter first for routing decisions, then `Current Gate Summary` or `ep-context.md` only when they need semantic context.

## 6. Strengthen Diff-First Code Review

Stage 10 already prefers minimal context. Make this rule stricter:

1. Read the diff first.
2. Read touched files only where needed.
3. Read immediate callers/callees only when the diff does not explain behavior.
4. Use `ep-context.md` Key Requirements and Acceptance Signals to identify relevant REQ/AC.
5. Use `ep-context.md` as the first epic-aware review input.
6. Open full epic artefacts only for the related requirements, acceptance criteria, or design sections.

This prevents code review from loading the full epic history for small changes.

## 7. Suggested Minimal Implementation

Start with a small process-only change:

1. Add `ep-context.md` to epic-level artefact naming in `pipeline.spec.md`.
2. Document the `ep-context.md` purpose and structure.
3. Add YAML front matter guidance for pipeline artefacts.
4. Update stage skills 3-8, 10, and 11 to read/update `ep-context.md`.
5. Add `Current Gate Summary` guidance to stage 7 and 10 review skills.
6. Update `pipeline.spec.md` traceability (§5) and summary diagram (§6) to include `ep-context.md`.
7. Add migration behavior: for existing epics without `ep-context.md`, agents should create it on first encounter during stages 3-11.
8. Defer validation tooling until the process proves useful.

## Estimated Token Savings

These are directional estimates, not benchmark results. They are based on the current pipeline shape, where later stages may read the full upstream chain, and on a target flow where agents read `ep-context.md` first and open full artefacts only for traceability or disputed details.

| Area | Expected saving | Basis |
|---|---:|---|
| Stage 6. System design | 20-35% | Reads compact scope, key REQ, and AC summaries before opening full requirement files selectively |
| Stage 7. System design review | 10-25% | Still reads full source artefacts for verification, but avoids full review-history reads through YAML metadata and Current Gate Summary |
| Stage 8. Implementation planning | 40-60% | Avoids loading full scope, requirements, AC, design, and all review iterations on every planning pass |
| Stage 10. Code review against an epic | 30-50% | Uses diff-first review plus `ep-context.md`; opens only touched REQ/AC/design sections |
| Stage 11. Audit | 40-70% | Uses gate summaries and links, then opens only artefacts needed to verify status |

The largest savings come from two patterns:

- replacing repeated full-document reads with one short `ep-context.md`;
- replacing full review-history reads with `Current Gate Summary` plus latest unresolved findings.

The savings will be smaller for very small epics and larger for epics with long requirements, multiple design-review iterations, or saved code-review histories.

## Measurement Plan

Use two complementary metrics:

1. `actual_total_tokens`: real token usage from Cursor Usage.
2. `read_context_tokens_estimate`: estimated tokens in files or sections read by the agent.

`actual_total_tokens` is the primary metric. It can be collected from the Cursor Usage dashboard by isolating a pipeline stage run:

1. Record the stage, mode, model, start time, and end time.
2. Run only one pipeline stage during that interval.
3. Avoid parallel agents or unrelated chats while measuring.
4. Sum all Usage rows in the interval.
5. Compare baseline and optimized runs for the same stage and similar input artefacts.

Example measurement table:

```markdown
| Run | Stage | Mode | Model | Start | End | Usage rows | Total tokens | Quality result |
|---|---|---|---|---|---|---:|---:|---|
| 1 | Stage 8 | baseline | gpt-5.5-medium | 12:40 | 12:45 | 3 | 342000 | Plan accepted |
| 2 | Stage 8 | optimized | gpt-5.5-medium | 13:00 | 13:03 | 2 | 118000 | Plan accepted |
```

Use `read_context_tokens_estimate` as the explanatory metric. It shows why token usage changed by tracking which artefacts were read and estimating their size. A simple estimate is `tokens ~= characters / 4` for English Markdown. This is not billing-accurate, but it is useful for comparing baseline and optimized context loading.

Success criteria should include both cost and quality:

- Stage 8: at least 30% lower `actual_total_tokens` without reducing implementation-plan quality.
- Stage 10: at least 20% lower `actual_total_tokens` for epic-aware review without missing Blocker or Major findings.
- Stage 11: at least 40% lower `actual_total_tokens`.
- No incorrect decisions caused by stale `ep-context.md`, stale YAML front matter, or outdated `Current Gate Summary`.

## Expected Benefits

- Lower token usage in stages 6-11.
- Faster handoff between agents and delegated reviewers.
- Less repeated reading of full requirements, acceptance criteria, designs, and review histories.
- Better current-state visibility without losing full traceability.

## 8. Align Human-in-the-loop with enforceability reality

Current documentation states Human-in-the-loop as a strict MUST behavior, while practical enforcement is mixed: some rules are machine-checkable, others are process-only. This proposal makes that explicit without weakening operator control.

### 8.1 Reframe Human-in-the-loop as policy + controls

Define Human-in-the-loop in two layers:

1. **Process policy (mandatory):** operator choice is required on defined decision points.
2. **Enforcement model (explicit):**
   - **Hard controls:** validated by tooling/CI;
   - **Soft controls:** validated by process discipline and review, not by CI.

This avoids over-claiming technical enforcement while keeping the rule normative.

### 8.1.1 HITL vs HOTL definitions for this pipeline

To avoid ambiguity, define both terms explicitly:

- **Human-in-the-loop (HITL):** blocking human participation in the decision path. The agent must stop and get operator choice before proceeding.
- **Human-on-the-loop (HOTL):** supervisory human oversight. The agent may proceed autonomously; the operator monitors and intervenes when needed.

Operational distinction:

- HITL is used where risk or irreversibility is high.
- HOTL is used for routine execution between explicit decision gates.

Recommended model for `ai-sdlc`: **hybrid control** — HITL on required decision points, HOTL for non-decision-path execution and routine validation steps.

### 8.2 Introduce required decision points

Replace broad wording like "when multiple valid choices exist" with a concrete list of decision points where agent must stop and ask the operator before proceeding.

Suggested required decision points:

- conflict resolution between competing artefact sources;
- review-gate override (stage 7 or 10) when blocking severities remain;
- iteration-cap override after bounded retry loops;
- skipping a prerequisite stage or proceeding with missing required input;
- migration or cutover strategy when multiple materially different paths exist.

For non-listed ambiguities, the orchestrator may proceed with defaults and log rationale.

### 8.3 Add an enforceability matrix to the spec

Add a table in `pipeline.spec.md`:

| Requirement | Enforcement | Evidence |
|---|---|---|
| Stage ordering and required artefacts | Hard | `validate pipeline`, file presence, front matter checks |
| Gate pass before stage progression | Hard | gate summary/open counts + validator |
| Subagent usage for stages 7/10 | Soft (currently) | run notes, operator review |
| Operator choice at decision points | Soft (currently) | decision log entries |

This makes expectations auditable and truthful.

### 8.4 Standardize decision logging

Add a minimal decision record format for any required decision point:

```markdown
Decision needed: <type>
Context: <one-line why>
Options: A | B | C
Operator choice: <selected option>
Rationale: <short>
```

Where to store:

- in chat output for transient decisions;
- in the affected epic artefact section when decision impacts long-lived traceability;
- optionally in a dedicated decision log file if frequent.

### 8.5 Documentation changes (minimal and recommended)

**Minimal (fast, low-risk):**

1. In `pipeline.spec.md`, clarify that Human-in-the-loop is policy-level and only partially machine-verifiable.
2. In `skills/README.md`, replace absolute wording with "required decision points" wording.
3. In `tools/validate/VALIDATION.md`, explicitly list what validators do and do not enforce.

**Recommended (stronger governance):**

1. Add the decision-point list and enforceability matrix to `pipeline.spec.md`.
2. Update stage skills 6-11 to use a consistent "Decision needed" output pattern.
3. Require decision-log evidence for gate overrides and iteration-cap overrides.

### 8.6 Success criteria

- No contradiction between normative docs and actual enforcement capability.
- Operators can clearly distinguish auto-validated vs process-validated requirements.
- Review and audit stages can verify that required operator decisions were recorded.
