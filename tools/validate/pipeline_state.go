package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// StageStatus describes the state of a single pipeline stage artefact.
type StageStatus struct {
	Stage   int    `json:"stage"`
	Name    string `json:"name"`
	File    string `json:"file"`
	Exists  bool   `json:"exists"`
	Status  string `json:"status"`
	Gate    string `json:"gate"`
	HasGaps bool   `json:"has_gaps"`
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

var gateSummaryHeaderRe = regexp.MustCompile(`(?i)##\s+Current\s+Gate\s+Summary`)

// parseGateSummary extracts blocker/major/medium/minor/nit counts from a
// "## Current Gate Summary" markdown table.
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

	// Find the "Count" row — it's the data row after the header row and separator.
	for i := headerIdx + 1; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if !strings.HasPrefix(trimmed, "|") {
			continue
		}
		cells := splitTableRow(trimmed)
		if len(cells) == 0 {
			continue
		}
		if strings.EqualFold(cells[0], "Count") {
			return parseCountCells(cells)
		}
	}
	return nil, fmt.Errorf("no Count row found in gate summary table")
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
		return nil, fmt.Errorf("Count row has %d cells, need 6", len(cells))
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

// checkPipelineState validates the pipeline state for a single epic directory.
func checkPipelineState(epicDir, epicID string) *PipelineResult {
	result := &PipelineResult{
		Epic:   epicID,
		Stages: make([]StageStatus, 0, len(pipelineStages)),
	}

	highestExistingIdx := -1

	for idx, ps := range pipelineStages {
		ss := StageStatus{
			Stage: ps.stage,
			Name:  ps.name,
			File:  ps.file,
		}

		filePath := epicDir + "/" + ps.file
		data, err := os.ReadFile(filePath)
		if err != nil {
			ss.Status = "missing"
			result.Stages = append(result.Stages, ss)
			continue
		}

		ss.Exists = true
		if idx > highestExistingIdx {
			highestExistingIdx = idx
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
				ss.Gate = fm.Gate
			}
		}

		if ps.hasGate && fmErr == nil {
			if gs, gsErr := parseGateSummary(content); gsErr == nil {
				if ss.Gate == "" {
					if gs.Blocker > 0 || gs.Major > 0 {
						ss.Gate = "fail"
					} else {
						ss.Gate = "pass"
					}
				}
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
			}
		}
	}

	// Gate violation: if stage 8 (impl plan) exists but stage 7 (review) gate != "pass"
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
		result.Stages[reviewIdx].Gate != "pass" {
		result.Stages[reviewIdx].HasGaps = true
		result.Errors++
		result.Findings = append(result.Findings, fmt.Sprintf(
			"stage 8 (%s) exists but stage 7 gate=%s (must be pass)",
			pipelineStages[implIdx].file, result.Stages[reviewIdx].Gate,
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
