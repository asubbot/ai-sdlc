# ai-sdlc Pipeline Proposals

Staging area for process ideas before they become normative. Nothing here is normative: the process is defined in [specification/pipeline.spec.md](specification/pipeline.spec.md) and the stage skills; shipped changes are recorded in [CHANGELOG.md](CHANGELOG.md).

Keep this file short. When a proposal is adopted, move it to the **Adopted** list below as a one-liner instead of keeping a second copy of the rules; the full text stays in git history (`git log -p -- proposals.md`).

## Open proposals

### 1. Verify the token-savings claims, or drop them

The three compact context layers were adopted on the basis of estimated savings of 20-70% per stage, with success criteria of at least 30% for stage 8, 20% for stage 10, and 40% for stage 11. No measurement has been recorded since v1.0.0.

Options: run one baseline/optimized pair on a comparable real epic (record stage, mode, model, time interval, total tokens from Cursor Usage, and whether the stage output was accepted), or drop the numbers and describe the layers as a design choice rather than a measured optimization.

Earlier estimate tables and the full manual measurement procedure are in this file's history up to v1.0.1.

### 2. Machine-check consistency between the compact layers

`validate` checks that YAML front matter and gate sections exist and parse, and gates pessimistically when blocking counts appear anywhere (`inferGate` in [tools/validate/pipeline_state.go](tools/validate/pipeline_state.go)). It never reports divergence, so two failure modes stay invisible:

1. front matter `open_counts`, `Current Gate Summary`, and the latest `## Review iteration N` disagree while all parsed values look non-blocking;
2. `ep-context.md` has an `updated_at` older than a source artefact it summarizes.

Both are comparisons between values the validator already parses. Implementing them would move `ep-context.md` staleness judgment and gate-summary accuracy from **Soft** to **Hard** in the enforcement matrix in [pipeline.spec.md](specification/pipeline.spec.md) §4.5 and [VALIDATION.md](tools/validate/VALIDATION.md).

## Adopted

Read the linked normative files, not this list.

| Proposal | Where it lives now |
|---|---|
| `ep-context.md` artefact; read context before full artefacts; refresh after key stages | [pipeline.spec.md](specification/pipeline.spec.md) (*Token-optimized context loading*, §4.3), stage skills 03-11, [skills/README.md](specification/skills/README.md) |
| `Current Gate Summary` in review artefacts | [07-system-design-review.skill.md](specification/skills/07-system-design-review.skill.md), [10-code-review.skill.md](specification/skills/10-code-review.skill.md), `validate structure` / `validate pipeline` |
| YAML front matter on pipeline artefacts | [pipeline.spec.md](specification/pipeline.spec.md), [VALIDATION.md](tools/validate/VALIDATION.md) |
| Diff-first code review | [10-code-review.skill.md](specification/skills/10-code-review.skill.md) |
| HOTL by default with required HITL decision points, enforceability matrix, decision record format | [pipeline.spec.md](specification/pipeline.spec.md) (*Required HITL decision points*, §4.5), [VALIDATION.md](tools/validate/VALIDATION.md) |
| Greenfield consumer bootstrap | [00-project-bootstrap.skill.md](specification/skills/00-project-bootstrap.skill.md), [templates/consumer/](templates/consumer/), [CHANGELOG.md](CHANGELOG.md) |
