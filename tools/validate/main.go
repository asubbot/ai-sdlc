package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type ACCode string

// acExclusionKind marks an AC as out of scope for mandatory automated test traces
// when declared in ep-acceptance-criteria.md (vision changed, superseded, or deferred work).
type acExclusionKind string

const (
	acExclusionNone     acExclusionKind = ""
	acExclusionDeferred acExclusionKind = "deferred"
	acExclusionObsolete acExclusionKind = "obsolete"
)

// CoverageRef is one test function reference with optional manual-only classification.
type CoverageRef struct {
	Ref    string `json:"ref"`
	Manual bool   `json:"manual"`
}

type Report struct {
	Epic                    string                   `json:"epic"`
	TotalACs                int                      `json:"total_acs"`
	DeferredACs             int                      `json:"deferred_acs"`
	ObsoleteACs             int                      `json:"obsolete_acs"`
	InScopeACs              int                      `json:"in_scope_acs"`
	AutomatedCoveredACs     int                      `json:"automated_covered_acs"`
	ManualOnlyTracedACs     int                      `json:"manual_only_traced_acs"`
	TraceabilityRatio       float64                  `json:"traceability_ratio"`
	AutomatedRatio          float64                  `json:"automated_ratio"`
	TestFuncsWithSkip       int                      `json:"test_funcs_with_skip"`
	TestsMissingACTrace     []string                 `json:"tests_missing_ac_trace,omitempty"`
	NolintGocycloViolations []string                 `json:"nolint_gocyclo_violations,omitempty"`
	Gaps                    []ACGap                  `json:"gaps"`
	Coverage                map[string][]CoverageRef `json:"ac_to_tests"`
}

type ACGap struct {
	Code      string `json:"code"`
	Criterion string `json:"criterion"`
	Status    string `json:"status"`
	Reason    string `json:"reason,omitempty"`
}

type EpicSummary struct {
	Epic                string  `json:"epic"`
	REQCount            int     `json:"req_count"`
	TotalACs            int     `json:"total_acs"`
	TestsCount          int     `json:"tests_count"`
	DeferredACs         int     `json:"deferred_acs"`
	ObsoleteACs         int     `json:"obsolete_acs"`
	InScopeACs          int     `json:"in_scope_acs"`
	AutomatedCoveredACs int     `json:"automated_covered_acs"`
	ManualOnlyTracedACs int     `json:"manual_only_traced_acs"`
	TraceabilityRatio   float64 `json:"traceability_ratio"`
	AutomatedRatio      float64 `json:"automated_ratio"`
}

// ProjectNotCoveredAC is one AC that has no test traceability comment (project-wide run).
type ProjectNotCoveredAC struct {
	Epic      string `json:"epic"`
	Code      string `json:"code"`
	Criterion string `json:"criterion,omitempty"`
}

type AllEpicsReport struct {
	Epics                   []EpicSummary         `json:"epics"`
	NotCoveredACs           []ProjectNotCoveredAC `json:"not_covered_acs"`
	NotCoveredCount         int                   `json:"not_covered_count"`
	TotalACs                int                   `json:"total_acs"`
	DeferredACs             int                   `json:"deferred_acs"`
	ObsoleteACs             int                   `json:"obsolete_acs"`
	InScopeACs              int                   `json:"in_scope_acs"`
	AutomatedCoveredACs     int                   `json:"automated_covered_acs"`
	ManualOnlyTracedACs     int                   `json:"manual_only_traced_acs"`
	TraceabilityRatio       float64               `json:"traceability_ratio"`
	AutomatedRatio          float64               `json:"automated_ratio"`
	TestFuncsWithSkip       int                   `json:"test_funcs_with_skip"`
	TestsMissingACTrace     []string              `json:"tests_missing_ac_trace,omitempty"`
	NolintGocycloViolations []string              `json:"nolint_gocyclo_violations,omitempty"`
	HasGaps                 bool                  `json:"has_gaps"`
}

