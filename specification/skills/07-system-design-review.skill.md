---
name: system-design-review
description: Review system design documents for SDLC epics (stage 7). Use when reviewing ep-system-design.md files, checking architecture quality, requirement traceability, or when the user asks for architecture or system design review.
---

# Stage 7: System design review

**Pipeline:** [pipeline.spec.md](../pipeline.spec.md)  
**Output:** ai-sdlc-artefacts/epics/<epic-id>/ep-system-design-review.md with YAML front matter, Current Gate Summary, and review iterations (see § Report Template below).

This skill guides systematic review of system design documents (`ep-system-design.md`) within the SDLC pipeline.

## Orchestrator brief (subagent mode)

When launched as a subagent by the pipeline orchestrator ([pipeline.spec.md](../pipeline.spec.md) §3, §4):

- **Required input:** epic ID, paths to `ep-scope.md`, `ep-requirements.md`, `ep-acceptance-criteria.md`, `ep-system-design.md`
- **Context:** `ep-context.md` (read for orientation, but do not use instead of source artefacts)
- **Gate check before launch:** `ep-system-design.md` must exist
- **Output signal:** `STAGE_7_COMPLETE: ai-sdlc-artefacts/epics/<epic-id>/ep-system-design-review.md [gate=<pass|fail|cap>, iteration <N>, blocker:<n> major:<n> medium:<n> minor:<n>]`
- **Validation after:** orchestrator checks gate status from signal; if fail → return to stage 6

## Mandatory delegation (pipeline stage 7)

When this skill is run as **pipeline stage 7** for an epic, execution MUST follow [pipeline.spec.md](../pipeline.spec.md) **§3**:

- **If you are the orchestrator** (you just helped author `ep-system-design.md` or earlier stages): **do not** perform this review yourself in the same session. **Delegate** to a **subagent** (Cursor Task / equivalent) or start a **new chat** with fresh context and a one-line brief: epic id, paths under `ai-sdlc-artefacts/epics/EP-XXX/`, instruction to run this skill end-to-end.
- **If you are the delegated reviewer**: treat inputs as read-only; produce and save `ep-system-design-review.md` when running inside the HOTL pipeline so the gate has durable evidence. Outside an orchestrated pipeline, draft in chat unless the user asked to save. Do **not** edit `ep-context.md`; the orchestrator applies any accepted gate summary update after the review is accepted or saved.

## When to Use

- Reviewing `ai-sdlc-artefacts/epics/EP-XXX/ep-system-design.md`
- User asks for architecture or design review
- Before implementation planning (stage 8)

## Design–review iteration ([pipeline.spec.md](../pipeline.spec.md) §2.1)

Stages **6** and **7** repeat until **zero** open findings in **Blocker**, **Major**, **Medium**, and **Minor**, or until the **operator decides** after the cap.

1. **Count iterations** — Each completed save of a **`## Review iteration N`** section in `ep-system-design-review.md` is one stage 7 iteration. **N** must not exceed **5** without an operator decision recorded in chat or in the review file (e.g. under the latest iteration).
2. **Single file** — Use one `ep-system-design-review.md` per epic. For iteration **N**, add (or replace only if the user agrees to discard a draft) a **top-level** heading `## Review iteration N` with a stable increasing **N**. **Retain** all prior `## Review iteration …` sections for history.
3. **Current Gate Summary** — On save, update YAML front matter, Current Gate Summary, and the latest review iteration atomically. The latest full `## Review iteration N` remains the source of truth.
4. **Exit loop** — After this iteration’s findings are recorded, set **`Iteration summary — open counts`** for Blocker / Major / Medium / Minor. If all four are **zero**, the iteration loop is **complete**; stage 8 may follow.
5. **Cap** — If **N = 5** and any **Blocker**, **Major**, **Medium**, or **Minor** count is still **> 0**, **stop**: list remaining issues and require an **operator decision** before further stage 6/7 work or stage 8.
6. **Return to stage 6** — When Blocker/Major/Medium/Minor > 0 and **N < 5**, the orchestrator runs **stage 6** again to revise `ep-system-design.md`, then runs **stage 7** again (new **delegated** session per pipeline §3).

