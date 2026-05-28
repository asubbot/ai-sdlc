---
artefact: ep-acceptance-criteria
epic_id: EP-000
status: approved
source_of_truth: true
updated_at: 2026-05-28
---

# EP-000: Acceptance Criteria

## AC Index

| AC | REQ | Summary |
|----|-----|---------|
| [AC-00.001](#ac-00001) | REQ-00.001 | Pin present |

## Scenarios

### AC-00.001 Pin (Trace: REQ-00.001)

Given the consumer repository
When `ai-sdlc.version` is read
Then it is non-empty

**Deferred:** Bootstrap gate; no product unit test.
