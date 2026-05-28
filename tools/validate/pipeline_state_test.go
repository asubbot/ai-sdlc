package main

import (
	"os"
	"path/filepath"
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
