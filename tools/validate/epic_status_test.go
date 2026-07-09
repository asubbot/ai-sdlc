package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseEpicStatusFromScope(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		want      string
		wantFound bool
	}{
		{
			name:      "done",
			content:   "| **Status** | DONE |\n",
			want:      "DONE",
			wantFound: true,
		},
		{
			name:      "new",
			content:   "| **Status** | NEW |\n",
			want:      "NEW",
			wantFound: true,
		},
		{
			name:      "canceled with note",
			content:   "| **Status** | CANCELED (UX is not good for the product) |\n",
			want:      "CANCELED (UX is not good for the product)",
			wantFound: true,
		},
		{
			name:      "cancel shorthand",
			content:   "| **Status** | CANCEL (Not necessary) |\n",
			want:      "CANCEL (Not necessary)",
			wantFound: true,
		},
		{
			name:      "missing",
			content:   "# Epic scope\n\nNo status table.\n",
			wantFound: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := parseEpicStatusFromScope(tt.content)
			if found != tt.wantFound {
				t.Fatalf("found=%v want %v", found, tt.wantFound)
			}
			if got != tt.want {
				t.Fatalf("status=%q want %q", got, tt.want)
			}
		})
	}
}

func TestIsSkippedEpicStatus(t *testing.T) {
	tests := []struct {
		status string
		want   bool
	}{
		{"NEW", true},
		{"CANCEL (foo)", true},
		{"CANCELED (foo)", true},
		{"DONE", false},
		{"IN_PROGRESS", false},
	}
	for _, tt := range tests {
		if got := isSkippedEpicStatus(tt.status); got != tt.want {
			t.Errorf("isSkippedEpicStatus(%q)=%v want %v", tt.status, got, tt.want)
		}
	}
}

func TestReadInScopeEpicNames(t *testing.T) {
	root := t.TempDir()
	epicsPath := filepath.Join(root, "epics")
	if err := os.MkdirAll(epicsPath, 0o755); err != nil {
		t.Fatal(err)
	}

	writeScope := func(epic, status string) {
		dir := filepath.Join(epicsPath, epic)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		body := "# scope\n\n| **Status** | " + status + " |\n"
		if err := os.WriteFile(filepath.Join(dir, "ep-scope.md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	writeScope("EP-001", "DONE")
	writeScope("EP-002", "NEW")
	writeScope("EP-003", "CANCELED (reason)")
	writeScope("EP-004", "IN_PROGRESS")

	inScope, skipped, err := readInScopeEpicNames(epicsPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(inScope) != 2 || inScope[0] != "EP-001" || inScope[1] != "EP-004" {
		t.Fatalf("inScope=%v", inScope)
	}
	if len(skipped) != 2 || skipped[0] != "EP-002" || skipped[1] != "EP-003" {
		t.Fatalf("skipped=%v", skipped)
	}
}
