package main

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestParseREQsFromFile(t *testing.T) {
	tests := []struct {
		name     string
		fixture  string
		wantREQs []string
		wantErr  bool
	}{
		{"healthy epic", "testdata/EP-099/ep-requirements.md", []string{"REQ-99.001", "REQ-99.002", "REQ-99.003", "REQ-99.004", "REQ-99.005", "REQ-99.006"}, false},
		{"broken epic", "testdata/EP-098/ep-requirements.md", []string{"REQ-98.001", "REQ-98.002", "REQ-98.003", "REQ-98.004"}, false},
		{"missing file", "testdata/nonexistent.md", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqs, err := parseREQsFromFile(tt.fixture)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			var gotCodes []string
			for code := range reqs {
				gotCodes = append(gotCodes, string(code))
			}
			sort.Strings(gotCodes)
			sort.Strings(tt.wantREQs)
			if !reflect.DeepEqual(gotCodes, tt.wantREQs) {
				t.Errorf("got %v, want %v", gotCodes, tt.wantREQs)
			}
		})
	}
}

func TestParseREQsFromFile_Summaries(t *testing.T) {
	reqs, err := parseREQsFromFile("testdata/EP-099/ep-requirements.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[REQCode]string{
		"REQ-99.001": "Fixture file loading",
		"REQ-99.002": "Golden file comparison",
		"REQ-99.003": "EARS compliance check",
		"REQ-99.004": "Structure validation",
		"REQ-99.005": "Pipeline gate checking",
		"REQ-99.006": "Broken link detection",
	}
	for code, wantSummary := range want {
		if got := reqs[code]; got != wantSummary {
			t.Errorf("reqs[%s] = %q, want %q", code, got, wantSummary)
		}
	}
}

