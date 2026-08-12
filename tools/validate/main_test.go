package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseACsFromFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "ac-test.md")

	content := `# EP-009 Acceptance Criteria

## Acceptance Criteria

### AC-09.001 First criterion
Trigger: something happens

### AC-09.002 Second criterion
Expected: something works

**AC-09.003** (Trace: REQ-09.003)
Some requirement
`

	if err := os.WriteFile(f, []byte(content), 0o644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	acs, excluded, err := parseACsFromFile(f)
	if err != nil {
		t.Fatalf("parseACsFromFile failed: %v", err)
	}

	if len(acs) == 0 {
		t.Fatal("Expected to find ACs, but got none")
	}

	if _, ok := acs[ACCode("AC-09.003")]; !ok {
		t.Error("Expected AC-09.003 to be found")
	}
	if len(excluded) != 0 {
		t.Errorf("Expected no excluded ACs, got %d", len(excluded))
	}
}

func TestGetEpicNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"EP-009", "09"},  // 3-digit → 2-digit
		{"EP-001", "01"},  // 3-digit → 2-digit
		{"EP-99", "99"},   // 2-digit stays 2-digit
		{"EP-100", "100"}, // 3-digit but not leading zero
		{"009", "09"},     // Just number, convert
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := getEpicNumber(tt.input)
			if result != tt.expected {
				t.Errorf("getEpicNumber(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGenerateReport(t *testing.T) {
	tests := []struct {
		name              string
		acs               map[ACCode]string
		excluded          map[ACCode]acExclusionKind
		coverage          map[ACCode][]CoverageRef
		totalACs          int
		inScopeACs        int
		automatedCovered  int
		manualOnly        int
		deferredACs       int
		obsoleteACs       int
		traceabilityRatio float64
		gapCount          int
		hasBlockingGaps   bool
	}{
		{
			name: "partial coverage",
			acs: map[ACCode]string{
				"AC-09.001": "First criterion",
				"AC-09.002": "Second criterion",
				"AC-09.003": "Third criterion",
			},
			excluded: map[ACCode]acExclusionKind{},
			coverage: map[ACCode][]CoverageRef{
				"AC-09.001": {{Ref: "tests/test.go::TestFunc1", Manual: false}},
				"AC-09.002": {{Ref: "tests/test.go::TestFunc2", Manual: false}},
			},
			totalACs:          3,
			inScopeACs:        3,
			automatedCovered:  2,
			manualOnly:        0,
			deferredACs:       0,
			obsoleteACs:       0,
			traceabilityRatio: 2.0 / 3.0,
			gapCount:          1,
			hasBlockingGaps:   true,
		},
		{
			name: "full coverage",
			acs: map[ACCode]string{
				"AC-09.001": "First criterion",
				"AC-09.002": "Second criterion",
			},
			excluded: map[ACCode]acExclusionKind{},
			coverage: map[ACCode][]CoverageRef{
				"AC-09.001": {{Ref: "tests/test.go::TestFunc1", Manual: false}},
				"AC-09.002": {{Ref: "tests/test.go::TestFunc2", Manual: false}},
			},
			totalACs:          2,
			inScopeACs:        2,
			automatedCovered:  2,
			manualOnly:        0,
			deferredACs:       0,
			obsoleteACs:       0,
			traceabilityRatio: 1.0,
			gapCount:          0,
			hasBlockingGaps:   false,
		},
		{
			name: "deferred AC",
			acs: map[ACCode]string{
				"AC-09.001": "First criterion",
				"AC-09.002": "Second criterion",
				"AC-09.003": "Third criterion",
			},
			excluded: map[ACCode]acExclusionKind{
				"AC-09.003": acExclusionDeferred,
			},
			coverage: map[ACCode][]CoverageRef{
				"AC-09.001": {{Ref: "tests/test.go::TestFunc1", Manual: false}},
				"AC-09.002": {{Ref: "tests/test.go::TestFunc2", Manual: false}},
			},
			totalACs:          3,
			inScopeACs:        2,
			automatedCovered:  2,
			manualOnly:        0,
			deferredACs:       1,
			obsoleteACs:       0,
			traceabilityRatio: 1.0,
			gapCount:          1,
			hasBlockingGaps:   false,
		},
		{
			name: "obsolete AC",
			acs: map[ACCode]string{
				"AC-09.001": "First",
				"AC-09.002": "Second",
			},
			excluded: map[ACCode]acExclusionKind{
				"AC-09.002": acExclusionObsolete,
			},
			coverage: map[ACCode][]CoverageRef{
				"AC-09.001": {{Ref: "tests/x.go::TestA", Manual: false}},
			},
			totalACs:          2,
			inScopeACs:        1,
			automatedCovered:  1,
			manualOnly:        0,
			deferredACs:       0,
			obsoleteACs:       1,
			traceabilityRatio: 1.0,
			gapCount:          1,
			hasBlockingGaps:   false,
		},
		{
			name:     "manual only",
			acs:      map[ACCode]string{"AC-09.001": "c1"},
			excluded: map[ACCode]acExclusionKind{},
			coverage: map[ACCode][]CoverageRef{
				"AC-09.001": {{Ref: "t.go::TestX", Manual: true}},
			},
			totalACs:          1,
			inScopeACs:        1,
			automatedCovered:  0,
			manualOnly:        1,
			deferredACs:       0,
			obsoleteACs:       0,
			traceabilityRatio: 1.0,
			gapCount:          0,
			hasBlockingGaps:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := generateReport("EP-009", tt.acs, tt.excluded, tt.coverage)

			if r.TotalACs != tt.totalACs {
				t.Errorf("TotalACs = %d, want %d", r.TotalACs, tt.totalACs)
			}
			if r.InScopeACs != tt.inScopeACs {
				t.Errorf("InScopeACs = %d, want %d", r.InScopeACs, tt.inScopeACs)
			}
			if r.AutomatedCoveredACs != tt.automatedCovered {
				t.Errorf("AutomatedCoveredACs = %d, want %d", r.AutomatedCoveredACs, tt.automatedCovered)
			}
			if r.ManualOnlyTracedACs != tt.manualOnly {
				t.Errorf("ManualOnlyTracedACs = %d, want %d", r.ManualOnlyTracedACs, tt.manualOnly)
			}
			if r.DeferredACs != tt.deferredACs {
				t.Errorf("DeferredACs = %d, want %d", r.DeferredACs, tt.deferredACs)
			}
			if r.ObsoleteACs != tt.obsoleteACs {
				t.Errorf("ObsoleteACs = %d, want %d", r.ObsoleteACs, tt.obsoleteACs)
			}
			if r.TraceabilityRatio != tt.traceabilityRatio {
				t.Errorf("TraceabilityRatio = %f, want %f", r.TraceabilityRatio, tt.traceabilityRatio)
			}
			if len(r.Gaps) != tt.gapCount {
				t.Errorf("Gaps = %d, want %d", len(r.Gaps), tt.gapCount)
			}
			if hasBlockingGaps(r) != tt.hasBlockingGaps {
				t.Errorf("hasBlockingGaps = %v, want %v", hasBlockingGaps(r), tt.hasBlockingGaps)
			}
		})
	}
}

