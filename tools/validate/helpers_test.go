package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// setupTestEpic creates an epic directory structure under dir with the given files.
// files is a map of relative paths (e.g., "ep-scope.md") to content strings.
// Returns the full path to the epic directory.
func setupTestEpic(t *testing.T, dir, epicID string, files map[string]string) string {
	t.Helper()
	epicDir := filepath.Join(dir, "ai-sdlc-artefacts", "epics", epicID)
	if err := os.MkdirAll(epicDir, 0o755); err != nil {
		t.Fatalf("setupTestEpic: mkdir %s: %v", epicDir, err)
	}
	for name, content := range files {
		p := filepath.Join(epicDir, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatalf("setupTestEpic: mkdir for %s: %v", name, err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatalf("setupTestEpic: write %s: %v", name, err)
		}
	}
	return epicDir
}

// readFixture reads a fixture file from testdata/ and returns its content.
func readFixture(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("readFixture: %v", err)
	}
	return string(data)
}

// assertJSONEqual compares got JSON string against a golden file.
// If UPDATE_GOLDEN=1 env var is set, it updates the golden file instead.
func assertJSONEqual(t *testing.T, got, goldenPath string) {
	t.Helper()

	var gotObj interface{}
	if err := json.Unmarshal([]byte(got), &gotObj); err != nil {
		t.Fatalf("assertJSONEqual: unmarshal got: %v", err)
	}
	gotNorm, err := json.MarshalIndent(gotObj, "", "  ")
	if err != nil {
		t.Fatalf("assertJSONEqual: marshal got: %v", err)
	}

	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0o755); err != nil {
			t.Fatalf("assertJSONEqual: mkdir for golden: %v", err)
		}
		if err := os.WriteFile(goldenPath, append(gotNorm, '\n'), 0o644); err != nil {
			t.Fatalf("assertJSONEqual: write golden: %v", err)
		}
		t.Logf("Updated golden file: %s", goldenPath)
		return
	}

	wantData, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("assertJSONEqual: read golden %s: %v (run with UPDATE_GOLDEN=1 to create)", goldenPath, err)
	}
	var wantObj interface{}
	if err := json.Unmarshal(wantData, &wantObj); err != nil {
		t.Fatalf("assertJSONEqual: unmarshal golden: %v", err)
	}
	wantNorm, err := json.MarshalIndent(wantObj, "", "  ")
	if err != nil {
		t.Fatalf("assertJSONEqual: marshal golden: %v", err)
	}

	if string(gotNorm) != string(wantNorm) {
		t.Errorf("JSON mismatch with golden %s\ngot:\n%s\nwant:\n%s", goldenPath, gotNorm, wantNorm)
	}
}
