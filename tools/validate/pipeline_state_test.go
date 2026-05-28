package main

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestParseGateSummary(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantErr     bool
		wantBlocker int
		wantMajor   int
	}{
		{
			"pass gate",
			"## Current Gate Summary\n" +
				"| Category | Blocker | Major | Medium | Minor | Nit |\n" +
				"|----------|---------|-------|--------|-------|-----|\n" +
				"| Count    | 0       | 0     | 0      | 0     | 1   |",
			false, 0, 0,
		},
		{
			"fail gate",
			"## Current Gate Summary\n" +
				"| Category | Blocker | Major | Medium | Minor | Nit |\n" +
				"|----------|---------|-------|--------|-------|-----|\n" +
				"| Count    | 0       | 2     | 1      | 0     | 0   |",
			false, 0, 2,
		},
		{
			"plain-text pass",
			"## Current Gate Summary\n\n" +
				"Gate: Pass\n" +
				"Latest iteration: 2\n" +
				"Last updated: 2026-05-28\n" +
				"Open counts: Blocker 0 | Major 0 | Medium 0 | Minor 0\n" +
				"Next action: Proceed to stage 8\n",
			false, 0, 0,
		},
		{
			"plain-text fail",
			"## Current Gate Summary\n\n" +
				"Gate: Fail\n" +
				"Latest iteration: 1\n" +
				"Last updated: 2026-05-28\n" +
				"Open counts: Blocker 1 | Major 3 | Medium 2 | Minor 1\n" +
				"Open findings:\n" +
				"- F-001 Blocker: Missing requirement coverage.\n" +
				"Next action: Return to stage 6\n",
			false, 1, 3,
		},
		{"no gate table", "# Just a heading\nSome text", true, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs, err := parseGateSummary(tt.content)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gs.Blocker != tt.wantBlocker {
				t.Errorf("Blocker = %d, want %d", gs.Blocker, tt.wantBlocker)
			}
			if gs.Major != tt.wantMajor {
				t.Errorf("Major = %d, want %d", gs.Major, tt.wantMajor)
			}
		})
	}
}

func TestCheckPipelineState_HealthyEpic(t *testing.T) {
	result := checkPipelineState("testdata/EP-099", "EP-099")
	if result.HasGaps {
		t.Errorf("expected no gaps for healthy epic, got errors=%d findings=%v", result.Errors, result.Findings)
	}
}

func TestCheckPipelineState_BrokenEpic(t *testing.T) {
	result := checkPipelineState("testdata/EP-098", "EP-098")
	if !result.HasGaps {
		t.Error("expected gaps for broken epic")
	}
	if result.Errors == 0 {
		t.Error("expected errors")
	}
}

func TestCheckPipelineState_OrderViolation(t *testing.T) {
	dir := t.TempDir()
	epicDir := filepath.Join(dir, "EP-TEST")
	os.MkdirAll(epicDir, 0o755)

	os.WriteFile(filepath.Join(epicDir, "ep-scope.md"),
		[]byte("---\nartefact: ep-scope\nepic_id: EP-TEST\nstatus: approved\nupdated_at: 2026-05-20\n---\n# Scope"), 0o644)

	os.WriteFile(filepath.Join(epicDir, "ep-implementation-plan.md"),
		[]byte("---\nartefact: ep-implementation-plan\nepic_id: EP-TEST\nstatus: approved\nupdated_at: 2026-05-20\n---\n# Plan"), 0o644)

	result := checkPipelineState(epicDir, "EP-TEST")
	if !result.HasGaps {
		t.Error("expected gaps for order violation")
	}
}

