package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// StageStatus describes the state of a single pipeline stage artefact.
type StageStatus struct {
	Stage            int    `json:"stage"`
	Name             string `json:"name"`
	File             string `json:"file"`
	Exists           bool   `json:"exists"`
	Status           string `json:"status"`
	Gate             string `json:"gate"`
	DecisionEvidence bool   `json:"decision_evidence,omitempty"`
	HasGaps          bool   `json:"has_gaps"`
}

// GateSummary holds counts from a "Current Gate Summary" table in a review file.
type GateSummary struct {
	Blocker int
	Major   int
	Medium  int
	Minor   int
	Nit     int
}

// PipelineResult is the aggregate result of pipeline state validation.
type PipelineResult struct {
	Epic     string        `json:"epic"`
	Stages   []StageStatus `json:"stages"`
	Errors   int           `json:"errors"`
	Warnings int           `json:"warnings"`
	Findings []string      `json:"findings,omitempty"`
	HasGaps  bool          `json:"has_gaps"`
}

var pipelineStages = []struct {
	stage    int
	name     string
	file     string
	required bool
	hasGate  bool
}{
	{3, "Epic planning", "ep-scope.md", true, false},
	{4, "Requirements", "ep-requirements.md", true, false},
	{5, "Acceptance criteria", "ep-acceptance-criteria.md", true, false},
	{6, "System design", "ep-system-design.md", true, false},
	{7, "System design review", "ep-system-design-review.md", true, true},
	{8, "Implementation planning", "ep-implementation-plan.md", true, false},
	{10, "Code review", "ep-code-review.md", false, true},
	{11, "Audit", "ep-audit-report.md", false, false},
}

var (
	gateSummaryHeaderRe   = regexp.MustCompile(`(?i)##\s+Current\s+Gate\s+Summary`)
	openCountsPlainTextRe = regexp.MustCompile(`(?i)Open\s+counts?:\s*Blocker\s+(\d+)\s*\|\s*Major\s+(\d+)\s*\|\s*Medium\s+(\d+)\s*\|\s*Minor\s+(\d+)`)
	decisionNeededRe      = regexp.MustCompile(`(?im)^\s*Decision needed:\s*\S`)
	operatorChoiceRe      = regexp.MustCompile(`(?im)^\s*Operator choice:\s*\S`)
)

// parseGateSummary extracts blocker/major/medium/minor counts from a
// "## Current Gate Summary" section. Supports both the plain-text pipe format
// ("Open counts: Blocker X | Major X | ...") and a markdown table with a "Count" row.
func parseGateSummary(content string) (*GateSummary, error) {
	lines := strings.Split(content, "\n")
	headerIdx := -1
	for i, line := range lines {
		if gateSummaryHeaderRe.MatchString(line) {
			headerIdx = i
			break
		}
	}
	if headerIdx < 0 {
		return nil, fmt.Errorf("no Current Gate Summary section found")
	}

	for i := headerIdx + 1; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])

		if strings.HasPrefix(trimmed, "#") && !gateSummaryHeaderRe.MatchString(trimmed) {
			break
		}

		if m := openCountsPlainTextRe.FindStringSubmatch(trimmed); len(m) == 5 {
			blocker, err1 := strconv.Atoi(m[1])
			major, err2 := strconv.Atoi(m[2])
			medium, err3 := strconv.Atoi(m[3])
			minor, err4 := strconv.Atoi(m[4])
			if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
				return nil, fmt.Errorf("parse open counts: %w", errors.Join(err1, err2, err3, err4))
			}
			return &GateSummary{Blocker: blocker, Major: major, Medium: medium, Minor: minor}, nil
		}

		if strings.HasPrefix(trimmed, "|") {
			cells := splitTableRow(trimmed)
			if len(cells) > 0 && strings.EqualFold(cells[0], "Count") {
				return parseCountCells(cells)
			}
		}
	}
	return nil, fmt.Errorf("no open counts found in gate summary")
}

func splitTableRow(row string) []string {
	row = strings.TrimSpace(row)
	row = strings.Trim(row, "|")
	parts := strings.Split(row, "|")
	var cells []string
	for _, p := range parts {
		cells = append(cells, strings.TrimSpace(p))
	}
	return cells
}

func parseCountCells(cells []string) (*GateSummary, error) {
	if len(cells) < 6 {
		return nil, fmt.Errorf("count row has %d cells, need 6", len(cells))
	}
	nums := make([]int, 5)
	for i := 0; i < 5; i++ {
		n, err := strconv.Atoi(cells[i+1])
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as integer: %w", cells[i+1], err)
		}
		nums[i] = n
	}
	return &GateSummary{
		Blocker: nums[0],
		Major:   nums[1],
		Medium:  nums[2],
		Minor:   nums[3],
		Nit:     nums[4],
	}, nil
}

func hasBlockingCounts(counts *SeverityCounts) bool {
	return counts != nil && (counts.Blocker > 0 || counts.Major > 0 || counts.Medium > 0 || counts.Minor > 0)
}