// validateAllEpics finds and validates all epics in ai-sdlc-artefacts/epics/
func validateAllEpics(jsonOutput bool) bool {
	cwd, err := os.Getwd()
	if err != nil {
		errLog.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	nolintViolations, err := findNolintGocycloViolations(cwd)
	if err != nil {
		errLog.Printf("Error scanning for nolint:gocyclo: %v\n", err)
		os.Exit(1)
	}

	epicsPath := filepath.Join(cwd, "ai-sdlc-artefacts", "epics")
	epics, err := readSortedEpicNames(epicsPath)
	if err != nil {
		errLog.Printf("Error reading epics directory: %v\n", err)
		os.Exit(1)
	}
	if len(epics) == 0 {
		errLog.Printf("No epics found in %s\n", epicsPath)
		os.Exit(1)
	}

	if !jsonOutput {
		writeStdout("🔍 Validating AC coverage for all %d epics...\n\n", len(epics))
	}

	globalCoverage, testFuncsWithSkip, testsMissingACTrace, err := findCoverageAndTestTrace(cwd)
	if err != nil {
		errLog.Printf("Error scanning codebase: %v\n", err)
		os.Exit(1)
	}

	results, projectNotCovered, hasGaps := scanEpicsAgainstCoverage(cwd, epics, globalCoverage)
	hasGaps = allEpicsProjectWideHasGaps(hasGaps, testsMissingACTrace, nolintViolations)

	sort.Slice(projectNotCovered, func(i, j int) bool {
		if projectNotCovered[i].Epic != projectNotCovered[j].Epic {
			return projectNotCovered[i].Epic < projectNotCovered[j].Epic
		}
		return projectNotCovered[i].Code < projectNotCovered[j].Code
	})

	totalACs, totalDeferred, totalObsolete, totalInScope, totalAuto, totalManual, traceRatio, autoRatio := aggregateEpicTotals(results)

	allReport := AllEpicsReport{
		Epics:                   results,
		NotCoveredACs:           projectNotCovered,
		NotCoveredCount:         len(projectNotCovered),
		TotalACs:                totalACs,
		DeferredACs:             totalDeferred,
		ObsoleteACs:             totalObsolete,
		InScopeACs:              totalInScope,
		AutomatedCoveredACs:     totalAuto,
		ManualOnlyTracedACs:     totalManual,
		TraceabilityRatio:       traceRatio,
		AutomatedRatio:          autoRatio,
		TestFuncsWithSkip:       testFuncsWithSkip,
		TestsMissingACTrace:     testsMissingACTrace,
		NolintGocycloViolations: nolintViolations,
		HasGaps:                 hasGaps,
	}

	if jsonOutput {
		return writeAllEpicsJSON(allReport, hasGaps)
	}

	return printAllEpicsHuman(results, projectNotCovered, totalACs, totalDeferred, totalObsolete, totalInScope, totalAuto, totalManual, traceRatio, autoRatio, testFuncsWithSkip, hasGaps, testsMissingACTrace, nolintViolations)
}

// getEpicNumber extracts the numeric part from "EP-009" → "09".
func getEpicNumber(epic string) string {
	num, err := normalizeEpicNumber(epic)
	if err != nil {
		return strings.TrimPrefix(epic, "EP-")
	}
	return num
}

// jsonOutputRequested returns true if JSON output is requested.
func jsonOutputRequested(flagVal bool, argvTail []string) bool {
	if flagVal {
		return true
	}
	for _, a := range argvTail {
		if a == "--json" {
			return true
		}
	}
	return false
}

func runSingleEpicValidation(epic, epicNum, cwd string, acs map[ACCode]string, excluded map[ACCode]acExclusionKind, nolintViolations []string, jsonOut bool) bool {
	coverage, testFuncsWithSkip, testsMissingACTrace, err := findCoverageAndTestTrace(cwd)
	if err != nil {
		errLog.Printf("Error scanning codebase: %v\n", err)
		os.Exit(1)
	}
	epicCoverage := filterCoverageForEpicNum(coverage, epicNum)
	r := generateReport(epic, acs, excluded, epicCoverage)
	r.TestFuncsWithSkip = testFuncsWithSkip
	r.TestsMissingACTrace = testsMissingACTrace
	r.NolintGocycloViolations = nolintViolations
	if jsonOut {
		data, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			errLog.Printf("Error marshaling JSON: %v\n", err)
			os.Exit(1)
		}
		writelnStdout(string(data))
	} else {
		printTable(r, acs, excluded, testFuncsWithSkip, testsMissingACTrace, r.NolintGocycloViolations)
	}
	return !(hasBlockingGaps(r) || len(testsMissingACTrace) > 0 || len(nolintViolations) > 0)
}