func TestCheckPipelineState_GateViolation(t *testing.T) {
	dir := t.TempDir()
	epicDir := filepath.Join(dir, "EP-GATE")
	os.MkdirAll(epicDir, 0o755)

	for _, f := range []struct{ name, content string }{
		{"ep-scope.md", "---\nartefact: ep-scope\nepic_id: EP-GATE\nstatus: approved\nupdated_at: 2026-05-20\n---\n"},
		{"ep-requirements.md", "---\nartefact: ep-requirements\nepic_id: EP-GATE\nstatus: approved\nupdated_at: 2026-05-20\n---\n"},
		{"ep-acceptance-criteria.md", "---\nartefact: ep-acceptance-criteria\nepic_id: EP-GATE\nstatus: approved\nupdated_at: 2026-05-20\n---\n"},
		{"ep-system-design.md", "---\nartefact: ep-system-design\nepic_id: EP-GATE\nstatus: approved\nupdated_at: 2026-05-20\n---\n"},
		{"ep-system-design-review.md", "---\nartefact: ep-system-design-review\nepic_id: EP-GATE\nstatus: in_progress\ngate: fail\nupdated_at: 2026-05-20\n---\n## Current Gate Summary\n| Category | Blocker | Major | Medium | Minor | Nit |\n|----------|---------|-------|--------|-------|-----|\n| Count    | 0       | 2     | 0      | 0     | 0   |\n"},
		{"ep-implementation-plan.md", "---\nartefact: ep-implementation-plan\nepic_id: EP-GATE\nstatus: approved\nupdated_at: 2026-05-20\n---\n"},
	} {
		os.WriteFile(filepath.Join(epicDir, f.name), []byte(f.content), 0o644)
	}

	result := checkPipelineState(epicDir, "EP-GATE")
	if !result.HasGaps {
		t.Error("expected gaps: stage 8 present but stage 7 gate=fail")
	}
}

func TestCheckPipelineState_GateOverrideDecisionAllowsProgression(t *testing.T) {
	dir := t.TempDir()
	epicDir := filepath.Join(dir, "EP-OVERRIDE")
	os.MkdirAll(epicDir, 0o755)

	for _, f := range []struct{ name, content string }{
		{"ep-scope.md", frontMatter("ep-scope", "EP-OVERRIDE")},
		{"ep-requirements.md", frontMatter("ep-requirements", "EP-OVERRIDE")},
		{"ep-acceptance-criteria.md", frontMatter("ep-acceptance-criteria", "EP-OVERRIDE")},
		{"ep-system-design.md", frontMatter("ep-system-design", "EP-OVERRIDE")},
		{"ep-system-design-review.md", reviewArtefact("ep-system-design-review", "EP-OVERRIDE", "fail", "return_to_stage_6", 0, 1, 0, 0) +
			"\nDecision needed: review-gate override\nOperator choice: accept residual risk\nRationale: tracked for follow-up.\n"},
		{"ep-implementation-plan.md", frontMatter("ep-implementation-plan", "EP-OVERRIDE")},
	} {
		os.WriteFile(filepath.Join(epicDir, f.name), []byte(f.content), 0o644)
	}

	result := checkPipelineState(epicDir, "EP-OVERRIDE")
	if result.HasGaps {
		t.Fatalf("expected recorded operator decision to allow progression, got findings=%v", result.Findings)
	}
}

func TestCheckPipelineState_CapGateRequiresDecisionEvidence(t *testing.T) {
	dir := t.TempDir()
	epicDir := filepath.Join(dir, "EP-CAP")
	os.MkdirAll(epicDir, 0o755)

	for _, f := range []struct{ name, content string }{
		{"ep-scope.md", frontMatter("ep-scope", "EP-CAP")},
		{"ep-requirements.md", frontMatter("ep-requirements", "EP-CAP")},
		{"ep-acceptance-criteria.md", frontMatter("ep-acceptance-criteria", "EP-CAP")},
		{"ep-system-design.md", frontMatter("ep-system-design", "EP-CAP")},
		{"ep-system-design-review.md", reviewArtefact("ep-system-design-review", "EP-CAP", "cap", "operator_decision_required", 0, 0, 1, 0)},
		{"ep-implementation-plan.md", frontMatter("ep-implementation-plan", "EP-CAP")},
	} {
		os.WriteFile(filepath.Join(epicDir, f.name), []byte(f.content), 0o644)
	}

	result := checkPipelineState(epicDir, "EP-CAP")
	if !result.HasGaps {
		t.Error("expected cap gate without operator decision evidence to block progression")
	}
}

