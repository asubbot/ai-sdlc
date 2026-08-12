package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScanEpicsAgainstCoverageFailsOnUnreadableACFile(t *testing.T) {
	dir := t.TempDir()
	epicDir := setupTestEpic(t, dir, "EP-042", nil)
	if err := os.Mkdir(filepath.Join(epicDir, "ep-acceptance-criteria.md"), 0o755); err != nil {
		t.Fatalf("mkdir placeholder: %v", err)
	}

	_, _, _, err := scanEpicsAgainstCoverage(dir, []string{"EP-042"}, nil)

	if err == nil {
		t.Fatal("scanEpicsAgainstCoverage returned no error for an unreadable AC file")
	}
	if !strings.Contains(err.Error(), "ep-acceptance-criteria.md") ||
		!strings.Contains(err.Error(), "is a directory") {
		t.Fatalf("error does not name the unreadable AC file: %v", err)
	}
}

func TestScanEpicsAgainstCoverageFailsOnUnreadableREQFile(t *testing.T) {
	dir := t.TempDir()
	epicDir := setupTestEpic(t, dir, "EP-042", map[string]string{
		"ep-acceptance-criteria.md": "### AC-42.001 First criterion\n",
	})
	if err := os.Mkdir(filepath.Join(epicDir, "ep-requirements.md"), 0o755); err != nil {
		t.Fatalf("mkdir placeholder: %v", err)
	}

	_, _, _, err := scanEpicsAgainstCoverage(dir, []string{"EP-042"}, nil)

	if err == nil {
		t.Fatal("scanEpicsAgainstCoverage returned no error for an unreadable requirements file")
	}
	if !strings.Contains(err.Error(), "ep-requirements.md") ||
		!strings.Contains(err.Error(), "is a directory") {
		t.Fatalf("error does not name the unreadable requirements file: %v", err)
	}
}

func TestScanEpicsAgainstCoverageFailsOnUnstattableREQFile(t *testing.T) {
	dir := t.TempDir()
	epicDir := setupTestEpic(t, dir, "EP-042", map[string]string{
		"ep-acceptance-criteria.md": "### AC-42.001 First criterion\n",
	})
	reqPath := filepath.Join(epicDir, "ep-requirements.md")
	if err := os.Symlink(reqPath, reqPath); err != nil {
		t.Fatalf("symlink loop: %v", err)
	}

	_, _, _, err := scanEpicsAgainstCoverage(dir, []string{"EP-042"}, nil)

	if err == nil {
		t.Fatal("scanEpicsAgainstCoverage returned no error for an unstattable requirements file")
	}
	if !strings.Contains(err.Error(), "ep-requirements.md") {
		t.Fatalf("error does not name the requirements file: %v", err)
	}
}

func TestScanEpicsAgainstCoverageAllowsMissingREQFile(t *testing.T) {
	dir := t.TempDir()
	setupTestEpic(t, dir, "EP-042", map[string]string{
		"ep-acceptance-criteria.md": "### AC-42.001 First criterion\n",
	})

	results, _, _, err := scanEpicsAgainstCoverage(dir, []string{"EP-042"}, nil)

	if err != nil {
		t.Fatalf("scanEpicsAgainstCoverage failed for an epic without requirements: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("got %d epic summaries, want 1", len(results))
	}
	if results[0].REQCount != 0 {
		t.Fatalf("REQCount = %d, want 0 for a missing requirements file", results[0].REQCount)
	}
}

func TestScanEpicsAgainstCoverageFailsOnNonNumericEpicID(t *testing.T) {
	dir := t.TempDir()
	_, _, _, err := scanEpicsAgainstCoverage(dir, []string{"EP-foo"}, nil)
	if err == nil {
		t.Fatal("scanEpicsAgainstCoverage returned no error for non-numeric epic id")
	}
	if !strings.Contains(err.Error(), "EP-foo") {
		t.Fatalf("error does not name the epic id: %v", err)
	}
}
