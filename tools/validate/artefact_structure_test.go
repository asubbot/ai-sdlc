package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFrontMatter(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantErr   bool
		wantField string
		wantValue string
	}{
		{
			name:      "valid",
			content:   "---\nartefact: ep-scope\nepic_id: EP-099\nstatus: approved\nupdated_at: 2026-05-20\nsource_of_truth: true\n---\n# Title",
			wantErr:   false,
			wantField: "EpicID",
			wantValue: "EP-099",
		},
		{
			name:    "no front matter",
			content: "# Just a heading\nSome text",
			wantErr: true,
		},
		{
			name:    "unclosed",
			content: "---\nartefact: ep-scope\n# Missing close",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, err := parseFrontMatter(tt.content)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			switch tt.wantField {
			case "EpicID":
				if fm.EpicID != tt.wantValue {
					t.Errorf("EpicID = %q, want %q", fm.EpicID, tt.wantValue)
				}
			case "Artefact":
				if fm.Artefact != tt.wantValue {
					t.Errorf("Artefact = %q, want %q", fm.Artefact, tt.wantValue)
				}
			}
		})
	}
}

func TestParseFrontMatter_AllFields(t *testing.T) {
	content := "---\nartefact: ep-scope\nepic_id: EP-099\nstatus: approved\nsource_of_truth: true\nupdated_at: 2026-05-20\ngate: pass\n---\n# Body"
	fm, err := parseFrontMatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fm.Artefact != "ep-scope" {
		t.Errorf("Artefact = %q, want ep-scope", fm.Artefact)
	}
	if fm.EpicID != "EP-099" {
		t.Errorf("EpicID = %q, want EP-099", fm.EpicID)
	}
	if fm.Status != "approved" {
		t.Errorf("Status = %q, want approved", fm.Status)
	}
	if !fm.SourceOfTruth {
		t.Error("SourceOfTruth = false, want true")
	}
	if fm.UpdatedAt != "2026-05-20" {
		t.Errorf("UpdatedAt = %q, want 2026-05-20", fm.UpdatedAt)
	}
	if fm.Gate != "pass" {
		t.Errorf("Gate = %q, want pass", fm.Gate)
	}
}

func TestValidateFrontMatter(t *testing.T) {
	t.Run("all valid", func(t *testing.T) {
		fm := &FrontMatter{
			Artefact:      "ep-scope",
			EpicID:        "EP-099",
			Status:        "approved",
			SourceOfTruth: true,
			UpdatedAt:     "2026-05-20",
		}
		findings := validateFrontMatter(fm, "EP-099", "ep-scope.md")
		for _, f := range findings {
			if f.Severity == "error" {
				t.Errorf("unexpected error finding: %s", f.Message)
			}
		}
	})

	t.Run("missing fields", func(t *testing.T) {
		fm := &FrontMatter{}
		findings := validateFrontMatter(fm, "EP-099", "ep-scope.md")
		errorCount := 0
		for _, f := range findings {
			if f.Severity == "error" {
				errorCount++
			}
		}
		if errorCount < 4 {
			t.Errorf("expected at least 4 errors for empty front matter, got %d", errorCount)
		}
	})

	t.Run("wrong epic_id", func(t *testing.T) {
		fm := &FrontMatter{
			Artefact:      "ep-scope",
			EpicID:        "EP-098",
			Status:        "approved",
			SourceOfTruth: true,
			UpdatedAt:     "2026-05-20",
		}
		findings := validateFrontMatter(fm, "EP-099", "ep-scope.md")
		found := false
		for _, f := range findings {
			if f.Check == "front_matter" && f.Severity == "error" &&
				contains(f.Message, "does not match") {
				found = true
			}
		}
		if !found {
			t.Error("expected error about epic_id mismatch")
		}
	})

	t.Run("bad date format", func(t *testing.T) {
		fm := &FrontMatter{
			Artefact:      "ep-scope",
			EpicID:        "EP-099",
			Status:        "approved",
			SourceOfTruth: true,
			UpdatedAt:     "May 26, 2026",
		}
		findings := validateFrontMatter(fm, "EP-099", "ep-scope.md")
		found := false
		for _, f := range findings {
			if f.Check == "front_matter" && contains(f.Message, "YYYY-MM-DD") {
				found = true
			}
		}
		if !found {
			t.Error("expected error about date format")
		}
	})

	t.Run("source_of_truth not set", func(t *testing.T) {
		fm := &FrontMatter{
			Artefact:      "ep-scope",
			EpicID:        "EP-099",
			Status:        "approved",
			SourceOfTruth: false,
			UpdatedAt:     "2026-05-20",
		}
		findings := validateFrontMatter(fm, "EP-099", "ep-scope.md")
		found := false
		for _, f := range findings {
			if f.Severity == "warning" && contains(f.Message, "source_of_truth") {
				found = true
			}
		}
		if !found {
			t.Error("expected warning about source_of_truth")
		}
	})

	t.Run("ep-context allows source_of_truth false", func(t *testing.T) {
		fm := &FrontMatter{
			Artefact:      "ep-context",
			EpicID:        "EP-099",
			Status:        "draft",
			SourceOfTruth: false,
			UpdatedAt:     "2026-05-20",
		}
		findings := validateFrontMatter(fm, "EP-099", "ep-context.md")
		for _, f := range findings {
			if f.Severity == "warning" && contains(f.Message, "source_of_truth") {
				t.Errorf("ep-context should not warn on source_of_truth=false, got %q", f.Message)
			}
		}
	})
}