func validateSingleEpic(epic string, jsonOut bool) bool {
	epicNum, err := normalizeEpicNumber(epic)
	if err != nil {
		errLog.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		errLog.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	nolintViolations, err := findNolintGocycloViolations(cwd)
	if err != nil {
		errLog.Printf("Error scanning for nolint:gocyclo: %v\n", err)
		os.Exit(1)
	}

	if !jsonOut {
		writeStdout("🔍 Validating AC coverage for %s...\n\n", epic)
	}

	acPath := filepath.Join(cwd, "ai-sdlc-artefacts", "epics", epic, "ep-acceptance-criteria.md")
	if _, err := os.Stat(acPath); os.IsNotExist(err) {
		errLog.Printf("Error: %s not found\n", acPath)
		writeStderr("\nUsage:\n")
		writeStderr("  validate          - Validate all epics (default)\n")
		writeStderr("  validate EP-009   - Validate single epic\n")
		writeStderr("  validate --json   - JSON output\n")
		os.Exit(1)
	}

	acs, excluded, err := parseACsFromFile(acPath)
	if err != nil {
		errLog.Printf("Error parsing AC file: %v\n", err)
		os.Exit(1)
	}

	return runSingleEpicValidation(epic, epicNum, cwd, acs, excluded, nolintViolations, jsonOut)
}

var subcommands = map[string]bool{
	"ac":        true,
	"req":       true,
	"pipeline":  true,
	"structure": true,
	"ears":      true,
}

func resolveSubcommand(args []string) (subcmd, epic string, explicit bool) {
	subcmd = "ac"
	epic = "all"
	explicit = false
	cmdArgs := args
	if len(args) > 0 && subcommands[args[0]] {
		subcmd = args[0]
		explicit = true
		cmdArgs = args[1:]
	}
	if len(cmdArgs) > 0 {
		epic = cmdArgs[0]
	}
	return subcmd, epic, explicit
}

func main() {
	if helpRequested(os.Args[1:]) {
		printUsageTo(stdout)
		return
	}

	jsonFlag := flag.Bool("json", false, "Output report as JSON instead of table")
	flag.Parse()
	jsonOut := jsonOutputRequested(*jsonFlag, os.Args[1:])

	subcmd, epic, explicit := resolveSubcommand(flag.Args())

	if !explicit {
		runFullValidation(epic, jsonOut)
		return
	}

	switch subcmd {
	case "ac":
		runACValidation(epic, jsonOut)
	case "req":
		runREQValidation(epic, jsonOut)
	case "pipeline":
		runPipelineValidation(epic, jsonOut)
	case "structure":
		runStructureValidation(epic, jsonOut)
	case "ears":
		runEARSValidation(epic, jsonOut)
	default:
		errLog.Printf("Unknown subcommand: %s\n", subcmd)
		printUsageTo(os.Stderr)
		os.Exit(1)
	}
}

func runACValidation(epic string, jsonOut bool) {
	if !acValidationPasses(epic, jsonOut) {
		os.Exit(1)
	}
}

func runPipelineValidation(epic string, jsonOut bool) {
	if epic == "all" {
		errLog.Printf("pipeline subcommand requires an epic ID (e.g., validate pipeline EP-009)\n")
		os.Exit(1)
	}
	cwd, err := os.Getwd()
	if err != nil {
		errLog.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}
	epicDir := filepath.Join(cwd, "ai-sdlc-artefacts", "epics", epic)
	result := checkPipelineState(epicDir, epic)
	if jsonOut {
		data, _ := json.MarshalIndent(result, "", "  ")
		writelnStdout(string(data))
	} else {
		printPipelineHuman(result)
	}
	if result.HasGaps {
		os.Exit(1)
	}
}

func runEARSValidation(epic string, jsonOut bool) {
	if !earsValidationPasses(epic, jsonOut) {
		os.Exit(1)
	}
}

// helpRequested reports whether argv asks for the usage text explicitly. An
// explicit request is the output the user asked for, so it belongs on stdout
// with exit code 0; usage printed because of a bad invocation stays on stderr.
func helpRequested(argv []string) bool {
	for _, arg := range argv {
		switch arg {
		case "-h", "-help", "--help", "help":
			return true
		}
	}
	return false
}

func printUsageTo(w io.Writer) {
	_, _ = fmt.Fprintf(w, `Usage: validate [subcommand] [EP-XXX] [--json]

Default (no subcommand):
  validate              AC + ears + req for all epics (ears/req skip NEW/CANCEL)
  validate EP-009       AC + ears + req for one epic (ears/req skipped when NEW/CANCEL)

Subcommands (single check):
  ac          AC coverage validation
  req         REQ <-> AC traceability check
  pipeline    Pipeline state and gate validation
  ears        EARS requirements linting
  structure   Artefact structure validation

Examples:
  validate                    # Full gate (all epics)
  validate EP-009             # Full gate (one epic)
  validate ac EP-009          # AC coverage only
  validate req EP-009         # REQ-AC traceability only
  validate req all            # REQ-AC for in-scope epics
  validate ears EP-009        # EARS linter only
  validate ears all           # EARS for in-scope epics
  validate --json             # Not supported for default gate
  validate req EP-009 --json  # JSON for explicit subcommand
`)
}

func normalizeEpicNumber(epic string) (string, error) {
	num := epic
	if strings.HasPrefix(epic, "EP-") {
		num = strings.TrimPrefix(epic, "EP-")
	}
	if num == "" {
		return "", fmt.Errorf("invalid epic id: %s", epic)
	}
	n, err := strconv.Atoi(num)
	if err != nil {
		return "", fmt.Errorf("invalid epic id: %s", epic)
	}
	return fmt.Sprintf("%02d", n), nil
}
