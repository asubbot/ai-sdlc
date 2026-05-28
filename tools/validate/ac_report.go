package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

func acCoverageCell(refs []CoverageRef, exclusion acExclusionKind) (status, testStr string) {
	switch exclusion {
	case acExclusionObsolete:
		return "↷", "OBSOLETE"
	case acExclusionDeferred:
		return "↷", "DEFERRED"
	}
	if len(refs) == 0 {
		return "✗", "NOT COVERED"
	}
	hasAuto := false
	for _, ref := range refs {
		if !ref.Manual {
			hasAuto = true
			break
		}
	}
	if !hasAuto {
		if len(refs) == 1 {
			return "✎", "MANUAL " + refs[0].Ref
		}
		return "✎", fmt.Sprintf("%d MANUAL", len(refs))
	}
	if len(refs) == 1 {
		return "✓", refs[0].Ref
	}
	return "✓", fmt.Sprintf("%d tests", len(refs))
}

// printTable prints a formatted table report for a single epic.
func printTable(r *Report, acs map[ACCode]string, excluded map[ACCode]acExclusionKind, testFuncsWithSkip int, testsMissingACTrace []string, gocycloSuppressViolations []string) {
	writeStdout("\n📋 AC Coverage Report for %s\n\n", r.Epic)
	writeStdout("%-15s %-50s %-30s\n", "AC Code", "Criterion", "Coverage")
	writelnStdout(strings.Repeat("─", 95))

	keys := make([]ACCode, 0, len(acs))
	for ac := range acs {
		keys = append(keys, ac)
	}
	sort.Slice(keys, func(i, j int) bool {
		return string(keys[i]) < string(keys[j])
	})

	for _, ac := range keys {
		criterion := acs[ac]
		if len(criterion) > 47 {
			criterion = criterion[:44] + "..."
		}

		refs := r.Coverage[string(ac)]
		status, testStr := acCoverageCell(refs, excluded[ac])

		if len(testStr) > 27 {
			testStr = testStr[:24] + "..."
		}

		writeStdout("%s %-15s %-50s %-30s\n", status, string(ac), criterion, testStr)
	}

	writelnStdout(strings.Repeat("─", 95))

	coverageStr := "❌"
	if r.TraceabilityRatio == 1.0 {
		coverageStr = "✅"
	} else if r.TraceabilityRatio >= 0.9 {
		coverageStr = "⚠️"
	}

	writeStdout("\n%s RESULT: in-scope %d/%d traced (%.1f%%), automated %d (%.1f%%), manual-only %d | deferred %d | obsolete %d | total ACs %d\n",
		coverageStr,
		r.AutomatedCoveredACs+r.ManualOnlyTracedACs, r.InScopeACs, r.TraceabilityRatio*100,
		r.AutomatedCoveredACs, r.AutomatedRatio*100,
		r.ManualOnlyTracedACs,
		r.DeferredACs, r.ObsoleteACs, r.TotalACs,
	)
	writeStdout("   Project-wide: Test functions with t.Skip: %d\n", testFuncsWithSkip)

	if hasBlockingGaps(r) {
		writeStdout("\n❌ Missing coverage for:\n")
		for _, gap := range r.Gaps {
			if gap.Status == "not_covered" {
				writeStdout("  • %s\n", gap.Code)
			}
		}
		writeStdout("\nAction: Add tests for missing ACs, or mark the AC **Obsolete** / **Deferred** in ep-acceptance-criteria.md (see VALIDATION.md)\n")
	} else {
		writeStdout("\n✅ All ACs covered — epic is ready for audit\n")
	}

	if len(testsMissingACTrace) > 0 {
		writeStdout("\n❌ Test functions without AC trace comment (project-wide): %d\n", len(testsMissingACTrace))
		for _, ref := range testsMissingACTrace {
			writeStdout("  • %s\n", ref)
		}
		writeStdout("\nAction: Add a trace line (e.g. // Covers AC-EE.NNN) bound to each Test* per VALIDATION.md\n")
	}
	printNolintGocycloViolationsHuman(gocycloSuppressViolations)
	writelnStdout("")
}

