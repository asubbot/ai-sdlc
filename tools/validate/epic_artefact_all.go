package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type AllEARSEpicSummary struct {
	Epic     string `json:"epic"`
	Total    int    `json:"total_reqs"`
	Errors   int    `json:"errors"`
	Warnings int    `json:"warnings"`
	HasGaps  bool   `json:"has_gaps"`
}

type AllEARSReport struct {
	InScopeEpics []string             `json:"in_scope_epics"`
	SkippedEpics []string             `json:"skipped_epics"`
	Epics        []AllEARSEpicSummary `json:"epics"`
	TotalREQs    int                  `json:"total_reqs"`
	TotalErrors  int                  `json:"total_errors"`
	TotalWarnings int                 `json:"total_warnings"`
	HasGaps      bool                 `json:"has_gaps"`
}

type AllREQEpicSummary struct {
	Epic          string  `json:"epic"`
	TotalREQs     int     `json:"total_reqs"`
	CoveredREQs   int     `json:"covered_reqs"`
	CoverageRatio float64 `json:"coverage_ratio"`
	OrphanRefs    int     `json:"orphan_refs"`
	ACsWithoutREQ int     `json:"acs_without_req"`
	HasGaps       bool    `json:"has_gaps"`
}

type AllREQReport struct {
	InScopeEpics []string            `json:"in_scope_epics"`
	SkippedEpics []string            `json:"skipped_epics"`
	Epics        []AllREQEpicSummary `json:"epics"`
	TotalREQs    int                 `json:"total_reqs"`
	CoveredREQs  int                 `json:"covered_reqs"`
	HasGaps      bool                `json:"has_gaps"`
}

func epicArtefactsRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, "ai-sdlc-artefacts", "epics"), nil
}

func loadInScopeEpics() (epicsPath string, inScope, skipped []string, err error) {
	epicsPath, err = epicArtefactsRoot()
	if err != nil {
		return "", nil, nil, err
	}
	inScope, skipped, err = readInScopeEpicNames(epicsPath)
	if err != nil {
		return "", nil, nil, err
	}
	if len(inScope) == 0 && len(skipped) == 0 {
		return "", nil, nil, fmt.Errorf("no epics found in %s", epicsPath)
	}
	return epicsPath, inScope, skipped, nil
}

func validateAllEARS(jsonOut bool) {
	epicsPath, inScope, skipped, err := loadInScopeEpics()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	report := AllEARSReport{
		InScopeEpics: inScope,
		SkippedEpics: skipped,
	}

	for _, epic := range inScope {
		reqPath := filepath.Join(epicsPath, epic, "ep-requirements.md")
		var result *EARSResult
		if _, statErr := os.Stat(reqPath); os.IsNotExist(statErr) {
			result = &EARSResult{
				Epic: epic,
				Findings: []EARSFinding{{
					Severity: "error",
					Rule:     "ears-file-error",
					Message:  fmt.Sprintf("missing ep-requirements.md for in-scope epic %s", epic),
				}},
				Errors:  1,
				HasGaps: true,
			}
		} else {
			result = lintEARSFile(reqPath)
			result.Epic = epic
		}
		report.Epics = append(report.Epics, AllEARSEpicSummary{
			Epic:     epic,
			Total:    result.Total,
			Errors:   result.Errors,
			Warnings: result.Warnings,
			HasGaps:  result.HasGaps,
		})
		report.TotalREQs += result.Total
		report.TotalErrors += result.Errors
		report.TotalWarnings += result.Warnings
		if result.HasGaps {
			report.HasGaps = true
		}
	}

	if jsonOut {
		data, marshalErr := json.MarshalIndent(report, "", "  ")
		if marshalErr != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", marshalErr)
			os.Exit(1)
		}
		writelnStdout(string(data))
	} else {
		printAllEARSHuman(report)
	}
	if report.HasGaps {
		os.Exit(1)
	}
}

