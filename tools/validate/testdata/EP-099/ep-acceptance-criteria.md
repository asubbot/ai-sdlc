---
artefact: ep-acceptance-criteria
epic_id: EP-099
status: approved
source_of_truth: true
updated_at: 2026-05-20
---

# EP-099: Acceptance Criteria

## AC Index

| AC | REQ | Summary |
|----|-----|---------|
| [AC-99.001](#ac-99001) | REQ-99.001 | Fixture loading works |
| [AC-99.002](#ac-99002) | REQ-99.002 | Golden comparison works |
| [AC-99.003](#ac-99003) | REQ-99.003 | EARS patterns validated |
| [AC-99.004](#ac-99004) | REQ-99.004 | Structure sections checked |
| [AC-99.005](#ac-99005) | REQ-99.005 | Pipeline gates parsed |
| [AC-99.006](#ac-99006) | REQ-99.006 | Broken links detected |

## Scenarios

### AC-99.001 Fixture loading (Trace: REQ-99.001)

Given a testdata directory with fixture markdown files
When the validator loads the fixtures
Then all files are parsed without errors

### AC-99.002 Golden file comparison (Trace: REQ-99.002)

Given a golden JSON file in testdata/golden/
When the validator produces JSON output
Then the output matches the golden file byte-for-byte

### AC-99.003 EARS compliance (Trace: REQ-99.003)

Given a requirements file with EARS-formatted requirements
When the EARS linter runs
Then valid EARS patterns are accepted and invalid ones are flagged

### AC-99.004 Structure validation (Trace: REQ-99.004)

Given an artefact markdown file
When the structure validator runs
Then missing required sections are reported as errors

### AC-99.005 Pipeline gate checking (Trace: REQ-99.005)

Given an epic directory with artefact front matter
When the pipeline gate checker runs
Then gate status is correctly parsed from YAML and gate summaries

### AC-99.006 Broken link detection (Trace: REQ-99.006)

Given a markdown file with internal links
When the structure validator checks links
Then links pointing to non-existent files are reported