func gateSummaryHasBlockingCounts(summary *GateSummary) bool {
	return summary != nil && (summary.Blocker > 0 || summary.Major > 0 || summary.Medium > 0 || summary.Minor > 0)
}

func normalizeGate(gate string) string {
	gate = strings.ToLower(strings.TrimSpace(gate))
	switch gate {
	case "pass", "fail", "cap":
		return gate
	default:
		return gate
	}
}

func inferGate(content string, fm *FrontMatter) string {
	gate := normalizeGate(fm.Gate)
	if hasBlockingCounts(fm.OpenCounts) {
		return "fail"
	}
	if fm.OpenCounts != nil && gate == "" {
		return "pass"
	}
	if gs, gsErr := parseGateSummary(content); gsErr == nil {
		if gateSummaryHasBlockingCounts(gs) {
			if gate == "cap" {
				return "cap"
			}
			return "fail"
		}
		if gate == "" {
			return "pass"
		}
	}
	return gate
}

func hasOperatorDecisionEvidence(content string) bool {
	return decisionNeededRe.MatchString(content) && operatorChoiceRe.MatchString(content)
}

var uncheckedTaskRe = regexp.MustCompile(`(?m)^\s*-\s*\[\s\]\s+`)

// countUncheckedPlanTasks returns the number of unchecked `- [ ]` items in the plan body.
// Stage 9 has no artefact file; unchecked tasks in ep-implementation-plan.md indicate incomplete work.
func countUncheckedPlanTasks(content string) int {
	inTasks := false
	count := 0
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			inTasks = strings.Contains(strings.ToLower(trimmed), "tasks")
			continue
		}
		if inTasks && uncheckedTaskRe.MatchString(line) {
			count++
		}
	}
	return count
}