func TestFindRequiredSections(t *testing.T) {
	t.Run("all sections present for ep-scope", func(t *testing.T) {
		content := "---\nartefact: ep-scope\n---\n## Glossary\n\n## Scope\n\n## Success criteria\n\n## Traceability\n"
		findings := findRequiredSections(content, "ep-scope")
		if len(findings) != 0 {
			t.Errorf("expected no findings, got %d: %+v", len(findings), findings)
		}
	})

	t.Run("missing section", func(t *testing.T) {
		content := "---\nartefact: ep-scope\n---\n## Glossary\n\n## Scope\n\n## Traceability\n"
		findings := findRequiredSections(content, "ep-scope")
		if len(findings) != 1 {
			t.Fatalf("expected 1 finding, got %d", len(findings))
		}
		if !contains(findings[0].Message, "Success criteria") {
			t.Errorf("expected mention of Success criteria, got %q", findings[0].Message)
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		content := "## glossary\n\n## scope\n\n## success criteria\n\n## traceability\n"
		findings := findRequiredSections(content, "ep-scope")
		if len(findings) != 0 {
			t.Errorf("expected case-insensitive match, got %d findings", len(findings))
		}
	})

	t.Run("h3 headings accepted", func(t *testing.T) {
		content := "### Glossary\n\n### Scope\n\n### Success criteria\n\n### Traceability\n"
		findings := findRequiredSections(content, "ep-scope")
		if len(findings) != 0 {
			t.Errorf("expected h3 headings to be accepted, got %d findings", len(findings))
		}
	})

	t.Run("unknown artefact type", func(t *testing.T) {
		findings := findRequiredSections("anything", "ep-unknown")
		if findings != nil {
			t.Errorf("expected nil for unknown artefact type, got %+v", findings)
		}
	})

	t.Run("ep-requirements all present", func(t *testing.T) {
		content := "## Introduction\n\n## Glossary\n\n## Requirements\n"
		findings := findRequiredSections(content, "ep-requirements")
		if len(findings) != 0 {
			t.Errorf("expected no findings, got %d", len(findings))
		}
	})

	t.Run("ep-acceptance-criteria scenarios", func(t *testing.T) {
		content := "## Scenarios\n\n### AC-99.001\n"
		findings := findRequiredSections(content, "ep-acceptance-criteria")
		if len(findings) != 0 {
			t.Errorf("expected no findings, got %d", len(findings))
		}
	})

	t.Run("ep-context all required sections present", func(t *testing.T) {
		content := "## Purpose\n\n## Current Scope\n\n## Open Questions\n\n## Links\n"
		findings := findRequiredSections(content, "ep-context")
		if len(findings) != 0 {
			t.Errorf("expected no findings, got %d: %+v", len(findings), findings)
		}
	})

	t.Run("ep-context missing Open Questions", func(t *testing.T) {
		content := "## Purpose\n\n## Current Scope\n\n## Links\n"
		findings := findRequiredSections(content, "ep-context")
		if len(findings) != 1 {
			t.Fatalf("expected 1 finding, got %d", len(findings))
		}
		if !contains(findings[0].Message, "Open Questions") {
			t.Errorf("expected mention of Open Questions, got %q", findings[0].Message)
		}
	})
}

func TestFindBrokenLinks(t *testing.T) {
	dir := t.TempDir()

	existingFile := filepath.Join(dir, "existing.md")
	if err := os.WriteFile(existingFile, []byte("# Existing"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Run("valid link", func(t *testing.T) {
		content := "[link](existing.md)"
		findings := findBrokenLinks(content, dir)
		if len(findings) != 0 {
			t.Errorf("expected no broken links, got %d", len(findings))
		}
	})

	t.Run("broken link", func(t *testing.T) {
		content := "[link](nonexistent.md)"
		findings := findBrokenLinks(content, dir)
		if len(findings) != 1 {
			t.Fatalf("expected 1 broken link, got %d", len(findings))
		}
		if !contains(findings[0].Message, "nonexistent.md") {
			t.Errorf("expected message about nonexistent.md, got %q", findings[0].Message)
		}
	})

	t.Run("URL skipped", func(t *testing.T) {
		content := "[link](https://example.com)"
		findings := findBrokenLinks(content, dir)
		if len(findings) != 0 {
			t.Errorf("expected URLs to be skipped, got %d findings", len(findings))
		}
	})

	t.Run("anchor skipped", func(t *testing.T) {
		content := "[link](#section)"
		findings := findBrokenLinks(content, dir)
		if len(findings) != 0 {
			t.Errorf("expected anchors to be skipped, got %d findings", len(findings))
		}
	})

	t.Run("file with anchor valid", func(t *testing.T) {
		content := "[link](existing.md#section)"
		findings := findBrokenLinks(content, dir)
		if len(findings) != 0 {
			t.Errorf("expected file+anchor to be valid, got %d findings", len(findings))
		}
	})

	t.Run("file with anchor broken", func(t *testing.T) {
		content := "[link](missing.md#section)"
		findings := findBrokenLinks(content, dir)
		if len(findings) != 1 {
			t.Fatalf("expected 1 broken link, got %d", len(findings))
		}
	})
}

func TestValidateArtefactStructure_HealthyEpic(t *testing.T) {
	result := validateArtefactStructure("testdata/EP-099", "EP-099")

	if result.Errors != 0 {
		t.Errorf("expected 0 errors for healthy epic, got %d", result.Errors)
		for _, f := range result.Findings {
			if f.Severity == "error" {
				t.Logf("  error: [%s] %s — %s", f.File, f.Check, f.Message)
			}
		}
	}
	if result.HasGaps {
		t.Error("expected HasGaps=false for healthy epic")
	}
}

func TestValidateArtefactStructure_BrokenEpic(t *testing.T) {
	result := validateArtefactStructure("testdata/EP-098", "EP-098")

	if result.Errors == 0 {
		t.Error("expected errors for broken epic (ep-scope.md has no front matter)")
	}
	if !result.HasGaps {
		t.Error("expected HasGaps=true for broken epic")
	}

	foundScopeFMError := false
	for _, f := range result.Findings {
		if f.File == "ep-scope.md" && f.Check == "front_matter" && f.Severity == "error" {
			foundScopeFMError = true
		}
	}
	if !foundScopeFMError {
		t.Error("expected front_matter error for ep-scope.md")
	}
}

func TestDateFormat(t *testing.T) {
	tests := []struct {
		date  string
		valid bool
	}{
		{"2026-05-26", true},
		{"May 26, 2026", false},
		{"2026/05/26", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.date, func(t *testing.T) {
			got := isValidDateFormat(tt.date)
			if got != tt.valid {
				t.Errorf("isValidDateFormat(%q) = %v, want %v", tt.date, got, tt.valid)
			}
		})
	}
}

func TestValidateArtefactStructure_WithSetupHelper(t *testing.T) {
	dir := t.TempDir()
	epicDir := setupTestEpic(t, dir, "EP-100", map[string]string{
		"ep-scope.md": "---\nartefact: ep-scope\nepic_id: EP-100\nstatus: draft\nsource_of_truth: true\nupdated_at: 2026-01-01\n---\n# Scope\n\n## Glossary\n\n## Scope\n\n## Success criteria\n\n## Traceability\n",
	})

	result := validateArtefactStructure(epicDir, "EP-100")
	foundScope := false
	for _, f := range result.Findings {
		if f.File == "ep-scope.md" && f.Severity == "error" {
			t.Errorf("unexpected error on ep-scope.md: %s", f.Message)
		}
		if f.File == "ep-scope.md" {
			foundScope = true
		}
	}
	_ = foundScope

	missingCount := 0
	for _, f := range result.Findings {
		if f.Check == "file_exists" {
			missingCount++
		}
	}
	if missingCount == 0 {
		t.Error("expected warnings for missing artefact files")
	}
}

func TestValidateArtefactStructure_FailsOnUnreadableArtefact(t *testing.T) {
	dir := t.TempDir()
	epicDir := filepath.Join(dir, "EP-101")
	if err := os.MkdirAll(epicDir, 0o755); err != nil {
		t.Fatalf("mkdir epic: %v", err)
	}
	// Directory at an artefact path: Stat succeeds for some checks; ReadFile fails (EISDIR).
	if err := os.Mkdir(filepath.Join(epicDir, "ep-scope.md"), 0o755); err != nil {
		t.Fatalf("mkdir ep-scope.md path: %v", err)
	}

	stderr := captureErrLog(t)
	result := validateArtefactStructure(epicDir, "EP-101")
	if result.Errors == 0 {
		t.Fatal("expected error for unreadable ep-scope.md")
	}
	if !result.HasGaps {
		t.Fatal("expected HasGaps=true for unreadable artefact")
	}
	found := false
	for _, f := range result.Findings {
		if f.File == "ep-scope.md" && f.Check == "file_readable" && f.Severity == "error" {
			found = true
			if !contains(f.Message, "cannot read") {
				t.Errorf("message = %q, want to contain %q", f.Message, "cannot read")
			}
		}
	}
	if !found {
		t.Fatalf("expected file_readable error for ep-scope.md, got %#v", result.Findings)
	}
	if !contains(stderr.String(), "ep-scope.md") || !contains(stderr.String(), "cannot read") {
		t.Fatalf("expected stderr log for unreadable artefact, got %q", stderr.String())
	}
}

func TestParseFrontMatter_RejectsNonIntegerSeverityCounts(t *testing.T) {
	content := "---\nartefact: ep-scope\nepic_id: EP-099\nstatus: draft\nsource_of_truth: true\nupdated_at: 2026-01-01\nopen_counts:\n  blocker: not-a-number\n  major: 1\n---\n# Body"
	_, err := parseFrontMatter(content)
	if err == nil {
		t.Fatal("expected error for non-integer open_counts.blocker")
	}
	if !contains(err.Error(), "open_counts.blocker") {
		t.Errorf("error = %q, want open_counts.blocker context", err)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstr(s, substr)
}

func findSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
