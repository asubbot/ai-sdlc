package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// findCoverageInCodebase searches for AC traceability comments in test files.
// The second return value is the number of Test* functions whose body contains t.Skip (project-wide).
func findCoverageInCodebase(codebasePath string) (map[ACCode][]CoverageRef, int, error) {
	coverage := make(map[ACCode][]CoverageRef)
	testFuncsWithSkip := 0

	root, err := os.OpenRoot(codebasePath)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = root.Close() }()

	fsys := root.FS()

	for _, dir := range []string{"tests", "internal", "cmd"} {
		n, werr := appendCoverageFromTestDir(root, fsys, dir, coverage)
		if werr != nil {
			return nil, 0, werr
		}
		testFuncsWithSkip += n
	}

	return coverage, testFuncsWithSkip, nil
}

// appendCoverageFromTestDir walks dir under root (tests, internal, or cmd) and merges AC coverage into coverage.
func appendCoverageFromTestDir(root *os.Root, fsys fs.FS, dir string, coverage map[ACCode][]CoverageRef) (int, error) {
	if _, err := root.Stat(dir); os.IsNotExist(err) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	var testFuncsWithSkip int
	err := fs.WalkDir(fsys, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk %s: %w", path, err)
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, "_test.go") {
			return nil
		}

		content, err := root.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read test file %s: %w", path, err)
		}

		index, err := parseTestASTIndex(content, path)
		if err != nil {
			return fmt.Errorf("parse test file %s: %w", path, err)
		}
		testFuncsWithSkip += index.countTestsWithSkip()

		fileContent := string(content)
		relPath := path

		lines := strings.Split(fileContent, "\n")

		for i, line := range lines {
			if !index.isActualCommentLine(i + 1) {
				continue
			}
			if !lineDeclaresACCoverage(line) {
				continue
			}
			acs := extractACsFromLine(line)
			if len(acs) == 0 {
				continue
			}
			testName, err := index.bindTraceLine(i + 1)
			if err != nil {
				return err
			}
			manual := lineDeclaresManualTrace(line) || index.hasDirectTSkip(testName)
			ref := CoverageRef{
				Ref:    fmt.Sprintf("%s::%s", relPath, testName),
				Manual: manual,
			}

			for _, ac := range acs {
				coverage[ac] = append(coverage[ac], ref)
			}
		}

		return nil
	})
	return testFuncsWithSkip, err
}

func filterCoverageForEpicNum(coverage map[ACCode][]CoverageRef, epicNum string) map[ACCode][]CoverageRef {
	out := make(map[ACCode][]CoverageRef)
	for ac, tests := range coverage {
		if strings.Contains(string(ac), "-"+epicNum+".") {
			out[ac] = tests
		}
	}
	return out
}

func uniqueTestRefsCount(coverage map[ACCode][]CoverageRef) int {
	seen := make(map[string]struct{})
	for _, refs := range coverage {
		for _, ref := range refs {
			seen[ref.Ref] = struct{}{}
		}
	}
	return len(seen)
}

// generateReport creates the AC coverage report.
func generateReport(epic string, acs map[ACCode]string, excluded map[ACCode]acExclusionKind, coverage map[ACCode][]CoverageRef) *Report {
	r := &Report{
		Epic:     epic,
		TotalACs: len(acs),
		Coverage: make(map[string][]CoverageRef),
	}

	for ac, testRefs := range coverage {
		r.Coverage[string(ac)] = append([]CoverageRef(nil), testRefs...)
	}

	for ac := range acs {
		switch excluded[ac] {
		case acExclusionObsolete:
			r.Gaps = append(r.Gaps, ACGap{
				Code:   string(ac),
				Status: "obsolete",
				Reason: "Obsolete in ep-acceptance-criteria.md",
			})
			r.ObsoleteACs++
			continue
		case acExclusionDeferred:
			r.Gaps = append(r.Gaps, ACGap{
				Code:   string(ac),
				Status: "deferred",
				Reason: "Deferred in ep-acceptance-criteria.md",
			})
			r.DeferredACs++
			continue
		}
		refs := coverage[ac]
		if len(refs) == 0 {
			r.Gaps = append(r.Gaps, ACGap{
				Code:   string(ac),
				Status: "not_covered",
			})
			continue
		}
		hasAuto := false
		for _, ref := range refs {
			if !ref.Manual {
				hasAuto = true
				break
			}
		}
		if hasAuto {
			r.AutomatedCoveredACs++
		} else {
			r.ManualOnlyTracedACs++
		}
	}

	r.InScopeACs = r.TotalACs - r.DeferredACs - r.ObsoleteACs
	if r.InScopeACs > 0 {
		traced := r.AutomatedCoveredACs + r.ManualOnlyTracedACs
		r.TraceabilityRatio = float64(traced) / float64(r.InScopeACs)
		r.AutomatedRatio = float64(r.AutomatedCoveredACs) / float64(r.InScopeACs)
	}

	sort.Slice(r.Gaps, func(i, j int) bool {
		return r.Gaps[i].Code < r.Gaps[j].Code
	})

	return r
}

