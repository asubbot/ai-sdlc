package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// SeverityCounts holds counts for review severity levels parsed from YAML front matter.
type SeverityCounts struct {
	Blocker int `json:"blocker"`
	Major   int `json:"major"`
	Medium  int `json:"medium"`
	Minor   int `json:"minor"`
}

// FrontMatter represents YAML front matter parsed from artefact markdown files.
type FrontMatter struct {
	Artefact          string          `yaml:"artefact"`
	EpicID            string          `yaml:"epic_id"`
	Status            string          `yaml:"status"`
	SourceOfTruth     bool            `yaml:"source_of_truth"`
	UpdatedAt         string          `yaml:"updated_at"`
	Gate              string          `yaml:"gate,omitempty"`
	LatestIteration   int             `yaml:"latest_iteration,omitempty"`
	NextAction        string          `yaml:"next_action,omitempty"`
	OpenCounts        *SeverityCounts `yaml:"open_counts,omitempty"`
	NonBlockingCounts *SeverityCounts `yaml:"non_blocking_counts,omitempty"`
}

// StructureFinding is a single validation finding for an artefact file.
type StructureFinding struct {
	File     string `json:"file"`
	Check    string `json:"check"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

// StructureResult aggregates all findings for an epic's artefact structure.
type StructureResult struct {
	Epic     string             `json:"epic"`
	Findings []StructureFinding `json:"findings"`
	Errors   int                `json:"errors"`
	Warnings int                `json:"warnings"`
	HasGaps  bool               `json:"has_gaps"`
}

var requiredSections = map[string][]string{
	"ep-scope":                {"Glossary", "Scope", "Success criteria", "Traceability"},
	"ep-requirements":         {"Introduction", "Glossary", "Requirements"},
	"ep-acceptance-criteria":  {"Scenarios"},
	"ep-system-design":        {"Overview", "Components", "Traceability"},
	"ep-implementation-plan":  {"Tasks"},
	"ep-context":              {"Purpose", "Current Scope", "Open Questions", "Links"},
	"ep-system-design-review": {"Current Gate Summary", "Review iteration"},
	"ep-code-review":          {"Current Gate Summary", "Review iteration"},
	"ep-audit-report":         {"Summary", "Implementation vs plan"},
}

var dateFormatRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

var mdLinkRe = regexp.MustCompile(`\[([^\]]*)\]\(([^)]+)\)`)

func parseFrontMatter(content string) (*FrontMatter, error) {
	lines := strings.Split(content, "\n")
	if len(lines) < 2 || strings.TrimSpace(lines[0]) != "---" {
		return nil, fmt.Errorf("no YAML front matter found")
	}
	endIdx := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			endIdx = i
			break
		}
	}
	if endIdx < 0 {
		return nil, fmt.Errorf("unclosed YAML front matter")
	}
	fm := &FrontMatter{}
	var parentKey string
	for _, line := range lines[1:endIdx] {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		isIndented := len(line) > 0 && (line[0] == ' ' || line[0] == '\t')
		parts := strings.SplitN(trimmed, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		if isIndented && parentKey != "" {
			parseFrontMatterSubKey(fm, parentKey, key, val)
			continue
		}

		parentKey = ""
		switch key {
		case "artefact":
			fm.Artefact = val
		case "epic_id":
			fm.EpicID = val
		case "status":
			fm.Status = val
		case "source_of_truth":
			fm.SourceOfTruth = val == "true"
		case "updated_at":
			fm.UpdatedAt = val
		case "gate":
			fm.Gate = val
		case "latest_iteration":
			if n, err := strconv.Atoi(val); err == nil {
				fm.LatestIteration = n
			}
		case "next_action":
			fm.NextAction = val
		case "open_counts":
			if val == "" {
				fm.OpenCounts = &SeverityCounts{}
				parentKey = "open_counts"
			}
		case "non_blocking_counts":
			if val == "" {
				fm.NonBlockingCounts = &SeverityCounts{}
				parentKey = "non_blocking_counts"
			}
		}
	}
	return fm, nil
}

func parseFrontMatterSubKey(fm *FrontMatter, parent, key, val string) {
	n, err := strconv.Atoi(val)
	if err != nil {
		return
	}
	var target *SeverityCounts
	switch parent {
	case "open_counts":
		target = fm.OpenCounts
	case "non_blocking_counts":
		target = fm.NonBlockingCounts
	}
	if target == nil {
		return
	}
	switch key {
	case "blocker":
		target.Blocker = n
	case "major":
		target.Major = n
	case "medium":
		target.Medium = n
	case "minor":
		target.Minor = n
	}
}

func isValidDateFormat(d string) bool {
	return dateFormatRe.MatchString(d)
}

func validateFrontMatter(fm *FrontMatter, epicID, filename string) []StructureFinding {
	var findings []StructureFinding
	if fm.Artefact == "" {
		findings = append(findings, StructureFinding{
			File: filename, Check: "front_matter", Severity: "error",
			Message: "missing required field: artefact",
		})
	}
	if fm.EpicID == "" {
		findings = append(findings, StructureFinding{
			File: filename, Check: "front_matter", Severity: "error",
			Message: "missing required field: epic_id",
		})
	}
	if fm.Status == "" {
		findings = append(findings, StructureFinding{
			File: filename, Check: "front_matter", Severity: "error",
			Message: "missing required field: status",
		})
	}
	if fm.UpdatedAt == "" {
		findings = append(findings, StructureFinding{
			File: filename, Check: "front_matter", Severity: "error",
			Message: "missing required field: updated_at",
		})
	}
	if fm.EpicID != "" && fm.EpicID != epicID {
		findings = append(findings, StructureFinding{
			File: filename, Check: "front_matter", Severity: "error",
			Message: fmt.Sprintf("epic_id %q does not match directory %q", fm.EpicID, epicID),
		})
	}
	if fm.UpdatedAt != "" && !isValidDateFormat(fm.UpdatedAt) {
		findings = append(findings, StructureFinding{
			File: filename, Check: "front_matter", Severity: "error",
			Message: fmt.Sprintf("updated_at %q is not in YYYY-MM-DD format", fm.UpdatedAt),
		})
	}
	if !fm.SourceOfTruth && fm.Artefact != "ep-context" {
		findings = append(findings, StructureFinding{
			File: filename, Check: "front_matter", Severity: "warning",
			Message: "source_of_truth is not set to true",
		})
	}
	return findings
}

func findRequiredSections(content, artefactType string) []StructureFinding {
	sections, ok := requiredSections[artefactType]
	if !ok {
		return nil
	}

	contentLower := strings.ToLower(content)
	var findings []StructureFinding
	for _, section := range sections {
		sectionLower := strings.ToLower(section)
		found := false
		for _, line := range strings.Split(contentLower, "\n") {
			trimmed := strings.TrimSpace(line)
			if (strings.HasPrefix(trimmed, "## ") || strings.HasPrefix(trimmed, "### ")) &&
				strings.Contains(trimmed, sectionLower) {
				found = true
				break
			}
		}
		if !found {
			findings = append(findings, StructureFinding{
				File: artefactType + ".md", Check: "required_section", Severity: "error",
				Message: fmt.Sprintf("missing required section: %s", section),
			})
		}
	}
	return findings
}

func findBrokenLinks(content, baseDir string) []StructureFinding {
	var findings []StructureFinding
	matches := mdLinkRe.FindAllStringSubmatch(content, -1)
	for _, m := range matches {
		linkTarget := m[2]
		if strings.HasPrefix(linkTarget, "http://") || strings.HasPrefix(linkTarget, "https://") {
			continue
		}
		if strings.HasPrefix(linkTarget, "#") {
			continue
		}
		// Strip anchor from file path
		path := linkTarget
		if idx := strings.Index(path, "#"); idx >= 0 {
			path = path[:idx]
		}
		if path == "" {
			continue
		}
		fullPath := filepath.Join(baseDir, path)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			findings = append(findings, StructureFinding{
				File: filepath.Base(baseDir), Check: "broken_link", Severity: "warning",
				Message: fmt.Sprintf("broken link: %s", linkTarget),
			})
		}
	}
	return findings
}

var expectedArtefacts = []string{
	"ep-scope.md",
	"ep-context.md",
	"ep-requirements.md",
	"ep-acceptance-criteria.md",
	"ep-system-design.md",
	"ep-system-design-review.md",
	"ep-implementation-plan.md",
	"ep-code-review.md",
	"ep-audit-report.md",
}

func validateArtefactStructure(epicDir, epicID string) *StructureResult {
	result := &StructureResult{
		Epic:     epicID,
		Findings: []StructureFinding{},
	}

	for _, filename := range expectedArtefacts {
		filePath := filepath.Join(epicDir, filename)
		data, err := os.ReadFile(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				result.Findings = append(result.Findings, StructureFinding{
					File: filename, Check: "file_exists", Severity: "warning",
					Message: fmt.Sprintf("artefact file not found: %s", filename),
				})
				result.Warnings++
			}
			continue
		}

		content := string(data)

		fm, fmErr := parseFrontMatter(content)
		if fmErr != nil {
			result.Findings = append(result.Findings, StructureFinding{
				File: filename, Check: "front_matter", Severity: "error",
				Message: fmErr.Error(),
			})
			result.Errors++
			continue
		}

		fmFindings := validateFrontMatter(fm, epicID, filename)
		for _, f := range fmFindings {
			if f.Severity == "error" {
				result.Errors++
			} else {
				result.Warnings++
			}
		}
		result.Findings = append(result.Findings, fmFindings...)

		artefactType := strings.TrimSuffix(filename, ".md")
		secFindings := findRequiredSections(content, artefactType)
		for _, f := range secFindings {
			f.File = filename
			if f.Severity == "error" {
				result.Errors++
			} else {
				result.Warnings++
			}
			result.Findings = append(result.Findings, f)
		}

		linkFindings := findBrokenLinks(content, epicDir)
		for i := range linkFindings {
			linkFindings[i].File = filename
		}
		for _, f := range linkFindings {
			if f.Severity == "error" {
				result.Errors++
			} else {
				result.Warnings++
			}
		}
		result.Findings = append(result.Findings, linkFindings...)
	}

	result.HasGaps = result.Errors > 0
	return result
}

func printStructureHuman(result *StructureResult) {
	writeStdout("\n📋 Structure Validation for %s\n\n", result.Epic)

	if len(result.Findings) == 0 {
		writelnStdout("✅ No issues found")
		writelnStdout("")
		return
	}

	currentFile := ""
	for _, f := range result.Findings {
		if f.File != currentFile {
			currentFile = f.File
			writeStdout("\n  %s\n", currentFile)
		}
		marker := "⚠️"
		if f.Severity == "error" {
			marker = "❌"
		}
		writeStdout("    %s [%s] %s\n", marker, f.Check, f.Message)
	}

	writelnStdout("")
	statusEmoji := "✅"
	if result.HasGaps {
		statusEmoji = "❌"
	}
	writeStdout("%s RESULT: %d error(s), %d warning(s)\n\n", statusEmoji, result.Errors, result.Warnings)
}

func runStructureValidation(epic string, jsonOut bool) {
	if epic == "all" {
		fmt.Fprintf(os.Stderr, "structure subcommand requires an epic ID\n")
		os.Exit(1)
	}
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}
	epicDir := filepath.Join(cwd, "ai-sdlc-artefacts", "epics", epic)
	result := validateArtefactStructure(epicDir, epic)
	if jsonOut {
		data, _ := json.MarshalIndent(result, "", "  ")
		writelnStdout(string(data))
	} else {
		printStructureHuman(result)
	}
	if result.HasGaps {
		os.Exit(1)
	}
}