func TestCheckPipelineState_CodeReviewGateBlocksAudit(t *testing.T) {
	dir := t.TempDir()
	epicDir := filepath.Join(dir, "EP-CODE")
	os.MkdirAll(epicDir, 0o755)

	for _, f := range []struct{ name, content string }{
		{"ep-scope.md", frontMatter("ep-scope", "EP-CODE")},
		{"ep-requirements.md", frontMatter("ep-requirements", "EP-CODE")},
		{"ep-acceptance-criteria.md", frontMatter("ep-acceptance-criteria", "EP-CODE")},
		{"ep-system-design.md", frontMatter("ep-system-design", "EP-CODE")},
		{"ep-system-design-review.md", reviewArtefact("ep-system-design-review", "EP-CODE", "pass", "proceed_to_stage_8", 0, 0, 0, 0)},
		{"ep-implementation-plan.md", frontMatter("ep-implementation-plan", "EP-CODE")},
		{"ep-code-review.md", reviewArtefact("ep-code-review", "EP-CODE", "fail", "return_to_stage_9", 0, 0, 0, 1)},
		{"ep-audit-report.md", frontMatter("ep-audit-report", "EP-CODE")},
	} {
		os.WriteFile(filepath.Join(epicDir, f.name), []byte(f.content), 0o644)
	}

	result := checkPipelineState(epicDir, "EP-CODE")
	if !result.HasGaps {
		t.Error("expected stage 10 gate fail to block stage 11")
	}
}

func TestCheckPipelineState_AuditRequiresCodeReviewEvidence(t *testing.T) {
	dir := t.TempDir()
	epicDir := filepath.Join(dir, "EP-MISSING-REVIEW")
	os.MkdirAll(epicDir, 0o755)

	for _, f := range []struct{ name, content string }{
		{"ep-scope.md", frontMatter("ep-scope", "EP-MISSING-REVIEW")},
		{"ep-requirements.md", frontMatter("ep-requirements", "EP-MISSING-REVIEW")},
		{"ep-acceptance-criteria.md", frontMatter("ep-acceptance-criteria", "EP-MISSING-REVIEW")},
		{"ep-system-design.md", frontMatter("ep-system-design", "EP-MISSING-REVIEW")},
		{"ep-system-design-review.md", reviewArtefact("ep-system-design-review", "EP-MISSING-REVIEW", "pass", "proceed_to_stage_8", 0, 0, 0, 0)},
		{"ep-implementation-plan.md", frontMatter("ep-implementation-plan", "EP-MISSING-REVIEW")},
		{"ep-audit-report.md", frontMatter("ep-audit-report", "EP-MISSING-REVIEW")},
	} {
		os.WriteFile(filepath.Join(epicDir, f.name), []byte(f.content), 0o644)
	}

	result := checkPipelineState(epicDir, "EP-MISSING-REVIEW")
	if !result.HasGaps {
		t.Error("expected stage 11 without stage 10 gate evidence to fail")
	}
}

func frontMatter(artefact, epic string) string {
	return "---\nartefact: " + artefact + "\nepic_id: " + epic + "\nstatus: approved\nsource_of_truth: true\nupdated_at: 2026-05-20\n---\n"
}

func reviewArtefact(artefact, epic, gate, nextAction string, blocker, major, medium, minor int) string {
	return "---\nartefact: " + artefact + "\nepic_id: " + epic + "\nstatus: approved\nsource_of_truth: true\ngate: " + gate + "\nlatest_iteration: 1\nopen_counts:\n  blocker: " +
		strconv.Itoa(blocker) + "\n  major: " + strconv.Itoa(major) + "\n  medium: " + strconv.Itoa(medium) + "\n  minor: " + strconv.Itoa(minor) + "\nnext_action: " + nextAction + "\nupdated_at: 2026-05-20\n---\n\n" +
		"## Current Gate Summary\n\nGate: " + gate + "\nLatest iteration: 1\nLast updated: 2026-05-20\nOpen counts: Blocker " +
		strconv.Itoa(blocker) + " | Major " + strconv.Itoa(major) + " | Medium " + strconv.Itoa(medium) + " | Minor " + strconv.Itoa(minor) + "\n"
}
