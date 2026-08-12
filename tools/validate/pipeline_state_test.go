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

func TestCheckPipelineState_MalformedReviewFrontMatterStillKeepsGateViolation(t *testing.T) {
	tests := []struct {
		name             string
		epic             string
		reviewFile       string
		reviewContent    string
		downstreamFile   string
		downstreamPrefix string
		gateSnippet      string
		files            []struct{ name, content string }
	}{
		{
			name:             "stage 7 review",
			epic:             "EP-MALFORMED-STAGE7",
			reviewFile:       "ep-system-design-review.md",
			reviewContent:    "---\nartefact: ep-system-design-review\nepic_id: EP-MALFORMED-STAGE7\nstatus: approved\ngate: pass\n",
			downstreamFile:   "ep-implementation-plan.md",
			downstreamPrefix: "stage 8 (ep-implementation-plan.md) exists but stage 7 gate=",
			gateSnippet:      "stage 8 (ep-implementation-plan.md) exists but stage 7 gate=",
			files: []struct{ name, content string }{
				{"ep-scope.md", frontMatter("ep-scope", "EP-MALFORMED-STAGE7")},
				{"ep-requirements.md", frontMatter("ep-requirements", "EP-MALFORMED-STAGE7")},
				{"ep-acceptance-criteria.md", frontMatter("ep-acceptance-criteria", "EP-MALFORMED-STAGE7")},
				{"ep-system-design.md", frontMatter("ep-system-design", "EP-MALFORMED-STAGE7")},
				{"ep-implementation-plan.md", frontMatter("ep-implementation-plan", "EP-MALFORMED-STAGE7")},
			},
		},
		{
			name:             "stage 10 review",
			epic:             "EP-MALFORMED-STAGE10",
			reviewFile:       "ep-code-review.md",
			reviewContent:    "---\nartefact: ep-code-review\nepic_id: EP-MALFORMED-STAGE10\nstatus: approved\ngate: pass\n",
			downstreamFile:   "ep-audit-report.md",
			downstreamPrefix: "stage 11 (ep-audit-report.md) exists but stage 10 gate=",
			gateSnippet:      "stage 11 (ep-audit-report.md) exists but stage 10 gate=",
			files: []struct{ name, content string }{
				{"ep-scope.md", frontMatter("ep-scope", "EP-MALFORMED-STAGE10")},
				{"ep-requirements.md", frontMatter("ep-requirements", "EP-MALFORMED-STAGE10")},
				{"ep-acceptance-criteria.md", frontMatter("ep-acceptance-criteria", "EP-MALFORMED-STAGE10")},
				{"ep-system-design.md", frontMatter("ep-system-design", "EP-MALFORMED-STAGE10")},
				{"ep-system-design-review.md", reviewArtefact("ep-system-design-review", "EP-MALFORMED-STAGE10", "pass", "proceed_to_stage_8", 0, 0, 0, 0)},
				{"ep-implementation-plan.md", frontMatter("ep-implementation-plan", "EP-MALFORMED-STAGE10")},
				{"ep-audit-report.md", frontMatter("ep-audit-report", "EP-MALFORMED-STAGE10") + "## Summary\n\n## Implementation vs plan\n"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			epicDir := filepath.Join(dir, tt.epic)
			os.MkdirAll(epicDir, 0o755)

			for _, f := range tt.files {
				os.WriteFile(filepath.Join(epicDir, f.name), []byte(f.content), 0o644)
			}
			os.WriteFile(filepath.Join(epicDir, tt.reviewFile), []byte(tt.reviewContent), 0o644)

			result := checkPipelineState(epicDir, tt.epic)
			if !result.HasGaps {
				t.Fatalf("expected gaps for malformed readable review artefact, got findings=%v", result.Findings)
			}

			foundParseError := false
			foundGateViolation := false
			for _, finding := range result.Findings {
				if contains(finding, "unclosed YAML front matter") {
					foundParseError = true
				}
				if contains(finding, tt.gateSnippet) {
					foundGateViolation = true
				}
			}
			if !foundParseError || !foundGateViolation {
				t.Fatalf("expected both parse error and gate violation, got %v", result.Findings)
			}
		})
	}
}

func TestCheckPipelineState_UnreadableReviewArtefactDoesNotAddRedundantGateFinding(t *testing.T) {
	tests := []struct {
		name             string
		epic             string
		reviewFile       string
		reviewStageLabel string
		gateSnippet      string
		files            []struct{ name, content string }
	}{
		{
			name:             "stage 7 review",
			epic:             "EP-UNREADABLE-STAGE7",
			reviewFile:       "ep-system-design-review.md",
			reviewStageLabel: "stage 7 (ep-system-design-review.md)",
			gateSnippet:      "stage 8 (ep-implementation-plan.md) exists but stage 7 gate=",
			files: []struct{ name, content string }{
				{"ep-scope.md", frontMatter("ep-scope", "EP-UNREADABLE-STAGE7")},
				{"ep-requirements.md", frontMatter("ep-requirements", "EP-UNREADABLE-STAGE7")},
				{"ep-acceptance-criteria.md", frontMatter("ep-acceptance-criteria", "EP-UNREADABLE-STAGE7")},
				{"ep-system-design.md", frontMatter("ep-system-design", "EP-UNREADABLE-STAGE7")},
				{"ep-implementation-plan.md", frontMatter("ep-implementation-plan", "EP-UNREADABLE-STAGE7")},
			},
		},
		{
			name:             "stage 10 review",
			epic:             "EP-UNREADABLE-STAGE10",
			reviewFile:       "ep-code-review.md",
			reviewStageLabel: "stage 10 (ep-code-review.md)",
			gateSnippet:      "stage 11 (ep-audit-report.md) exists but stage 10 gate=",
			files: []struct{ name, content string }{
				{"ep-scope.md", frontMatter("ep-scope", "EP-UNREADABLE-STAGE10")},
				{"ep-requirements.md", frontMatter("ep-requirements", "EP-UNREADABLE-STAGE10")},
				{"ep-acceptance-criteria.md", frontMatter("ep-acceptance-criteria", "EP-UNREADABLE-STAGE10")},
				{"ep-system-design.md", frontMatter("ep-system-design", "EP-UNREADABLE-STAGE10")},
				{"ep-system-design-review.md", reviewArtefact("ep-system-design-review", "EP-UNREADABLE-STAGE10", "pass", "proceed_to_stage_8", 0, 0, 0, 0)},
				{"ep-implementation-plan.md", frontMatter("ep-implementation-plan", "EP-UNREADABLE-STAGE10")},
				{"ep-audit-report.md", frontMatter("ep-audit-report", "EP-UNREADABLE-STAGE10") + "## Summary\n\n## Implementation vs plan\n"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			epicDir := filepath.Join(dir, tt.epic)
			os.MkdirAll(epicDir, 0o755)

			for _, f := range tt.files {
				os.WriteFile(filepath.Join(epicDir, f.name), []byte(f.content), 0o644)
			}
			if err := os.Mkdir(filepath.Join(epicDir, tt.reviewFile), 0o755); err != nil {
				t.Fatalf("mkdir review path: %v", err)
			}

			result := checkPipelineState(epicDir, tt.epic)
			if !result.HasGaps {
				t.Fatalf("expected gaps for unreadable review artefact, got findings=%v", result.Findings)
			}

			foundReadError := false
			foundGateViolation := false
			for _, finding := range result.Findings {
				if contains(finding, tt.reviewStageLabel) && contains(finding, "is a directory") {
					foundReadError = true
				}
				if contains(finding, tt.gateSnippet) {
					foundGateViolation = true
				}
			}
			if !foundReadError {
				t.Fatalf("expected contextual read error, got %v", result.Findings)
			}
			if foundGateViolation {
				t.Fatalf("did not expect redundant gate violation for unreadable review artefact, got %v", result.Findings)
			}
		})
	}
}

func TestCheckPipelineState_FutureMissingStageRemainsMissingWithoutIOError(t *testing.T) {
	dir := t.TempDir()
	epicDir := filepath.Join(dir, "EP-FUTURE-MISSING")
	os.MkdirAll(epicDir, 0o755)

	for _, f := range []struct{ name, content string }{
		{"ep-scope.md", frontMatter("ep-scope", "EP-FUTURE-MISSING")},
		{"ep-requirements.md", frontMatter("ep-requirements", "EP-FUTURE-MISSING")},
		{"ep-acceptance-criteria.md", frontMatter("ep-acceptance-criteria", "EP-FUTURE-MISSING")},
		{"ep-system-design.md", frontMatter("ep-system-design", "EP-FUTURE-MISSING")},
		{"ep-system-design-review.md", reviewArtefact("ep-system-design-review", "EP-FUTURE-MISSING", "pass", "proceed_to_stage_8", 0, 0, 0, 0)},
		{"ep-implementation-plan.md", frontMatter("ep-implementation-plan", "EP-FUTURE-MISSING") + "## Tasks\n\n- [x] Done task\n"},
	} {
		os.WriteFile(filepath.Join(epicDir, f.name), []byte(f.content), 0o644)
	}

	result := checkPipelineState(epicDir, "EP-FUTURE-MISSING")
	if result.HasGaps {
		t.Fatalf("expected no gaps when only future optional stages are absent, got findings=%v", result.Findings)
	}

	stage10 := stageByNumber(t, result, 10)
	if stage10.Status != "missing" {
		t.Fatalf("stage 10 status = %q, want missing", stage10.Status)
	}
	for _, finding := range result.Findings {
		if contains(finding, "ep-code-review.md") && contains(finding, "stage 10") {
			t.Fatalf("unexpected direct I/O finding for absent future stage: %q", finding)
		}
	}
}

func TestCheckPipelineState_ReadFileErrorMarksStageErrorAndGap(t *testing.T) {
	dir := t.TempDir()
	epicDir := filepath.Join(dir, "EP-READFILE-ERROR")
	os.MkdirAll(epicDir, 0o755)

	for _, f := range []struct{ name, content string }{
		{"ep-scope.md", frontMatter("ep-scope", "EP-READFILE-ERROR")},
		{"ep-requirements.md", frontMatter("ep-requirements", "EP-READFILE-ERROR")},
		{"ep-acceptance-criteria.md", frontMatter("ep-acceptance-criteria", "EP-READFILE-ERROR")},
		{"ep-system-design.md", frontMatter("ep-system-design", "EP-READFILE-ERROR")},
		{"ep-system-design-review.md", reviewArtefact("ep-system-design-review", "EP-READFILE-ERROR", "pass", "proceed_to_stage_8", 0, 0, 0, 0)},
	} {
		os.WriteFile(filepath.Join(epicDir, f.name), []byte(f.content), 0o644)
	}
	if err := os.Mkdir(filepath.Join(epicDir, "ep-implementation-plan.md"), 0o755); err != nil {
		t.Fatalf("mkdir plan path: %v", err)
	}

	stderr := captureErrLog(t)
	result := checkPipelineState(epicDir, "EP-READFILE-ERROR")
	if !result.HasGaps {
		t.Fatal("expected gaps for unreadable implementation plan")
	}

	stage8 := stageByNumber(t, result, 8)
	if stage8.Status != "error" {
		t.Fatalf("stage 8 status = %q, want error", stage8.Status)
	}
	if !stage8.HasGaps {
		t.Fatal("expected stage 8 HasGaps for unreadable implementation plan")
	}

	found := false
	for _, finding := range result.Findings {
		if contains(finding, "stage 8 (ep-implementation-plan.md)") && contains(finding, "is a directory") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected contextual read error finding, got %v", result.Findings)
	}
	logged := stderr.String()
	if !contains(logged, "stage 8 (ep-implementation-plan.md)") || !contains(logged, "cannot read artefact") {
		t.Fatalf("expected stderr log for unreadable artefact, got %q", logged)
	}
}

func TestCheckPipelineState_MalformedPlanFrontMatterStillChecksUncheckedTasksForAudit(t *testing.T) {
	dir := t.TempDir()
	epicDir := filepath.Join(dir, "EP-MALFORMED-PLAN")
	os.MkdirAll(epicDir, 0o755)

	plan := "---\nartefact: ep-implementation-plan\nepic_id: EP-MALFORMED-PLAN\nstatus: approved\n## Tasks\n\n- [ ] Pending task\n"
	for _, f := range []struct{ name, content string }{
		{"ep-scope.md", frontMatter("ep-scope", "EP-MALFORMED-PLAN")},
		{"ep-requirements.md", frontMatter("ep-requirements", "EP-MALFORMED-PLAN")},
		{"ep-acceptance-criteria.md", frontMatter("ep-acceptance-criteria", "EP-MALFORMED-PLAN")},
		{"ep-system-design.md", frontMatter("ep-system-design", "EP-MALFORMED-PLAN")},
		{"ep-system-design-review.md", reviewArtefact("ep-system-design-review", "EP-MALFORMED-PLAN", "pass", "proceed_to_stage_8", 0, 0, 0, 0)},
		{"ep-implementation-plan.md", plan},
		{"ep-code-review.md", reviewArtefact("ep-code-review", "EP-MALFORMED-PLAN", "pass", "proceed_to_stage_11", 0, 0, 0, 0)},
		{"ep-audit-report.md", frontMatter("ep-audit-report", "EP-MALFORMED-PLAN") + "## Summary\n\n## Implementation vs plan\n"},
	} {
		os.WriteFile(filepath.Join(epicDir, f.name), []byte(f.content), 0o644)
	}

	result := checkPipelineState(epicDir, "EP-MALFORMED-PLAN")
	if !result.HasGaps {
		t.Fatal("expected gaps for malformed plan front matter and unchecked plan tasks")
	}

	foundFrontMatter := false
	foundUnchecked := false
	for _, finding := range result.Findings {
		if contains(finding, "stage 8 (ep-implementation-plan.md): unclosed YAML front matter") {
			foundFrontMatter = true
		}
		if contains(finding, "unchecked task") {
			foundUnchecked = true
		}
	}
	if !foundFrontMatter || !foundUnchecked {
		t.Fatalf("expected both front matter and unchecked task findings, got %v", result.Findings)
	}
}

func TestCheckPipelineState_AuditHardFailsWhenPlanIsUnreadable(t *testing.T) {
	dir := t.TempDir()
	epicDir := filepath.Join(dir, "EP-PLAN-UNREADABLE")
	os.MkdirAll(epicDir, 0o755)

	for _, f := range []struct{ name, content string }{
		{"ep-scope.md", frontMatter("ep-scope", "EP-PLAN-UNREADABLE")},
		{"ep-requirements.md", frontMatter("ep-requirements", "EP-PLAN-UNREADABLE")},
		{"ep-acceptance-criteria.md", frontMatter("ep-acceptance-criteria", "EP-PLAN-UNREADABLE")},
		{"ep-system-design.md", frontMatter("ep-system-design", "EP-PLAN-UNREADABLE")},
		{"ep-system-design-review.md", reviewArtefact("ep-system-design-review", "EP-PLAN-UNREADABLE", "pass", "proceed_to_stage_8", 0, 0, 0, 0)},
		{"ep-code-review.md", reviewArtefact("ep-code-review", "EP-PLAN-UNREADABLE", "pass", "proceed_to_stage_11", 0, 0, 0, 0)},
		{"ep-audit-report.md", frontMatter("ep-audit-report", "EP-PLAN-UNREADABLE") + "## Summary\n\n## Implementation vs plan\n"},
	} {
		os.WriteFile(filepath.Join(epicDir, f.name), []byte(f.content), 0o644)
	}
	if err := os.Mkdir(filepath.Join(epicDir, "ep-implementation-plan.md"), 0o755); err != nil {
		t.Fatalf("mkdir plan path: %v", err)
	}

	result := checkPipelineState(epicDir, "EP-PLAN-UNREADABLE")
	if !result.HasGaps {
		t.Fatal("expected audit to hard-fail when implementation plan is unreadable")
	}

	foundUnreadable := false
	for _, finding := range result.Findings {
		if contains(finding, "stage 8 (ep-implementation-plan.md)") && contains(finding, "is a directory") {
			foundUnreadable = true
		}
	}
	if !foundUnreadable {
		t.Fatalf("expected unreadable-plan finding, got %v", result.Findings)
	}
}

func TestCheckPipelineState_EmptyReadablePlanDoesNotTriggerUnreadableFallback(t *testing.T) {
	dir := t.TempDir()
	epicDir := filepath.Join(dir, "EP-EMPTY-PLAN")
	os.MkdirAll(epicDir, 0o755)

	for _, f := range []struct{ name, content string }{
		{"ep-scope.md", frontMatter("ep-scope", "EP-EMPTY-PLAN")},
		{"ep-requirements.md", frontMatter("ep-requirements", "EP-EMPTY-PLAN")},
		{"ep-acceptance-criteria.md", frontMatter("ep-acceptance-criteria", "EP-EMPTY-PLAN")},
		{"ep-system-design.md", frontMatter("ep-system-design", "EP-EMPTY-PLAN")},
		{"ep-system-design-review.md", reviewArtefact("ep-system-design-review", "EP-EMPTY-PLAN", "pass", "proceed_to_stage_8", 0, 0, 0, 0)},
		{"ep-implementation-plan.md", ""},
		{"ep-code-review.md", reviewArtefact("ep-code-review", "EP-EMPTY-PLAN", "pass", "proceed_to_stage_11", 0, 0, 0, 0)},
		{"ep-audit-report.md", frontMatter("ep-audit-report", "EP-EMPTY-PLAN") + "## Summary\n\n## Implementation vs plan\n"},
	} {
		os.WriteFile(filepath.Join(epicDir, f.name), []byte(f.content), 0o644)
	}

	result := checkPipelineState(epicDir, "EP-EMPTY-PLAN")
	if !result.HasGaps {
		t.Fatal("expected gaps for empty plan front matter/order violations")
	}

	foundFrontMatter := false
	foundFallback := false
	for _, finding := range result.Findings {
		if contains(finding, "stage 8 (ep-implementation-plan.md): no YAML front matter found") {
			foundFrontMatter = true
		}
		if contains(finding, "must be readable before checkbox evaluation") {
			foundFallback = true
		}
	}
	if !foundFrontMatter {
		t.Fatalf("expected empty plan to keep existing front matter error, got %v", result.Findings)
	}
	if foundFallback {
		t.Fatalf("did not expect unreadable fallback finding for empty but readable plan, got %v", result.Findings)
	}
}

func TestCountUncheckedPlanTasks(t *testing.T) {
	plan := "---\nartefact: ep-implementation-plan\n---\n## Tasks\n\n- [x] Done task\n- [ ] Pending task\n"
	if got := countUncheckedPlanTasks(plan); got != 1 {
		t.Errorf("countUncheckedPlanTasks = %d, want 1", got)
	}
}

func TestCheckPipelineState_UncheckedTasksWarnWhenCodeReviewExists(t *testing.T) {
	dir := t.TempDir()
	epicDir := filepath.Join(dir, "EP-PLAN-WARN")
	os.MkdirAll(epicDir, 0o755)

	plan := frontMatter("ep-implementation-plan", "EP-PLAN-WARN") + "## Tasks\n\n- [ ] Pending task\n"
	for _, f := range []struct{ name, content string }{
		{"ep-scope.md", frontMatter("ep-scope", "EP-PLAN-WARN")},
		{"ep-requirements.md", frontMatter("ep-requirements", "EP-PLAN-WARN")},
		{"ep-acceptance-criteria.md", frontMatter("ep-acceptance-criteria", "EP-PLAN-WARN")},
		{"ep-system-design.md", frontMatter("ep-system-design", "EP-PLAN-WARN")},
		{"ep-system-design-review.md", reviewArtefact("ep-system-design-review", "EP-PLAN-WARN", "pass", "proceed_to_stage_8", 0, 0, 0, 0)},
		{"ep-implementation-plan.md", plan},
		{"ep-code-review.md", reviewArtefact("ep-code-review", "EP-PLAN-WARN", "pass", "proceed_to_stage_11", 0, 0, 0, 0)},
	} {
		os.WriteFile(filepath.Join(epicDir, f.name), []byte(f.content), 0o644)
	}

	result := checkPipelineState(epicDir, "EP-PLAN-WARN")
	if result.Warnings == 0 {
		t.Fatal("expected warning for unchecked tasks when stage 10 exists")
	}
	found := false
	for _, f := range result.Findings {
		if contains(f, "unchecked task") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected unchecked task finding, got %v", result.Findings)
	}
}

func TestCheckPipelineState_UncheckedTasksErrorWhenAuditExists(t *testing.T) {
	dir := t.TempDir()
	epicDir := filepath.Join(dir, "EP-PLAN-ERR")
	os.MkdirAll(epicDir, 0o755)

	plan := frontMatter("ep-implementation-plan", "EP-PLAN-ERR") + "## Tasks\n\n- [ ] Pending task\n"
	for _, f := range []struct{ name, content string }{
		{"ep-scope.md", frontMatter("ep-scope", "EP-PLAN-ERR")},
		{"ep-requirements.md", frontMatter("ep-requirements", "EP-PLAN-ERR")},
		{"ep-acceptance-criteria.md", frontMatter("ep-acceptance-criteria", "EP-PLAN-ERR")},
		{"ep-system-design.md", frontMatter("ep-system-design", "EP-PLAN-ERR")},
		{"ep-system-design-review.md", reviewArtefact("ep-system-design-review", "EP-PLAN-ERR", "pass", "proceed_to_stage_8", 0, 0, 0, 0)},
		{"ep-implementation-plan.md", plan},
		{"ep-code-review.md", reviewArtefact("ep-code-review", "EP-PLAN-ERR", "pass", "proceed_to_stage_11", 0, 0, 0, 0)},
		{"ep-audit-report.md", frontMatter("ep-audit-report", "EP-PLAN-ERR") + "## Summary\n\n## Implementation vs plan\n"},
	} {
		os.WriteFile(filepath.Join(epicDir, f.name), []byte(f.content), 0o644)
	}

	result := checkPipelineState(epicDir, "EP-PLAN-ERR")
	if !result.HasGaps {
		t.Error("expected gaps when audit exists with unchecked plan tasks")
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

func stageByNumber(t *testing.T, result *PipelineResult, stage int) StageStatus {
	t.Helper()
	for _, ss := range result.Stages {
		if ss.Stage == stage {
			return ss
		}
	}
	t.Fatalf("stage %d not found in result", stage)
	return StageStatus{}
}
