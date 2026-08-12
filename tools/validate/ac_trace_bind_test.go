package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTestSource(t *testing.T, root, relPath, content string) {
	t.Helper()
	fullPath := filepath.Join(root, relPath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatalf("mkdir for %s: %v", relPath, err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", relPath, err)
	}
}

func writeLoopTestSymlink(t *testing.T, root, relPath string) {
	t.Helper()
	fullPath := filepath.Join(root, relPath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatalf("mkdir for %s: %v", relPath, err)
	}
	if err := os.Symlink(filepath.Base(fullPath), fullPath); err != nil {
		t.Skipf("self-referential symlink fixture unsupported: %v", err)
	}
}

func TestParseTestASTIndex_ExcludesMethods(t *testing.T) {
	src := []byte(`package tests

import "testing"

type suite struct{}

// Covers AC-09.008
func (suite) TestMethod(t *testing.T) {}

// Covers AC-09.009
func TestTopLevel(t *testing.T) {}
`)

	idx, err := parseTestASTIndex(src, "x_test.go")
	if err != nil {
		t.Fatalf("parseTestASTIndex returned error: %v", err)
	}

	names := idx.topLevelTestNames()
	if len(names) != 1 || names[0] != "TestTopLevel" {
		t.Fatalf("expected only top-level TestTopLevel, got %v", names)
	}

	if _, err := idx.bindTraceLine(7); err == nil {
		t.Fatal("expected method doc comment line to be rejected as orphan trace")
	}

	name, err := idx.bindTraceLine(10)
	if err != nil {
		t.Fatalf("expected top-level doc comment to bind, got error: %v", err)
	}
	if name != "TestTopLevel" {
		t.Fatalf("expected top-level doc comment to bind to TestTopLevel, got %q", name)
	}
}

func TestFindCoverageInCodebase_BindsDocCommentToFollowingTest(t *testing.T) {
	tmpDir := t.TempDir()
	writeTestSource(t, tmpDir, "tests/doc_bind_test.go", `package tests

import "testing"

// Covers AC-09.001
func TestDocBound(t *testing.T) {}

func TestOther(t *testing.T) {}
`)

	coverage, _, err := findCoverageInCodebase(tmpDir)
	if err != nil {
		t.Fatalf("findCoverageInCodebase returned error: %v", err)
	}

	refs := coverage[ACCode("AC-09.001")]
	if len(refs) != 1 {
		t.Fatalf("expected one ref for AC-09.001, got %+v", refs)
	}
	if refs[0].Ref != "tests/doc_bind_test.go::TestDocBound" {
		t.Fatalf("expected doc comment to bind to TestDocBound, got %+v", refs[0])
	}
	if refs[0].Manual {
		t.Fatalf("expected automated ref, got %+v", refs[0])
	}
}

func TestFindCoverageInCodebase_BindsInlineTraceToEnclosingTest(t *testing.T) {
	tmpDir := t.TempDir()
	writeTestSource(t, tmpDir, "tests/inline_bind_test.go", `package tests

import "testing"

func TestA(t *testing.T) {
	// Covers AC-09.002
}

func TestB(t *testing.T) {}
`)

	coverage, _, err := findCoverageInCodebase(tmpDir)
	if err != nil {
		t.Fatalf("findCoverageInCodebase returned error: %v", err)
	}

	refs := coverage[ACCode("AC-09.002")]
	if len(refs) != 1 {
		t.Fatalf("expected one ref for AC-09.002, got %+v", refs)
	}
	if refs[0].Ref != "tests/inline_bind_test.go::TestA" {
		t.Fatalf("expected inline trace to bind to TestA, got %+v", refs[0])
	}
}

func TestFindCoverageInCodebase_IgnoresRawStringPseudoCommentAtPackageLevel(t *testing.T) {
	tmpDir := t.TempDir()
	writeTestSource(t, tmpDir, "tests/raw_string_package_level_test.go", "package tests\n\nimport \"testing\"\n\nconst fixture = `\n// Covers AC-09.012\n`\n\n// Covers AC-09.013\nfunc TestRealTrace(t *testing.T) {}\n")

	coverage, skipCount, err := findCoverageInCodebase(tmpDir)
	if err != nil {
		t.Fatalf("expected package-level raw string pseudo-comment to be ignored, got error: %v", err)
	}
	if skipCount != 0 {
		t.Fatalf("expected no skipped test functions, got %d", skipCount)
	}
	if refs := coverage[ACCode("AC-09.012")]; len(refs) != 0 {
		t.Fatalf("expected no coverage from package-level raw string pseudo-comment, got %+v", refs)
	}
	refs := coverage[ACCode("AC-09.013")]
	if len(refs) != 1 || refs[0].Ref != "tests/raw_string_package_level_test.go::TestRealTrace" {
		t.Fatalf("expected real comment trace to remain valid, got %+v", refs)
	}
}

func TestFindCoverageAndTestTrace_IgnoresRawStringPseudoCommentInsideTest(t *testing.T) {
	tmpDir := t.TempDir()
	writeTestSource(t, tmpDir, "tests/raw_string_in_test_test.go", "package tests\n\nimport \"testing\"\n\nfunc TestRawStringOnly(t *testing.T) {\n\tfixture := `\n// Covers AC-09.014\n`\n\t_ = fixture\n}\n\n// Covers AC-09.015\nfunc TestRealTrace(t *testing.T) {}\n")

	coverage, skipCount, missing, err := findCoverageAndTestTrace(tmpDir)
	if err != nil {
		t.Fatalf("expected in-test raw string pseudo-comment to be ignored, got error: %v", err)
	}
	if skipCount != 0 {
		t.Fatalf("expected no skipped test functions, got %d", skipCount)
	}
	if refs := coverage[ACCode("AC-09.014")]; len(refs) != 0 {
		t.Fatalf("expected no coverage from in-test raw string pseudo-comment, got %+v", refs)
	}
	refs := coverage[ACCode("AC-09.015")]
	if len(refs) != 1 || refs[0].Ref != "tests/raw_string_in_test_test.go::TestRealTrace" {
		t.Fatalf("expected real comment trace to remain valid, got %+v", refs)
	}
	wantMissing := "tests/raw_string_in_test_test.go::TestRawStringOnly"
	if len(missing) != 1 || missing[0] != wantMissing {
		t.Fatalf("expected raw-string-only test to remain missing AC trace, got %v", missing)
	}
}

func TestFindCoverageInCodebase_IgnoresOrphanCoverageTextWithoutParseableAC(t *testing.T) {
	tmpDir := t.TempDir()
	writeTestSource(t, tmpDir, "tests/non_parseable_orphan_test.go", `package tests

import "testing"

// Covers integration behavior
const ignored = 1

// Covers AC-09.ABC
const alsoIgnored = 2
`)

	coverage, skipCount, err := findCoverageInCodebase(tmpDir)
	if err != nil {
		t.Fatalf("expected non-parseable orphan coverage text to be ignored, got error: %v", err)
	}
	if skipCount != 0 {
		t.Fatalf("expected no skipped test functions, got %d", skipCount)
	}
	if len(coverage) != 0 {
		t.Fatalf("expected no coverage from ignored orphan lines, got %+v", coverage)
	}
}

func TestFindCoverageInCodebase_OrphanTraceReturnsContextualError(t *testing.T) {
	tmpDir := t.TempDir()
	writeTestSource(t, tmpDir, "tests/orphan_test.go", `package tests

import "testing"

// Covers AC-09.003
const orphan = 1
`)

	_, _, err := findCoverageInCodebase(tmpDir)
	if err == nil {
		t.Fatal("expected orphan trace error, got nil")
	}
	if !strings.Contains(err.Error(), "tests/orphan_test.go:5") {
		t.Fatalf("expected contextual orphan error, got %v", err)
	}
}

func TestFindCoverageInCodebase_MalformedGoReturnsError(t *testing.T) {
	tmpDir := t.TempDir()
	writeTestSource(t, tmpDir, "tests/bad_syntax_test.go", `package tests

import "testing"

func TestBroken(t *testing.T) {
	// Covers AC-09.004
`)

	_, _, err := findCoverageInCodebase(tmpDir)
	if err == nil {
		t.Fatal("expected malformed Go error, got nil")
	}
	if !strings.Contains(err.Error(), "tests/bad_syntax_test.go") {
		t.Fatalf("expected filename in parse error, got %v", err)
	}
}

func TestFindCoverageInCodebase_ReadFileErrorPropagates(t *testing.T) {
	tmpDir := t.TempDir()
	writeLoopTestSymlink(t, tmpDir, "tests/loop_test.go")

	_, _, err := findCoverageInCodebase(tmpDir)
	if err == nil {
		t.Fatal("expected read-file error, got nil")
	}
	if !strings.Contains(err.Error(), "tests/loop_test.go") {
		t.Fatalf("expected unreadable path in error, got %v", err)
	}
	if !strings.Contains(err.Error(), "read test file") {
		t.Fatalf("expected root.ReadFile wrapper in error, got %v", err)
	}
}

func TestFindTestsMissingACTrace_ReadFileErrorPropagates(t *testing.T) {
	tmpDir := t.TempDir()
	writeLoopTestSymlink(t, tmpDir, "tests/loop_test.go")

	_, err := findTestsMissingACTrace(tmpDir)
	if err == nil {
		t.Fatal("expected read-file error, got nil")
	}
	if !strings.Contains(err.Error(), "tests/loop_test.go") {
		t.Fatalf("expected unreadable path in error, got %v", err)
	}
	if !strings.Contains(err.Error(), "read test file") {
		t.Fatalf("expected root.ReadFile wrapper in error, got %v", err)
	}
}

func TestFindCoverageAndTestTrace_UsesSameBindingRules(t *testing.T) {
	tmpDir := t.TempDir()
	writeTestSource(t, tmpDir, "tests/consistency_test.go", `package tests

import "testing"

// Covers AC-09.005
func TestDocBound(t *testing.T) {}

func TestInlineBound(t *testing.T) {
	// Covers AC-09.006
}

// Covers AC-09.007
func TestManualSkip(t *testing.T) {
	t.Skip("manual")
}

func TestMissing(t *testing.T) {}
`)

	coverage, skipCount, missing, err := findCoverageAndTestTrace(tmpDir)
	if err != nil {
		t.Fatalf("findCoverageAndTestTrace returned error: %v", err)
	}

	if skipCount != 1 {
		t.Fatalf("expected one skipped test function, got %d", skipCount)
	}

	refs005 := coverage[ACCode("AC-09.005")]
	if len(refs005) != 1 || refs005[0].Ref != "tests/consistency_test.go::TestDocBound" || refs005[0].Manual {
		t.Fatalf("expected AC-09.005 to bind to TestDocBound automatically, got %+v", refs005)
	}

	refs006 := coverage[ACCode("AC-09.006")]
	if len(refs006) != 1 || refs006[0].Ref != "tests/consistency_test.go::TestInlineBound" || refs006[0].Manual {
		t.Fatalf("expected AC-09.006 to bind to TestInlineBound automatically, got %+v", refs006)
	}

	refs007 := coverage[ACCode("AC-09.007")]
	if len(refs007) != 1 || refs007[0].Ref != "tests/consistency_test.go::TestManualSkip" || !refs007[0].Manual {
		t.Fatalf("expected AC-09.007 to bind to TestManualSkip manually, got %+v", refs007)
	}

	wantMissing := "tests/consistency_test.go::TestMissing"
	if len(missing) != 1 || missing[0] != wantMissing {
		t.Fatalf("expected missing trace %q, got %v", wantMissing, missing)
	}
}