func hasBlockingGaps(r *Report) bool {
	for _, gap := range r.Gaps {
		if gap.Status == "not_covered" {
			return true
		}
	}
	return false
}

// epicREQCount counts REQ ids in an epic's requirements file. A missing file is a
// valid state — not every epic has requirements yet — but any other failure to
// stat or read it must surface instead of being reported as zero requirements.
func epicREQCount(reqPath string) (int, error) {
	if _, err := os.Stat(reqPath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return 0, nil
		}
		return 0, fmt.Errorf("stat %s: %w", reqPath, err)
	}
	n, err := parseREQCountFromFile(reqPath)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", reqPath, err)
	}
	return n, nil
}

func scanEpicsAgainstCoverage(cwd string, epics []string, globalCoverage map[ACCode][]CoverageRef) ([]EpicSummary, []ProjectNotCoveredAC, bool, error) {
	var results []EpicSummary
	var projectNotCovered []ProjectNotCoveredAC
	hasGaps := false

	for _, epic := range epics {
		epicNum, err := getEpicNumber(epic)
		if err != nil {
			return nil, nil, false, fmt.Errorf("epic %q: %w", epic, err)
		}
		acPath := filepath.Join(cwd, "ai-sdlc-artefacts", "epics", epic, "ep-acceptance-criteria.md")
		if _, err := os.Stat(acPath); os.IsNotExist(err) {
			continue
		}
		reqPath := filepath.Join(cwd, "ai-sdlc-artefacts", "epics", epic, "ep-requirements.md")
		acs, excluded, err := parseACsFromFile(acPath)
		if err != nil {
			return nil, nil, false, fmt.Errorf("parse %s: %w", acPath, err)
		}
		reqCount, err := epicREQCount(reqPath)
		if err != nil {
			return nil, nil, false, err
		}
		epicCoverage := filterCoverageForEpicNum(globalCoverage, epicNum)
		testsCount := uniqueTestRefsCount(epicCoverage)
		r := generateReport(epic, acs, excluded, epicCoverage)
		results = append(results, EpicSummary{
			Epic:                epic,
			REQCount:            reqCount,
			TotalACs:            r.TotalACs,
			TestsCount:          testsCount,
			DeferredACs:         r.DeferredACs,
			ObsoleteACs:         r.ObsoleteACs,
			InScopeACs:          r.InScopeACs,
			AutomatedCoveredACs: r.AutomatedCoveredACs,
			ManualOnlyTracedACs: r.ManualOnlyTracedACs,
			TraceabilityRatio:   r.TraceabilityRatio,
			AutomatedRatio:      r.AutomatedRatio,
		})
		if hasBlockingGaps(r) {
			hasGaps = true
		}
		for _, gap := range r.Gaps {
			if gap.Status != "not_covered" {
				continue
			}
			crit := gap.Criterion
			if crit == "" {
				crit = acs[ACCode(gap.Code)]
			}
			crit = normalizeCriterionPreview(crit)
			projectNotCovered = append(projectNotCovered, ProjectNotCoveredAC{
				Epic:      epic,
				Code:      gap.Code,
				Criterion: crit,
			})
		}
	}
	return results, projectNotCovered, hasGaps, nil
}

// findCoverageAndTestTrace runs both codebase scans used for validation.
func findCoverageAndTestTrace(cwd string) (map[ACCode][]CoverageRef, int, []string, error) {
	globalCoverage, testFuncsWithSkip, err := findCoverageInCodebase(cwd)
	if err != nil {
		return nil, 0, nil, err
	}
	testsMissingACTrace, err := findTestsMissingACTrace(cwd)
	if err != nil {
		return nil, 0, nil, err
	}
	return globalCoverage, testFuncsWithSkip, testsMissingACTrace, nil
}

// readSortedEpicNames lists EP-* directory names under epicsPath, sorted.
func readSortedEpicNames(epicsPath string) ([]string, error) {
	entries, err := os.ReadDir(epicsPath)
	if err != nil {
		return nil, err
	}
	var epics []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "EP-") {
			epics = append(epics, entry.Name())
		}
	}
	sort.Strings(epics)
	return epics, nil
}

func allEpicsProjectWideHasGaps(hasGaps bool, testsMissingACTrace, nolintViolations []string) bool {
	if len(testsMissingACTrace) > 0 || len(nolintViolations) > 0 {
		return true
	}
	return hasGaps
}

func aggregateEpicTotals(results []EpicSummary) (totalACs, totalDeferred, totalObsolete, totalInScope, totalAuto, totalManual int, traceRatio, autoRatio float64) {
	for _, res := range results {
		totalACs += res.TotalACs
		totalDeferred += res.DeferredACs
		totalObsolete += res.ObsoleteACs
		totalInScope += res.InScopeACs
		totalAuto += res.AutomatedCoveredACs
		totalManual += res.ManualOnlyTracedACs
	}
	if totalInScope > 0 {
		traceRatio = float64(totalAuto+totalManual) / float64(totalInScope)
		autoRatio = float64(totalAuto) / float64(totalInScope)
	}
	return totalACs, totalDeferred, totalObsolete, totalInScope, totalAuto, totalManual, traceRatio, autoRatio
}
