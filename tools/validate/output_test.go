package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// captureErrLog redirects errLog for the duration of the test and returns a
// buffer holding everything written to it.
func captureErrLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	prev := errLog.Writer()
	errLog.SetOutput(&buf)
	t.Cleanup(func() { errLog.SetOutput(prev) })
	return &buf
}

func TestErrLogHasNoPrefixOrTimestamp(t *testing.T) {
	buf := captureErrLog(t)

	errLog.Printf("Error: %v\n", errors.New("boom"))

	if got, want := buf.String(), "Error: boom\n"; got != want {
		t.Fatalf("errLog output = %q, want %q", got, want)
	}
}

func TestValidateSingleEpicEARSReportsUnreadableScopeOnErrLog(t *testing.T) {
	t.Chdir(t.TempDir())
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	buf := captureErrLog(t)

	if validateSingleEpicEARS("EP-999", false) {
		t.Fatal("validateSingleEpicEARS returned true for an epic without ep-scope.md")
	}

	want := "Error validating EP-999: read ep-scope.md: open " +
		filepath.Join(cwd, "ai-sdlc-artefacts", "epics", "EP-999", "ep-scope.md") +
		": no such file or directory\n"
	if got := buf.String(); got != want {
		t.Fatalf("errLog output = %q, want %q", got, want)
	}
}