// checkPipelineState validates the pipeline state for a single epic directory.
func checkPipelineState(epicDir, epicID string) *PipelineResult {
	result := &PipelineResult{
		Epic:   epicID,
		Stages: make([]StageStatus, 0, len(pipelineStages)),
	}

	highestExistingIdx := -1
	stageReadFailed := make([]bool, len(pipelineStages))
	var implPlanRaw []byte
	implPlanReadOK := false
	implPlanHardFailureRecorded := false

	for idx, ps := range pipelineStages {
		ss := StageStatus{
			Stage: ps.stage,
			Name:  ps.name,
			File:  ps.file,
		}

		filePath := epicDir + "/" + ps.file
		data, err := os.ReadFile(filePath)
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				stageReadFailed[idx] = true
				ss.Exists = true
				ss.Status = "error"
				ss.HasGaps = true
				result.Errors++
				msg := fmt.Sprintf(
					"stage %d (%s): cannot read artefact: %v",
					ps.stage, ps.file, err,
				)
				errLog.Printf("%s\n", msg)
				result.Findings = append(result.Findings, msg)
				if ps.stage == 8 {
					implPlanHardFailureRecorded = true
				}
				if idx > highestExistingIdx {
					highestExistingIdx = idx
				}
			} else {
				ss.Status = "missing"
			}
			result.Stages = append(result.Stages, ss)
			continue
		}

		ss.Exists = true
		if idx > highestExistingIdx {
			highestExistingIdx = idx
		}
		if ps.stage == 8 {
			implPlanReadOK = true
			implPlanRaw = append([]byte(nil), data...)
		}

		content := string(data)
		fm, fmErr := parseFrontMatter(content)
		if fmErr != nil {
			ss.Status = "error"
			ss.HasGaps = true
			result.Errors++
			result.Findings = append(result.Findings, fmt.Sprintf("stage %d (%s): %s", ps.stage, ps.file, fmErr.Error()))
		} else {
			ss.Status = fm.Status
			if ps.hasGate {
				ss.Gate = inferGate(content, fm)
				ss.DecisionEvidence = hasOperatorDecisionEvidence(content)
			}
		}

		result.Stages = append(result.Stages, ss)
	}

	// Ordering check: if a later stage exists, all earlier required stages must exist.
	if highestExistingIdx >= 0 {
		for idx := 0; idx < highestExistingIdx; idx++ {
			ps := pipelineStages[idx]
			if !ps.required {
				continue
			}
			ss := &result.Stages[idx]
			if !ss.Exists {
				ss.HasGaps = true
				result.Errors++
				result.Findings = append(result.Findings, fmt.Sprintf(
					"stage %d (%s) is missing but later stage %d (%s) exists",
					ps.stage, ps.file,
					pipelineStages[highestExistingIdx].stage, pipelineStages[highestExistingIdx].file,
				))
				if ps.stage == 8 {
					implPlanHardFailureRecorded = true
				}
			}
		}
	}

	// Gate violation: if stage 8 (impl plan) exists but stage 7 (review)
	// gate is not pass and has no recorded operator decision.
	reviewIdx := -1
	implIdx := -1
	for idx, ps := range pipelineStages {
		if ps.stage == 7 {
			reviewIdx = idx
		}
		if ps.stage == 8 {
			implIdx = idx
		}
	}
	if reviewIdx >= 0 && implIdx >= 0 &&
		result.Stages[implIdx].Exists && result.Stages[reviewIdx].Exists &&
		!stageReadFailed[reviewIdx] &&
		result.Stages[reviewIdx].Gate != "pass" && !result.Stages[reviewIdx].DecisionEvidence {
		result.Stages[reviewIdx].HasGaps = true
		result.Errors++
		result.Findings = append(result.Findings, fmt.Sprintf(
			"stage 8 (%s) exists but stage 7 gate=%s (must be pass or have recorded operator decision)",
			pipelineStages[implIdx].file, result.Stages[reviewIdx].Gate,
		))
	}

	// Gate violation: if stage 11 (audit) exists, stage 10 (code review)
	// gate evidence must exist and be pass or have a recorded operator decision.
	codeReviewIdx := -1
	auditIdx := -1
	for idx, ps := range pipelineStages {
		if ps.stage == 10 {
			codeReviewIdx = idx
		}
		if ps.stage == 11 {
			auditIdx = idx
		}
	}
	if codeReviewIdx >= 0 && auditIdx >= 0 && result.Stages[auditIdx].Exists {
		if !result.Stages[codeReviewIdx].Exists {
			result.Stages[codeReviewIdx].HasGaps = true
			result.Errors++
			result.Findings = append(result.Findings, fmt.Sprintf(
				"stage 11 (%s) exists but stage 10 (%s) gate evidence is missing",
				pipelineStages[auditIdx].file, pipelineStages[codeReviewIdx].file,
			))
		} else if !stageReadFailed[codeReviewIdx] &&
			result.Stages[codeReviewIdx].Gate != "pass" && !result.Stages[codeReviewIdx].DecisionEvidence {
			result.Stages[codeReviewIdx].HasGaps = true
			result.Errors++
			result.Findings = append(result.Findings, fmt.Sprintf(
				"stage 11 (%s) exists but stage 10 gate=%s (must be pass or have recorded operator decision)",
				pipelineStages[auditIdx].file, result.Stages[codeReviewIdx].Gate,
			))
		}
	}

	// Stage 9 is the codebase (no artefact). When a later stage exists, flag unchecked plan tasks.
	auditExists := auditIdx >= 0 && result.Stages[auditIdx].Exists
	codeReviewExists := codeReviewIdx >= 0 && result.Stages[codeReviewIdx].Exists
	if implPlanReadOK {
		unchecked := countUncheckedPlanTasks(string(implPlanRaw))
		if unchecked > 0 {
			if auditExists {
				result.Errors++
				result.HasGaps = true
				result.Findings = append(result.Findings, fmt.Sprintf(
					"stage 11 (%s) exists but ep-implementation-plan.md has %d unchecked task(s)",
					pipelineStages[auditIdx].file, unchecked,
				))
			} else if codeReviewExists {
				result.Warnings++
				result.Findings = append(result.Findings, fmt.Sprintf(
					"stage 10 (%s) exists but ep-implementation-plan.md has %d unchecked task(s)",
					pipelineStages[codeReviewIdx].file, unchecked,
				))
			}
		}
	} else if (auditExists || codeReviewExists) && !implPlanHardFailureRecorded && implIdx >= 0 {
		// Belt-and-braces: ordering and direct read-error paths should already have
		// recorded the actionable failure before we reach checkbox evaluation.
		result.Stages[implIdx].HasGaps = true
		result.Errors++
		result.Findings = append(result.Findings, fmt.Sprintf(
			"stage 8 (%s) must be readable before checkbox evaluation for later stages",
			pipelineStages[implIdx].file,
		))
	}

	result.HasGaps = result.Errors > 0
	return result
}

func printPipelineHuman(result *PipelineResult) {
	writeStdout("\n📋 Pipeline State for %s\n\n", result.Epic)

	writeStdout("%-7s %-24s %-31s %-11s %s\n", "Stage", "Name", "File", "Status", "Gate")
	writelnStdout(strings.Repeat("─", 82))

	optionalMissing := 0
	for _, ss := range result.Stages {
		status := ss.Status
		if !ss.Exists {
			status = "MISSING"
		}
		gate := "—"
		if ss.Gate != "" {
			gate = ss.Gate
		}

		writeStdout("%3d    %-24s %-31s %-11s %s\n", ss.Stage, ss.Name, ss.File, status, gate)

		if !ss.Exists {
			for _, ps := range pipelineStages {
				if ps.stage == ss.Stage && !ps.required {
					optionalMissing++
					break
				}
			}
		}
	}

	writelnStdout("")
	if result.HasGaps {
		writeStdout("❌ %d error(s) found", result.Errors)
		if len(result.Findings) > 0 {
			writelnStdout(":")
			for _, f := range result.Findings {
				writeStdout("  • %s\n", f)
			}
		} else {
			writelnStdout("")
		}
	} else {
		suffix := ""
		if optionalMissing > 0 {
			suffix = fmt.Sprintf(" | %d optional stages missing", optionalMissing)
		}
		writeStdout("✅ No gate violations%s\n", suffix)
	}
	writelnStdout("")
}