// printProjectNotCoveredHuman prints the AC-not-covered block (may be zero items).
func printProjectNotCoveredHuman(projectNotCovered []ProjectNotCoveredAC) {
	writelnStdout("")
	writeStdout("❌ AC not covered by tests (project-wide): %d\n", len(projectNotCovered))
	if len(projectNotCovered) == 0 {
		return
	}
	currentEpic := ""
	for _, item := range projectNotCovered {
		if item.Epic != currentEpic {
			currentEpic = item.Epic
			writeStdout("\n%s\n", currentEpic)
		}
		line := fmt.Sprintf("  • %s", item.Code)
		if c := strings.TrimSpace(item.Criterion); c != "" {
			if len(c) > 72 {
				c = c[:69] + "..."
			}
			line += " — " + c
		}
		writelnStdout(line)
	}
	writelnStdout("")
}

// printTestsMissingACTraceHuman prints the Test*-without-trace block.
func printTestsMissingACTraceHuman(testsMissingACTrace []string) {
	writelnStdout("")
	writeStdout("❌ Test functions without AC trace comment (project-wide): %d\n", len(testsMissingACTrace))
	for _, ref := range testsMissingACTrace {
		writeStdout("  • %s\n", ref)
	}
	writelnStdout("")
	writeStdout("Action: Add a trace line (e.g. // Covers AC-EE.NNN) bound to each Test* per VALIDATION.md\n")
}

func printAllEpicsHuman(
	results []EpicSummary,
	projectNotCovered []ProjectNotCoveredAC,
	totalACs, totalDeferred, totalObsolete, totalInScope, totalAuto, totalManual int,
	traceRatio, autoRatio float64,
	testFuncsWithSkip int,
	hasGaps bool,
	testsMissingACTrace []string,
	nolintGocycloViolations []string,
) {
	writelnStdout("📋 Epic Validation Summary")
	writelnStdout("")
	writeStdout("%-10s %-8s %-6s %-6s %-8s\n", "Epic", "Trace", "REQ", "AC", "Tests")
	writelnStdout(strings.Repeat("─", 42))

	totalREQ := 0
	totalTests := 0
	for _, res := range results {
		pct := int(res.TraceabilityRatio * 100)
		writeStdout("%-10s %3d%%     %-6d %-6d %-8d\n", res.Epic, pct, res.REQCount, res.InScopeACs, res.TestsCount)
		totalREQ += res.REQCount
		totalTests += res.TestsCount
	}

	writelnStdout(strings.Repeat("─", 42))
	writeStdout("%-10s %3d%%     %-6d %-6d %-8d\n", "TOTAL", int(traceRatio*100), totalREQ, totalInScope, totalTests)
	writelnStdout(strings.Repeat("─", 42))

	hasTestTraceGaps := len(testsMissingACTrace) > 0
	hasNolintGocyclo := len(nolintGocycloViolations) > 0
	overallFail := hasGaps || hasTestTraceGaps || hasNolintGocyclo

	statusEmoji := "✅"
	if overallFail {
		statusEmoji = "❌"
	}

	writeStdout("\n%s OVERALL: in-scope %d/%d traced (%.1f%%), automated %d (%.1f%%), manual-only %d | deferred %d | obsolete %d | total ACs %d\n",
		statusEmoji, totalAuto+totalManual, totalInScope, traceRatio*100,
		totalAuto, autoRatio*100, totalManual, totalDeferred, totalObsolete, totalACs)
	writeStdout("   Project-wide: Test functions with t.Skip: %d\n", testFuncsWithSkip)

	if !overallFail {
		return
	}

	if hasGaps {
		printProjectNotCoveredHuman(projectNotCovered)
	}
	if hasTestTraceGaps {
		printTestsMissingACTraceHuman(testsMissingACTrace)
	}
	if hasNolintGocyclo {
		printNolintGocycloViolationsHuman(nolintGocycloViolations)
	}

	writelnStdout("Tip: run `./tools/validate/validate EP-XXX` for per-AC detail and test refs.")
	os.Exit(1)
}

func writeAllEpicsJSON(allReport AllEpicsReport, hasGaps bool) {
	data, err := json.MarshalIndent(allReport, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}
	writelnStdout(string(data))
	if hasGaps {
		os.Exit(1)
	}
}