func validateAllREQ(jsonOut bool) {
	epicsPath, inScope, skipped, err := loadInScopeEpics()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	report := AllREQReport{
		InScopeEpics: inScope,
		SkippedEpics: skipped,
	}

	for _, epic := range inScope {
		result, validateErr := validateREQForEpic(epicsPath, epic)
		if validateErr != nil {
			fmt.Fprintf(os.Stderr, "Error validating %s: %v\n", epic, validateErr)
			os.Exit(1)
		}
		report.Epics = append(report.Epics, AllREQEpicSummary{
			Epic:          epic,
			TotalREQs:     result.TotalREQs,
			CoveredREQs:   result.CoveredREQs,
			CoverageRatio: result.CoverageRatio,
			OrphanRefs:    len(result.OrphanACRefs),
			ACsWithoutREQ: len(result.ACsWithoutREQ),
			HasGaps:       result.HasGaps,
		})
		report.TotalREQs += result.TotalREQs
		report.CoveredREQs += result.CoveredREQs
		if result.HasGaps {
			report.HasGaps = true
		}
	}

	if jsonOut {
		data, marshalErr := json.MarshalIndent(report, "", "  ")
		if marshalErr != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", marshalErr)
			os.Exit(1)
		}
		writelnStdout(string(data))
	} else {
		printAllREQHuman(report)
	}
	if report.HasGaps {
		os.Exit(1)
	}
}

func validateREQForEpic(epicsPath, epic string) (*REQACTraceResult, error) {
	epicDir := filepath.Join(epicsPath, epic)
	reqPath := filepath.Join(epicDir, "ep-requirements.md")
	acPath := filepath.Join(epicDir, "ep-acceptance-criteria.md")

	if _, err := os.Stat(reqPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("missing ep-requirements.md")
	}
	if _, err := os.Stat(acPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("missing ep-acceptance-criteria.md")
	}

	reqs, err := parseREQsFromFile(reqPath)
	if err != nil {
		return nil, err
	}
	acs, _, err := parseACsFromFile(acPath)
	if err != nil {
		return nil, err
	}
	acReqRefs, err := parseREQRefsFromACFile(acPath)
	if err != nil {
		return nil, err
	}
	result := checkREQACTraceability(reqs, acReqRefs, acs)
	result.Epic = epic
	return result, nil
}

func printAllEARSHuman(report AllEARSReport) {
	writeStdout("🔍 Validating EARS for %d in-scope epics", len(report.InScopeEpics))
	if len(report.SkippedEpics) > 0 {
		writeStdout(" (%d skipped: %s)", len(report.SkippedEpics), strings.Join(report.SkippedEpics, ", "))
	}
	writelnStdout("...\n")

	writeStdout("%-10s %8s %8s %8s %s\n", "Epic", "REQs", "Errors", "Warn", "Status")
	writelnStdout(strings.Repeat("─", 52))
	for _, epic := range report.Epics {
		icon := "✓"
		if epic.HasGaps {
			icon = "✗"
		}
		writeStdout("%-10s %8d %8d %8d %s\n", epic.Epic, epic.Total, epic.Errors, epic.Warnings, icon)
	}
	writelnStdout(strings.Repeat("─", 52))

	statusIcon := "✅"
	if report.HasGaps {
		statusIcon = "❌"
	}
	writeStdout("\n%s OVERALL: %d requirements, %d errors, %d warnings\n\n",
		statusIcon, report.TotalREQs, report.TotalErrors, report.TotalWarnings)
}

func printAllREQHuman(report AllREQReport) {
	writeStdout("🔍 Validating REQ↔AC for %d in-scope epics", len(report.InScopeEpics))
	if len(report.SkippedEpics) > 0 {
		writeStdout(" (%d skipped: %s)", len(report.SkippedEpics), strings.Join(report.SkippedEpics, ", "))
	}
	writelnStdout("...\n")

	writeStdout("%-10s %12s %8s %8s %s\n", "Epic", "REQ coverage", "Orphans", "No REQ", "Status")
	writelnStdout(strings.Repeat("─", 58))
	for _, epic := range report.Epics {
		icon := "✓"
		if epic.HasGaps {
			icon = "✗"
		}
		cov := fmt.Sprintf("%d/%d", epic.CoveredREQs, epic.TotalREQs)
		writeStdout("%-10s %12s %8d %8d %s\n", epic.Epic, cov, epic.OrphanRefs, epic.ACsWithoutREQ, icon)
	}
	writelnStdout(strings.Repeat("─", 58))

	statusIcon := "✅"
	if report.HasGaps {
		statusIcon = "❌"
	}
	writeStdout("\n%s OVERALL: %d/%d REQs covered across in-scope epics\n\n",
		statusIcon, report.CoveredREQs, report.TotalREQs)
}