## Review Workflow

### Step 1: Read Related Documents

Read in order:
1. YAML front matter and ep-context.md, if present, to understand the current epic state and identify areas that may have changed since the last review iteration
2. `ep-scope.md` — understand feature scope and glossary
3. `ep-requirements.md` — verify all requirements are addressed
4. `ep-acceptance-criteria.md` — ensure testability alignment
5. `ep-system-design.md` — the document under review

Do not use ep-context.md as a replacement for the source artefacts in this review stage.

### Step 2: Structural Check

Verify the design document contains:
- [ ] Overview with scope reference
- [ ] Architecture diagram (C4 C2 or equivalent)
- [ ] Module boundaries table
- [ ] Components and interfaces table
- [ ] Data models
- [ ] Error handling
- [ ] Testing strategy
- [ ] Risks and trade-offs
- [ ] Requirement traceability table

### Step 3: Requirement Traceability

For each requirement in `ep-requirements.md`:
- Verify explicit coverage in traceability table
- Confirm design component is identified
- Check acceptance criteria alignment

### Step 4: Architecture Quality

Assess:
- **KISS**: Is the solution as simple as possible?
- **Fail fast**: Are errors caught early with clear messages?
- **Security**: Are security controls adequate?
- **Testability**: Can components be tested in isolation?
- **Modularity**: Are boundaries clear and dependencies minimal?
- **Architecture patterns**: Does every architecturally significant decision carry an `architecture-pattern:` field with `<pattern-id>` or `n/a — <one-line reason>` in the **Design decisions** section of `ep-system-design.md`? (A record present only in `ep-context.md` does not count; that file is not a source of truth for this review.) For `<pattern-id>`: chosen / rejected / why is recorded and the decision does not contradict the card's `when_not` / `kiss_default` (catalog: `reference/architecture-patterns/` at the ai-sdlc checkout root). For `n/a`: the reason must hold — not a brush-off where a card's `when` clearly applies. A violation of any of these conditions is a **Medium** finding (see Step 5). Additionally, for each `architecture-pattern: <pattern-id>` (not `n/a`): confirm the Design decisions entry records **chosen / rejected / why** and an upstream **https://** link; when the choice rejects applying the pattern, the **why** should be consistent with the card's `kiss_default` or `when_not` (Medium if missing or contradictory). If `ai-sdlc-artefacts/architecture-patterns-playbook.md` exists and the decision contradicts the playbook default without an epic-local explanation in **why**, treat as **Medium** (documentation/consistency), not Blocker.

### Step 5: Identify Issues

Categorize every finding into **one** severity (definitions align with pipeline **§2.1** exit counts):

| Severity | Criteria |
|----------|----------|
| **Blocker** | Missing requirement coverage, unacceptable security gap, data loss or integrity risk, design that cannot meet must-have REQ/AC |
| **Major** | Wrong or missing component/contract, traceability break, missing error-handling strategy for required flows, testability blocker |
| **Medium** | Unclear interfaces, incomplete non-critical specs, inconsistent structure, gaps that should be fixed before implementation (including a missing or violated `architecture-pattern:` record on an architecturally significant decision, per Step 4) |
| **Minor** | Documentation polish, optional consistency, low-risk improvements |

**Loop exit:** **Blocker = 0 AND Major = 0 AND Medium = 0 AND Minor = 0** (only open items in this iteration; resolved items from prior iterations do not need re-listing unless regressed).

### Step 6: Output Report

Generate or update **ep-system-design-review.md** in the same epic folder as `ep-system-design.md` (unless the user agrees another path). Add **`## Review iteration N`** per **Design–review iteration** above and refresh the Current Gate Summary at the top. Follow the template below. On first review, **N = 1**; on subsequent passes after stage 6 fixes, increment **N** (max **5** without operator decision).

---

## Report Template

