package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

type REQCode string

type REQACTraceResult struct {
	Epic          string              `json:"epic"`
	TotalREQs     int                 `json:"total_reqs"`
	TotalACs      int                 `json:"total_acs"`
	CoveredREQs   int                 `json:"covered_reqs"`
	UncoveredREQs []string            `json:"uncovered_reqs,omitempty"`
	OrphanACRefs  []string            `json:"orphan_ac_refs,omitempty"`
	ACsWithoutREQ []string            `json:"acs_without_req,omitempty"`
	REQToACs      map[string][]string `json:"req_to_acs"`
	CoverageRatio float64             `json:"coverage_ratio"`
	HasGaps       bool                `json:"has_gaps"`
}

var reqHeadingPattern = regexp.MustCompile(`^###\s+(REQ-\d{2,3}\.\d{3})\s*[—–]\s*(.*)`)

// parseREQsFromFile extracts REQ codes and their summary text from headings like
// ### REQ-EE.NNN — Summary
func parseREQsFromFile(path string) (map[REQCode]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	reqs := make(map[REQCode]string)
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		matches := reqHeadingPattern.FindStringSubmatch(line)
		if len(matches) < 3 {
			continue
		}
		code := REQCode(matches[1])
		summary := strings.TrimSpace(matches[2])
		if _, exists := reqs[code]; !exists {
			reqs[code] = summary
		}
	}

	return reqs, nil
}

var reqCodePattern = regexp.MustCompile(`REQ-(\d{2,3})\.(\d{3})`)

// parseREQRefsFromACFile parses ep-acceptance-criteria.md and extracts REQ references
// from each AC block. Looks for patterns like:
// - (Trace: REQ-EE.NNN) or (Trace: [REQ-EE.NNN](...))
// - REQ-EE.NNN on the same line as AC-EE.NNN
// - Lines containing both a REQ code and an AC code
func parseREQRefsFromACFile(path string) (map[ACCode][]REQCode, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	acPattern := regexp.MustCompile(`AC-(\d{2,3})\.(\d{3})`)
	result := make(map[ACCode][]REQCode)
	lines := strings.Split(string(content), "\n")

	var currentAC ACCode

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		acMatches := acPattern.FindAllStringSubmatch(line, -1)
		reqMatches := reqCodePattern.FindAllStringSubmatch(line, -1)

		// Any markdown heading starts a new section; unless it is an AC heading, close AC block.
		if strings.HasPrefix(trimmed, "#") && !strings.HasPrefix(trimmed, "### AC-") {
			currentAC = ""
		}

		// Track current AC from heading lines (### AC-EE.NNN ...)
		if strings.HasPrefix(trimmed, "###") || strings.HasPrefix(trimmed, "**AC-") {
			// Enter a new AC block when the heading contains AC code.
			if len(acMatches) > 0 {
				code := fmt.Sprintf("AC-%s.%s", acMatches[0][1], acMatches[0][2])
				currentAC = ACCode(code)
				if _, exists := result[currentAC]; !exists {
					result[currentAC] = nil
				}
			} else {
				// Any non-AC heading closes the previous AC block.
				currentAC = ""
			}
		}

		// If line has both AC and REQ codes, associate them
		if len(acMatches) > 0 && len(reqMatches) > 0 {
			for _, acMatch := range acMatches {
				ac := ACCode(fmt.Sprintf("AC-%s.%s", acMatch[1], acMatch[2]))
				for _, reqMatch := range reqMatches {
					req := REQCode(fmt.Sprintf("REQ-%s.%s", reqMatch[1], reqMatch[2]))
					if !containsREQ(result[ac], req) {
						result[ac] = append(result[ac], req)
					}
				}
			}
		} else if len(reqMatches) > 0 && currentAC != "" {
			// REQ on a line within the current AC block (e.g. Trace line)
			for _, reqMatch := range reqMatches {
				req := REQCode(fmt.Sprintf("REQ-%s.%s", reqMatch[1], reqMatch[2]))
				if !containsREQ(result[currentAC], req) {
					result[currentAC] = append(result[currentAC], req)
				}
			}
		}
	}

	return result, nil
}

func containsREQ(slice []REQCode, code REQCode) bool {
	for _, c := range slice {
		if c == code {
			return true
		}
	}
	return false
}

