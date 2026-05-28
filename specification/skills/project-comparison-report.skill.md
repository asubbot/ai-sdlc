---
name: project-comparison-report
description: Produce analytical comparison reports between an external project (e.g. open-source repo) and the consumer product design. Use when the user asks to analyse a project, compare a project with the product, generate an analytics report on a repo, or research a similar/competing project.
---

# Project comparison report (analytics)

**Output:** ai-sdlc-artefacts/analytics/<target-project-slug>/<report-name>.md  
**Reference example:** [picoclaw-analysis.md](../../../ai-sdlc-artefacts/analytics/picoclaw/picoclaw-analysis.md)

---

## 1. Goal and scope

Produce a **single markdown report** that:

- Describes the **external project** (architecture, components, security, reliability, engineering culture) from its **source code and docs**.
- **Compares** it with the **consumer product** baseline: both **design** (epic scope, requirements, system design) and **actual implementation** (product source code in the workspace). State gaps, overlaps, and recommendations for the product.
- Stays **factual and traceable**: date of analysis, exact repository URL and **commit/revision** for both the external repo and the PA repo.

**When to use:** User asks to analyse a project, compare a repo with the product, or produce an analytics report on a similar/competing project (e.g. “сделай отчёт по проекту X”, “сравни PicoClaw с нашим дизайном”, “research project Y”).

---

## 2. Inputs

- **Target project:** Repository URL (e.g. GitHub); optionally a specific branch, tag, or commit. If the user did not specify a revision, use the default branch and record the commit hash at the time of analysis.
- **Product baseline (design):** Epic to compare against (user specifies the epic). Required artefacts:
  - ai-sdlc-artefacts/epics/<epic-id>/ep-scope.md
  - ai-sdlc-artefacts/epics/<epic-id>/ep-requirements.md
  - ai-sdlc-artefacts/epics/<epic-id>/ep-system-design.md
- **Product codebase:** The consumer product repo (workspace). The report must analyse **product source code** (package layout, entrypoints, message/request flow, security and reliability mechanisms, tests, CI) — not only the design artefacts — so that comparisons reflect “PA as implemented” as well as “PA as designed”.
- **Product repo revision:** At analysis time, record the current git commit of the product repo (e.g. `git rev-parse HEAD`) and include it in the report header.

---

## 3. Workflow

1. **Resolve target repo** — Clone or use the provided path to the external repo. Resolve the revision (branch/tag/commit) and record the exact commit hash and, if applicable, a direct link to that commit (e.g. `https://github.com/org/repo/commit/<hash>`).
2. **Resolve product baseline** — Confirm ep-scope, ep-requirements, and ep-system-design exist for the chosen epic. If not, ask the user which epic to use or create a minimal baseline.
3. **Analyse product codebase** — Inspect the product repo: package/directory layout, entrypoints (e.g. `cmd/`), main flow (adapters → core → LLM/tools), config loading, security (sandbox, exec, allowlists, secrets), reliability (validation, fallbacks, tests), and CI/quality. Use this to compare “PA as implemented” with the external project; where PA code is missing or differs from the design, note it in the report.
4. **Analyse target codebase** — Inspect layout (packages, modules, entrypoints), data flow (e.g. message bus, agent loop), config, security (sandbox, exec, allowlists), reliability (config validation, fallbacks, tests), and engineering culture (CONTRIBUTING, CI, quality gates). Use the target’s README, CONTRIBUTING, and code; avoid inventing details.
5. **Draft report** — Compose the report in chat (or write to file if the user prefers). Use the [output structure](#4-output-structure) below. Include Mermaid diagrams where they clarify architecture or comparison. Comparisons must reflect both PA design (artefacts) and PA implementation (code) where applicable.
6. **Write to artefact path** — Save the report under `ai-sdlc-artefacts/analytics/<target-project-slug>/<report-name>.md` only after the user approves (e.g. “save”, “lgtm”, “write”). Use a short slug for the project (e.g. `picoclaw`, `openclaw`). Use a short report name (e.g. `picoclaw-analysis.md`). Fix relative links to product artefacts (e.g. `../../epics/EP-104/ep-requirements.md`) so they resolve from the report’s directory.

**Options when in doubt:** If multiple valid choices exist (e.g. epic to compare against, depth of security analysis, which sections to include), present options and ask the user to choose. Do not change source files in the target repo; analysis is read-only.

---

## 4. Output structure

Use the following section layout. Adapt section titles and depth to the target project; omit sections that do not apply.

| Section | Content |
|--------|--------|
| **Header** | Report title; **date of analysis**; **external repo** URL and **revision analysed** (commit hash + link); **PA revision analysed** (commit hash). **Purpose** and **PA design reference** (links to ep-scope, ep-requirements, ep-system-design). |
| **1. High-level architecture** | How the external project is structured (processes, entrypoints, main components). One or more Mermaid diagrams (e.g. flowchart) for data flow. |
| **2. Package layout and module boundaries** | Directory/package structure; dependency rules (if any). Comparison with the product design and, where PA code exists, with the product implementation (e.g. internal/core, adapters). |
| **3. Message processing / main flow** | How requests or messages are handled (e.g. bus → agent loop → LLM/tools). Comparison with the product design and PA code (actual flow). |
| **4. Security** | Access control, sandbox, exec restrictions, allowlists, secrets handling. Comparison with the product requirements and with the product code (as implemented). |
| **5. Reliability and error handling** | Config validation, fallbacks, logging, tests. Comparison with the product strategy and implementation. |
| **6. Component comparison table** | Table: aspect vs external project vs PA (design and/or implementation, as applicable). |
| **7. Flow comparison (Mermaid)** | Sequence or flowchart for external project and for the product (design and/or code). |
| **8. Security analysis** | Deeper security section if relevant: assets, trust boundaries, gaps vs PA (e.g. redaction, LLM audit). Optional Mermaid. |
| **9. Summary** | Short bullets: architecture, security, reliability, fit for the product. |
| **10. Engineering culture comparison** | Process, requirements, quality gates, testing, docs. Optional **code quality** subsection and comparison table. |
| **11. Recommendations for the product** | Features or practices from the external project worth considering for the product **without reducing reliability or security**; grouped by theme. Optional “what not to adopt”. |

**Diagrams:** Use Mermaid (flowchart, sequenceDiagram) in fenced code blocks. Avoid reserved tokens in labels (e.g. no participant named `Loop` if a `loop` block is used; prefer `Agent`). Use relative links only to paths under `ai-sdlc-artefacts/`; ensure every linked file exists.

**Language:** All section titles and body text in English.

---

## 5. Conventions

- **Links:** All links to product artefacts must be **relative** from the report file. From `ai-sdlc-artefacts/analytics/<slug>/report.md`, use `../../epics/<epic-id>/ep-requirements.md` etc.
- **Revisions:** Always record and display the exact commit for the **external** repo and for the **PA** repo so the report is reproducible.
- **Neutral tone:** Describe the external project accurately; compare with the product without implying one is “better”. Recommendations are additive and conditional (“could consider”, “if PA adds X”).
- **No edits to target repo:** Do not modify the external project’s code or config; analysis is read-only.
- **User approval before write:** Do not create or overwrite the report file until the user explicitly approves the draft (e.g. “save”, “lgtm”, “write to file”).

---

## 6. Reference

For a full example of structure, sections, and tone, see:  
[ai-sdlc-artefacts/analytics/picoclaw/picoclaw-analysis.md](../../../ai-sdlc-artefacts/analytics/picoclaw/picoclaw-analysis.md).