Use **one** `ep-system-design-review.md` per epic. **First** iteration: create the file with YAML front matter, document title, Current Gate Summary, and `## Review iteration 1`. **Later** iterations: update YAML front matter and Current Gate Summary, then **append** a new top-level `## Review iteration N` block at the end; **do not remove** prior iteration sections.

```markdown
---
artefact: ep-system-design-review
epic_id: EP-XXX
status: draft
source_of_truth: true
gate: fail
latest_iteration: N
open_counts:
  blocker: X
  major: X
  medium: X
  minor: X
next_action: return_to_stage_6
updated_at: YYYY-MM-DD
---

# Architecture Review — EP-XXX [optional title]

**Reviewer:** [AI Agent / Name]

---

## Current Gate Summary

Gate: Pass | Fail | Cap
Latest iteration: N
Last updated: YYYY-MM-DD
Open counts: Blocker X | Major X | Medium X | Minor X
Open findings:
- F-001 Major: One-line finding title.
- M1 Medium: DD-3 uses architecture-pattern: rate-limiting but omits upstream https link and kiss_default rationale.
Next action: Proceed to stage 8 | Return to stage 6 | Operator decision required

---

## Review iteration N

**Review date:** YYYY-MM-DD
**Stage 7 iteration:** N of max 5
**Document reviewed:** [ep-system-design.md](ep-system-design.md)
**Iteration summary — open counts:** Blocker: X | Major: X | Medium: X | Minor: X
**Gate:** Pass (Blocker/Major/Medium/Minor all zero) | Fail (any Blocker/Major/Medium/Minor > 0) | Cap (N = 5 and Blocker/Major/Medium/Minor still > 0 — operator decision required)

### Overall assessment

[2–3 sentences for this iteration]

**Verdict:** [Pass gate / Fail gate / Cap — stop for operator]

### Strengths

- [Specific strength with reference]

### Issues and recommendations

#### Blocker

| # | Issue | Context | Recommendation |
|---|-------|---------|----------------|

#### Major

| # | Issue | Context | Recommendation |
|---|-------|---------|----------------|

#### Medium

| # | Issue | Context | Recommendation |
|---|-------|---------|----------------|

#### Minor

| # | Issue | Context | Recommendation |
|---|-------|---------|----------------|

### Architectural decisions (optional)

| Decision | Justification |
|----------|---------------|
| | |

### NFR coverage (optional)

| NFR | Coverage | Status |
|-----|----------|--------|
| | | OK / Needs work |

### Project rules compliance (optional)

| Rule | Compliance |
|------|------------|
| KISS | ✅ / ❌ / ⚠️ |
| Fail fast | ✅ / ❌ / ⚠️ |
| Security | ✅ / ❌ / ⚠️ |
| Testability | ✅ / ❌ / ⚠️ |

### Traceability (this iteration)

- **Architecture:** [ep-system-design.md](ep-system-design.md)
- **Requirements:** [ep-requirements.md](ep-requirements.md)
- **Acceptance criteria:** [ep-acceptance-criteria.md](ep-acceptance-criteria.md)
- **Scope:** [ep-scope.md](ep-scope.md)
```

---

## Checklist

Before completing review:
- [ ] Iteration number **N** is set (1–5); prior iteration sections preserved when **N > 1**
- [ ] **Iteration summary — open counts** filled for Blocker / Major / Medium / Minor
- [ ] YAML front matter, Current Gate Summary, and latest `## Review iteration N` are consistent
- [ ] Gate reflects **§2.1** (Pass only if Blocker/Major/Medium/Minor all zero; Cap if **N = 5** and any of those > 0)
- [ ] All requirements have traceability entries for this iteration
- [ ] Blocker / Major / Medium / Minor issues have clear recommendations
- [ ] Report follows template structure under `## Review iteration N`
- [ ] Severity ratings are consistent with the definitions in Step 5
- [ ] Action items are specific and actionable
- [ ] Delegated reviewer did not edit ep-context.md directly; orchestrator owns that update