func TestParseREQRefsFromACFile(t *testing.T) {
	refs, err := parseREQRefsFromACFile("testdata/EP-099/ep-acceptance-criteria.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Each AC should reference at least one REQ
	for ac, reqRefs := range refs {
		if len(reqRefs) == 0 {
			t.Errorf("AC %s has no REQ references", ac)
		}
	}
	// Verify specific mappings
	wantMappings := map[ACCode]REQCode{
		"AC-99.001": "REQ-99.001",
		"AC-99.002": "REQ-99.002",
		"AC-99.003": "REQ-99.003",
		"AC-99.004": "REQ-99.004",
		"AC-99.005": "REQ-99.005",
		"AC-99.006": "REQ-99.006",
	}
	for ac, wantReq := range wantMappings {
		reqRefs, ok := refs[ac]
		if !ok {
			t.Errorf("AC %s not found in refs", ac)
			continue
		}
		found := false
		for _, r := range reqRefs {
			if r == wantReq {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("AC %s: expected ref to %s, got %v", ac, wantReq, reqRefs)
		}
	}
}

func TestParseREQRefsFromACFile_MissingFile(t *testing.T) {
	_, err := parseREQRefsFromACFile("testdata/nonexistent.md")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestParseREQRefsFromACFile_DoesNotLeakAcrossSections(t *testing.T) {
	dir := t.TempDir()
	acPath := filepath.Join(dir, "ep-acceptance-criteria.md")
	content := `# ACs

### AC-99.001 — First criterion
Trace: REQ-99.001

## Notes
Do not map this to ACs: REQ-99.999
`
	if err := os.WriteFile(acPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	refs, err := parseREQRefsFromACFile(acPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := refs["AC-99.001"]
	want := []REQCode{"REQ-99.001"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("AC-99.001 refs = %v, want %v", got, want)
	}
}

func TestCheckREQACTraceability_FullCoverage(t *testing.T) {
	reqs, err := parseREQsFromFile("testdata/EP-099/ep-requirements.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	acs, _, err := parseACsFromFile("testdata/EP-099/ep-acceptance-criteria.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	acReqRefs, err := parseREQRefsFromACFile("testdata/EP-099/ep-acceptance-criteria.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := checkREQACTraceability(reqs, acReqRefs, acs)
	if result.HasGaps {
		t.Errorf("expected no gaps for healthy epic, got uncovered=%v orphan=%v acsWithoutReq=%v",
			result.UncoveredREQs, result.OrphanACRefs, result.ACsWithoutREQ)
	}
	if result.CoverageRatio != 1.0 {
		t.Errorf("coverage ratio = %f, want 1.0", result.CoverageRatio)
	}
	if result.TotalREQs != 6 {
		t.Errorf("TotalREQs = %d, want 6", result.TotalREQs)
	}
	if result.CoveredREQs != 6 {
		t.Errorf("CoveredREQs = %d, want 6", result.CoveredREQs)
	}
}

func TestCheckREQACTraceability_EmptyREQs(t *testing.T) {
	reqs := map[REQCode]string{}
	acReqRefs := map[ACCode][]REQCode{
		"AC-99.001": {"REQ-99.001"},
	}
	acs := map[ACCode]string{
		"AC-99.001": "Some criterion",
	}
	result := checkREQACTraceability(reqs, acReqRefs, acs)
	if result.TotalREQs != 0 {
		t.Errorf("TotalREQs = %d, want 0", result.TotalREQs)
	}
	if result.CoverageRatio != 0 {
		t.Errorf("CoverageRatio = %f, want 0", result.CoverageRatio)
	}
	// REQ-99.001 referenced by AC but doesn't exist → orphan
	if len(result.OrphanACRefs) != 1 || result.OrphanACRefs[0] != "REQ-99.001" {
		t.Errorf("OrphanACRefs = %v, want [REQ-99.001]", result.OrphanACRefs)
	}
}

func TestCheckREQACTraceability_OrphanREQRef(t *testing.T) {
	reqs := map[REQCode]string{
		"REQ-99.001": "Existing requirement",
	}
	acReqRefs := map[ACCode][]REQCode{
		"AC-99.001": {"REQ-99.001", "REQ-99.999"},
	}
	acs := map[ACCode]string{
		"AC-99.001": "Some criterion",
	}
	result := checkREQACTraceability(reqs, acReqRefs, acs)
	if len(result.OrphanACRefs) != 1 || result.OrphanACRefs[0] != "REQ-99.999" {
		t.Errorf("OrphanACRefs = %v, want [REQ-99.999]", result.OrphanACRefs)
	}
	if result.CoverageRatio != 1.0 {
		t.Errorf("CoverageRatio = %f, want 1.0", result.CoverageRatio)
	}
	if !result.HasGaps {
		t.Error("expected HasGaps=true due to orphan ref")
	}
}

func TestCheckREQACTraceability_UncoveredREQ(t *testing.T) {
	reqs := map[REQCode]string{
		"REQ-99.001": "Covered requirement",
		"REQ-99.002": "Uncovered requirement",
	}
	acReqRefs := map[ACCode][]REQCode{
		"AC-99.001": {"REQ-99.001"},
	}
	acs := map[ACCode]string{
		"AC-99.001": "Some criterion",
	}
	result := checkREQACTraceability(reqs, acReqRefs, acs)
	if result.CoveredREQs != 1 {
		t.Errorf("CoveredREQs = %d, want 1", result.CoveredREQs)
	}
	if len(result.UncoveredREQs) != 1 || result.UncoveredREQs[0] != "REQ-99.002" {
		t.Errorf("UncoveredREQs = %v, want [REQ-99.002]", result.UncoveredREQs)
	}
	if result.CoverageRatio != 0.5 {
		t.Errorf("CoverageRatio = %f, want 0.5", result.CoverageRatio)
	}
	if !result.HasGaps {
		t.Error("expected HasGaps=true due to uncovered REQ")
	}
}

func TestCheckREQACTraceability_ACWithoutREQ(t *testing.T) {
	reqs := map[REQCode]string{
		"REQ-99.001": "A requirement",
	}
	acReqRefs := map[ACCode][]REQCode{
		"AC-99.001": {"REQ-99.001"},
	}
	acs := map[ACCode]string{
		"AC-99.001": "Traced criterion",
		"AC-99.002": "Orphan criterion with no REQ reference",
	}
	result := checkREQACTraceability(reqs, acReqRefs, acs)
	if len(result.ACsWithoutREQ) != 1 || result.ACsWithoutREQ[0] != "AC-99.002" {
		t.Errorf("ACsWithoutREQ = %v, want [AC-99.002]", result.ACsWithoutREQ)
	}
	if !result.HasGaps {
		t.Error("expected HasGaps=true due to AC without REQ")
	}
}

func TestCheckREQACTraceability_REQToACsMapping(t *testing.T) {
	reqs := map[REQCode]string{
		"REQ-99.001": "Requirement one",
	}
	acReqRefs := map[ACCode][]REQCode{
		"AC-99.001": {"REQ-99.001"},
		"AC-99.002": {"REQ-99.001"},
	}
	acs := map[ACCode]string{
		"AC-99.001": "Criterion one",
		"AC-99.002": "Criterion two",
	}
	result := checkREQACTraceability(reqs, acReqRefs, acs)
	got := result.REQToACs["REQ-99.001"]
	want := []string{"AC-99.001", "AC-99.002"}
	sort.Strings(got)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("REQToACs[REQ-99.001] = %v, want %v", got, want)
	}
}