// checkREQACTraceability performs three checks:
// 1. Forward: every REQ should be covered by at least one AC
// 2. Reverse: every REQ referenced by an AC should exist in requirements
// 3. Orphan: ACs that don't reference any REQ
func checkREQACTraceability(reqs map[REQCode]string, acReqRefs map[ACCode][]REQCode, acs map[ACCode]string) *REQACTraceResult {
	result := &REQACTraceResult{
		TotalREQs: len(reqs),
		TotalACs:  len(acs),
		REQToACs:  make(map[string][]string),
	}

	// Build REQ → ACs mapping
	for ac, reqRefs := range acReqRefs {
		for _, req := range reqRefs {
			result.REQToACs[string(req)] = append(result.REQToACs[string(req)], string(ac))
		}
	}

	// Sort AC lists for deterministic output
	for req := range result.REQToACs {
		sort.Strings(result.REQToACs[req])
	}

	// Forward check: every REQ should be covered by at least one AC
	coveredCount := 0
	var uncoveredREQs []string
	for req := range reqs {
		if refs, ok := result.REQToACs[string(req)]; ok && len(refs) > 0 {
			coveredCount++
		} else {
			uncoveredREQs = append(uncoveredREQs, string(req))
		}
	}
	sort.Strings(uncoveredREQs)
	result.CoveredREQs = coveredCount
	result.UncoveredREQs = uncoveredREQs

	// Reverse check: every REQ referenced by an AC should exist
	var orphanRefs []string
	seenOrphan := make(map[string]bool)
	for _, reqRefs := range acReqRefs {
		for _, req := range reqRefs {
			if _, exists := reqs[req]; !exists {
				if !seenOrphan[string(req)] {
					orphanRefs = append(orphanRefs, string(req))
					seenOrphan[string(req)] = true
				}
			}
		}
	}
	sort.Strings(orphanRefs)
	result.OrphanACRefs = orphanRefs

	// Orphan ACs: ACs that don't reference any REQ
	var acsWithoutREQ []string
	for ac := range acs {
		refs := acReqRefs[ac]
		if len(refs) == 0 {
			acsWithoutREQ = append(acsWithoutREQ, string(ac))
		}
	}
	sort.Strings(acsWithoutREQ)
	result.ACsWithoutREQ = acsWithoutREQ

	// Coverage ratio
	if result.TotalREQs > 0 {
		result.CoverageRatio = float64(result.CoveredREQs) / float64(result.TotalREQs)
	}

	// HasGaps if any issue found
	result.HasGaps = len(result.UncoveredREQs) > 0 || len(result.OrphanACRefs) > 0 || len(result.ACsWithoutREQ) > 0

	return result
}

// printREQACTraceHuman outputs a human-readable traceability report.
func printREQACTraceHuman(result *REQACTraceResult) {
	writeStdout("\n📋 REQ↔AC Traceability Report for %s\n\n", result.Epic)

	// REQ → ACs table
	writeStdout("%-15s %s\n", "REQ Code", "Covered by ACs")
	writelnStdout(strings.Repeat("─", 60))

	reqCodes := make([]string, 0, len(result.REQToACs))
	for req := range result.REQToACs {
		reqCodes = append(reqCodes, req)
	}
	sort.Strings(reqCodes)
	for _, req := range reqCodes {
		acs := result.REQToACs[req]
		writeStdout("%-15s %s\n", req, strings.Join(acs, ", "))
	}

	writelnStdout(strings.Repeat("─", 60))

	// Uncovered REQs
	if len(result.UncoveredREQs) > 0 {
		writeStdout("\n❌ Uncovered REQs (%d):\n", len(result.UncoveredREQs))
		for _, req := range result.UncoveredREQs {
			writeStdout("  • %s\n", req)
		}
	}

	// Orphan AC refs
	if len(result.OrphanACRefs) > 0 {
		writeStdout("\n⚠️  Orphan REQ references in ACs (REQ not in requirements):\n")
		for _, ref := range result.OrphanACRefs {
			writeStdout("  • %s\n", ref)
		}
	}

	// ACs without REQ
	if len(result.ACsWithoutREQ) > 0 {
		writeStdout("\n⚠️  ACs without REQ reference (%d):\n", len(result.ACsWithoutREQ))
		for _, ac := range result.ACsWithoutREQ {
			writeStdout("  • %s\n", ac)
		}
	}

	// Summary
	coverageStr := "❌"
	if result.CoverageRatio == 1.0 && !result.HasGaps {
		coverageStr = "✅"
	} else if result.CoverageRatio >= 0.8 {
		coverageStr = "⚠️"
	}

	writeStdout("\n%s RESULT: %d/%d REQs covered (%.1f%%) | %d ACs total | %d orphan refs | %d ACs without REQ\n",
		coverageStr,
		result.CoveredREQs, result.TotalREQs, result.CoverageRatio*100,
		result.TotalACs,
		len(result.OrphanACRefs),
		len(result.ACsWithoutREQ),
	)
	writelnStdout("")
}

func runREQValidation(epic string, jsonOut bool) {
	if epic == "all" || epic == "" {
		validateAllREQ(jsonOut)
		return
	}
	epicsPath, err := epicArtefactsRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}
	result, err := validateREQForEpic(epicsPath, epic)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error validating %s: %v\n", epic, err)
		os.Exit(1)
	}

	if jsonOut {
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
			os.Exit(1)
		}
		writelnStdout(string(data))
	} else {
		printREQACTraceHuman(result)
	}
	if result.HasGaps {
		os.Exit(1)
	}
}