func TestFindCoverageInCodebase(t *testing.T) {
	tmpDir := t.TempDir()

	testsDir := filepath.Join(tmpDir, "tests")
	if err := os.Mkdir(testsDir, 0o755); err != nil {
		t.Fatalf("Failed to create tests dir: %v", err)
	}

	testFile := filepath.Join(testsDir, "example_test.go")
	testContent := `package tests

// Covers AC-09.001: test for criterion 1
func TestFunc1(t *testing.T) {
	// test code
}

// Covers AC-09.002, AC-09.003: test for criteria 2 and 3
func TestFunc2(t *testing.T) {
	// test code
}
`

	if err := os.WriteFile(testFile, []byte(testContent), 0o644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	coverage, skipCount, err := findCoverageInCodebase(tmpDir)
	if err != nil {
		t.Fatalf("findCoverageInCodebase failed: %v", err)
	}
	if skipCount != 0 {
		t.Errorf("skipCount = %d, want 0", skipCount)
	}

	if tests, ok := coverage[ACCode("AC-09.001")]; !ok {
		t.Error("Expected AC-09.001 to be found in coverage")
	} else if len(tests) == 0 {
		t.Error("Expected AC-09.001 to have at least one test")
	} else if tests[0].Manual {
		t.Error("Expected non-manual trace for AC-09.001")
	}
}

func TestParseACsFromFile_Exclusions(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		wantACCode    ACCode
		wantExclusion acExclusionKind
	}{
		{
			name: "deferred",
			content: `# EP-009 Acceptance Criteria

**AC-09.005** (Trace: [REQ-09.005](ep-requirements.md#docker-sandbox-execution))

Given something
When something
Then something
**Status:** Deferred to operations team.
`,
			wantACCode:    "AC-09.005",
			wantExclusion: acExclusionDeferred,
		},
		{
			name: "obsolete",
			content: `# EP-009 Acceptance Criteria

| [AC-09.006](#ac-09-006) | REQ | **Obsolete:** superseded by refactor. |
`,
			wantACCode:    "AC-09.006",
			wantExclusion: acExclusionObsolete,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			f := filepath.Join(dir, "ac.md")
			if err := os.WriteFile(f, []byte(tt.content), 0o644); err != nil {
				t.Fatal(err)
			}

			acs, excluded, err := parseACsFromFile(f)
			if err != nil {
				t.Fatalf("parseACsFromFile failed: %v", err)
			}
			if len(acs) < 1 {
				t.Fatal("expected at least 1 AC")
			}
			if excluded[tt.wantACCode] != tt.wantExclusion {
				t.Fatalf("want %q exclusion for %s, got %q", tt.wantExclusion, tt.wantACCode, excluded[tt.wantACCode])
			}
		})
	}
}

func TestNormalizeEpicNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		ok       bool
	}{
		{"EP-009", "09", true},
		{"EP-9", "09", true},
		{"9", "09", true},
		{"EP-100", "100", true},
		{"bad", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := normalizeEpicNumber(tt.input)
			if tt.ok && err != nil {
				t.Fatalf("normalizeEpicNumber returned unexpected error: %v", err)
			}
			if !tt.ok && err == nil {
				t.Fatal("normalizeEpicNumber expected error, got nil")
			}
			if tt.ok && got != tt.expected {
				t.Fatalf("normalizeEpicNumber(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestLineDeclaresACCoverage(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"// Covers AC-09.001", true},
		{"// covers AC-01.029, AC-01.030", true},
		{"// Supporting AC-06.002", true},
		{"// supporting AC-06.002", true},
		{"// EP-008 AC-08.001 / REQ-08.001: default_temperature", true},
		{"// AC-30.013: tools.text_based_enabled rejection", true},
		{"// each … (AC-06.005, AC-06.010 / REQ-06.013).", true},
		{"// (AC-06.005, AC-06.006, AC-06.010 / REQ-06.013)", true},
		{`message := "this covers AC-09.001"`, false},
		{"// no keyword and no pattern", false},
		{"func foo() {}", false},
	}
	for _, tt := range tests {
		if got := lineDeclaresACCoverage(tt.line); got != tt.want {
			t.Errorf("lineDeclaresACCoverage(%q) = %v, want %v", tt.line, got, tt.want)
		}
	}
}

func TestParseACsFromFile_ThreeDigitEpic(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "ac-100.md")
	content := `# EP-100 Acceptance Criteria

### AC-100.001 First three-digit epic criterion
Trace: REQ-100.001
`
	if err := os.WriteFile(f, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	acs, _, err := parseACsFromFile(f)
	if err != nil {
		t.Fatalf("parseACsFromFile failed: %v", err)
	}
	if _, ok := acs[ACCode("AC-100.001")]; !ok {
		t.Fatalf("expected AC-100.001 to be parsed, got %v", acs)
	}
}

func TestFindCoverageInCodebase_CmdAndFormats(t *testing.T) {
	tmpDir := t.TempDir()

	for _, sub := range []string{"tests", "internal", "cmd"} {
		if err := os.Mkdir(filepath.Join(tmpDir, sub), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	cmdTest := filepath.Join(tmpDir, "cmd", "pa", "x_test.go")
	if err := os.MkdirAll(filepath.Join(tmpDir, "cmd", "pa"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cmdTest, []byte(`package pa

// EP-008 AC-08.001 / REQ-08.001: body
func TestA(t *testing.T) {}
`), 0o644); err != nil {
		t.Fatal(err)
	}

	intTest := filepath.Join(tmpDir, "internal", "pkg", "y_test.go")
	if err := os.MkdirAll(filepath.Join(tmpDir, "internal", "pkg"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(intTest, []byte(`package pkg

// TestFoo covers AC-01.029
func TestFoo(t *testing.T) {}
`), 0o644); err != nil {
		t.Fatal(err)
	}

	cov, _, err := findCoverageInCodebase(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := cov[ACCode("AC-08.001")]; !ok {
		t.Error("expected AC-08.001 from cmd/ EP-008 style comment")
	}
	if _, ok := cov[ACCode("AC-01.029")]; !ok {
		t.Error("expected AC-01.029 from lowercase covers comment")
	}
}

func TestNormalizeCriterionPreview(t *testing.T) {
	if got := normalizeCriterionPreview("| col | col |"); got != "" {
		t.Errorf("table row should be empty, got %q", got)
	}
	if got := normalizeCriterionPreview("Given a user\nWhen x"); got == "" {
		t.Error("Gherkin line should be kept")
	}
}

func TestJsonOutputRequested(t *testing.T) {
	if !jsonOutputRequested(true, []string{"EP-009"}) {
		t.Error("flag true should request JSON")
	}
	if !jsonOutputRequested(false, []string{"EP-009", "--json"}) {
		t.Error("tail --json should request JSON (Go flag stops at first non-flag)")
	}
	if jsonOutputRequested(false, []string{"EP-009"}) {
		t.Error("no --json should not request JSON")
	}
}

func TestLineDeclaresManualTrace(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"// manual Covers AC-09.001", true},
		{"// Covers manual AC-01.004", true},
		{"// MANUAL only", true},
		{"// Covers AC-09.001", false},
	}
	for _, tt := range tests {
		if got := lineDeclaresManualTrace(tt.line); got != tt.want {
			t.Errorf("lineDeclaresManualTrace(%q) = %v, want %v", tt.line, got, tt.want)
		}
	}
}

func TestFindCoverageInCodebase_ManualAndSkip(t *testing.T) {
	tmpDir := t.TempDir()
	testsDir := filepath.Join(tmpDir, "tests")
	if err := os.Mkdir(testsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := `package tests

import "testing"

// manual Covers AC-09.010
func TestManualLine(t *testing.T) {}

// Covers AC-09.011
func TestWithSkip(t *testing.T) {
	t.Skip("integration")
}
`
	f := filepath.Join(testsDir, "m_test.go")
	if err := os.WriteFile(f, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cov, skipCnt, err := findCoverageInCodebase(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	if skipCnt != 1 {
		t.Errorf("Test functions with t.Skip: got %d want 1", skipCnt)
	}
	refs10 := cov[ACCode("AC-09.010")]
	if len(refs10) != 1 || !refs10[0].Manual {
		t.Errorf("AC-09.010: want one manual ref, got %+v", refs10)
	}
	refs11 := cov[ACCode("AC-09.011")]
	if len(refs11) != 1 || !refs11[0].Manual {
		t.Errorf("AC-09.011: want manual via t.Skip, got %+v", refs11)
	}
}

func TestParseTestFuncsWithTSkip(t *testing.T) {
	src := []byte(`package p

import "testing"

func TestNoSkip(t *testing.T) {}

func TestHasSkip(t *testing.T) {
	t.Skip("x")
}
`)
	m, err := parseTestFuncsWithTSkip(src, "x_test.go")
	if err != nil {
		t.Fatal(err)
	}
	if m["TestNoSkip"] {
		t.Error("TestNoSkip should not have skip")
	}
	if !m["TestHasSkip"] {
		t.Error("TestHasSkip should have skip")
	}
}

func TestParseTestFuncsWithTSkip_IgnoresNestedFuncLiteral(t *testing.T) {
	src := []byte(`package p

import "testing"

func TestNestedFuncLiteralSkip(t *testing.T) {
	run := func() {
		t.Skip("nested")
	}
	run()
}
`)

	m, err := parseTestFuncsWithTSkip(src, "x_test.go")
	if err != nil {
		t.Fatal(err)
	}
	if m["TestNestedFuncLiteralSkip"] {
		t.Error("nested function literal t.Skip should not mark enclosing test as skipped")
	}
}

func TestTestFuncsMissingACTraceInFile_tracedAndMissing(t *testing.T) {
	okFile := `package p

import "testing"

// Covers AC-09.001
func TestOK(t *testing.T) {}
`
	got, err := testFuncsMissingACTraceInFile("internal/foo/ok_test.go", okFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("want no missing, got %v", got)
	}

	two := `package p

import "testing"

// Covers AC-09.001
func TestOK(t *testing.T) {}

func TestBad(t *testing.T) {}
`
	got, err = testFuncsMissingACTraceInFile("internal/foo/b_test.go", two)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"internal/foo/b_test.go::TestBad"}
	if len(got) != 1 || got[0] != want[0] {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestTestFuncsMissingACTraceInFile_coversLineWithoutACCode(t *testing.T) {
	src := `package p

import "testing"

// Covers integration behaviour only (no AC-EE.NNN on this line)
func TestNoAC(t *testing.T) {}
`
	got, err := testFuncsMissingACTraceInFile("c_test.go", src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0] != "c_test.go::TestNoAC" {
		t.Fatalf("got %v", got)
	}
}

func TestTestFuncsMissingACTraceInFile_skipsTestMain(t *testing.T) {
	src := `package p

import "testing"

func TestMain(m *testing.M) {}

func TestReal(t *testing.T) {}
`
	got, err := testFuncsMissingACTraceInFile("m_test.go", src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0] != "m_test.go::TestReal" {
		t.Fatalf("got %v", got)
	}
}

func TestResolveSubcommand(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		wantCmd      string
		wantEpic     string
		wantExplicit bool
	}{
		{"no args", []string{}, "ac", "all", false},
		{"epic only", []string{"EP-009"}, "ac", "EP-009", false},
		{"explicit ac", []string{"ac"}, "ac", "all", true},
		{"explicit ac with epic", []string{"ac", "EP-009"}, "ac", "EP-009", true},
		{"req subcommand", []string{"req", "EP-009"}, "req", "EP-009", true},
		{"pipeline subcommand", []string{"pipeline", "EP-009"}, "pipeline", "EP-009", true},
		{"structure subcommand", []string{"structure", "EP-009"}, "structure", "EP-009", true},
		{"ears subcommand", []string{"ears", "EP-009"}, "ears", "EP-009", true},
		{"req no epic", []string{"req"}, "req", "all", true},
		{"all keyword", []string{"all"}, "ac", "all", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, epic, explicit := resolveSubcommand(tt.args)
			if cmd != tt.wantCmd {
				t.Errorf("resolveSubcommand(%v) cmd = %q, want %q", tt.args, cmd, tt.wantCmd)
			}
			if epic != tt.wantEpic {
				t.Errorf("resolveSubcommand(%v) epic = %q, want %q", tt.args, epic, tt.wantEpic)
			}
			if explicit != tt.wantExplicit {
				t.Errorf("resolveSubcommand(%v) explicit = %v, want %v", tt.args, explicit, tt.wantExplicit)
			}
		})
	}
}

func TestFindTestsMissingACTrace_tempDir(t *testing.T) {
	tmpDir := t.TempDir()
	testsDir := filepath.Join(tmpDir, "tests")
	if err := os.Mkdir(testsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	good := filepath.Join(testsDir, "good_test.go")
	if err := os.WriteFile(good, []byte(`package tests

import "testing"

// Covers AC-09.001
func TestOK(t *testing.T) {}
`), 0o644); err != nil {
		t.Fatal(err)
	}
	badPath := filepath.Join(testsDir, "bad_test.go")
	if err := os.WriteFile(badPath, []byte(`package tests

import "testing"

func TestBad(t *testing.T) {}
`), 0o644); err != nil {
		t.Fatal(err)
	}
	refs, err := findTestsMissingACTrace(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	want := "tests/bad_test.go::TestBad"
	var found bool
	for _, r := range refs {
		if r == want {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("want ref %q in %v", want, refs)
	}
	for _, r := range refs {
		if strings.Contains(r, "TestOK") {
			t.Fatalf("did not want TestOK in missing list: %v", refs)
		}
	}
}
