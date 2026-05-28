package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// EARSFinding is a single issue found by the EARS linter.
type EARSFinding struct {
	REQ      string `json:"req"`
	Line     int    `json:"line"`
	Severity string `json:"severity"` // "error" | "warning"
	Rule     string `json:"rule"`
	Message  string `json:"message"`
}

// EARSResult aggregates all findings for one requirements file.
type EARSResult struct {
	Epic     string        `json:"epic"`
	Total    int           `json:"total_reqs"`
	Errors   int           `json:"errors"`
	Warnings int           `json:"warnings"`
	Findings []EARSFinding `json:"findings"`
	HasGaps  bool          `json:"has_gaps"`
}

type reqBlock struct {
	Code    string
	Summary string
	Body    string
	Line    int
}

var reqHeadingRe = regexp.MustCompile(`^###\s+REQ-(\d{2,3})\.(\d{3})\s*[—–-]\s*(.*)$`)

var (
	reEARSShall = regexp.MustCompile(`(?i)\bSHALL\b`)
	reEARSWhen  = regexp.MustCompile(`(?i)\bWHEN\b`)
	reEARSWhile = regexp.MustCompile(`(?i)\bWHILE\b`)
	reEARSIf    = regexp.MustCompile(`(?i)\bIF\b`)
	reEARSThen  = regexp.MustCompile(`(?i)\bTHEN\b`)
	reEARSWhere = regexp.MustCompile(`(?i)\bWHERE\b`)
	reEARSThe   = regexp.MustCompile(`(?i)\bTHE\b`)
)

var bannedWords = []string{
	"efficiently",
	"appropriately",
	"reasonable",
	"if possible",
	"adequate",
	"etc.",
	"and/or",
	"as needed",
}

var (
	reShouldWord       = regexp.MustCompile(`(?i)\bshould\b`)
	rePassiveVoice     = regexp.MustCompile(`(?i)\b(is|are|was|were|been)\s+\w+ed\b`)
	reMultipleThoughts = regexp.MustCompile(`(?i)\bSHALL\b.+?\band\b.+?\bSHALL\b`)
	reDoubleQuoted     = regexp.MustCompile(`"[^"]*"`)
	reSingleQuoted     = regexp.MustCompile(`'[^']*'`)
)

func parseREQBlocks(content string) []reqBlock {
	lines := strings.Split(content, "\n")
	var blocks []reqBlock
	var current *reqBlock
	var bodyLines []string

	for i, line := range lines {
		if m := reqHeadingRe.FindStringSubmatch(line); len(m) > 0 {
			if current != nil {
				current.Body = strings.TrimSpace(strings.Join(bodyLines, "\n"))
				blocks = append(blocks, *current)
			}
			current = &reqBlock{
				Code:    fmt.Sprintf("REQ-%s.%s", m[1], m[2]),
				Summary: strings.TrimSpace(m[3]),
				Line:    i + 1,
			}
			bodyLines = nil
			continue
		}
		if current != nil {
			if strings.HasPrefix(line, "### ") {
				current.Body = strings.TrimSpace(strings.Join(bodyLines, "\n"))
				blocks = append(blocks, *current)
				current = nil
				bodyLines = nil
				continue
			}
			bodyLines = append(bodyLines, line)
		}
	}
	if current != nil {
		current.Body = strings.TrimSpace(strings.Join(bodyLines, "\n"))
		blocks = append(blocks, *current)
	}
	return blocks
}

func checkEARSPattern(block reqBlock) []EARSFinding {
	text := strings.Join(strings.Fields(block.Body), " ")

	if !reEARSShall.MatchString(text) {
		return []EARSFinding{{
			REQ:      block.Code,
			Line:     block.Line,
			Severity: "error",
			Rule:     "ears-missing-shall",
			Message:  "requirement does not use SHALL",
		}}
	}

	hasWhen := reEARSWhen.MatchString(text)
	hasWhile := reEARSWhile.MatchString(text)
	hasIf := reEARSIf.MatchString(text)
	hasThen := reEARSThen.MatchString(text)
	hasWhere := reEARSWhere.MatchString(text)
	hasThe := reEARSThe.MatchString(text)

	switch {
	case hasWhile && hasWhen && hasThe:
		return nil
	case hasIf && hasThen && hasThe:
		return nil
	case hasWhen && hasThe:
		return nil
	case hasWhile && hasThe:
		return nil
	case hasWhere && hasThe:
		return nil
	case hasThe:
		return nil
	}

	return []EARSFinding{{
		REQ:      block.Code,
		Line:     block.Line,
		Severity: "error",
		Rule:     "ears-no-pattern",
		Message:  "does not match any EARS pattern",
	}}
}

func stripQuotedText(s string) string {
	s = reDoubleQuoted.ReplaceAllString(s, "")
	s = reSingleQuoted.ReplaceAllString(s, "")
	return s
}

func checkBannedWords(block reqBlock) []EARSFinding {
	lower := strings.ToLower(block.Body)
	var findings []EARSFinding

	for _, word := range bannedWords {
		if strings.Contains(lower, word) {
			findings = append(findings, EARSFinding{
				REQ:      block.Code,
				Line:     block.Line,
				Severity: "error",
				Rule:     "ears-banned-word",
				Message:  fmt.Sprintf("banned word: %q", word),
			})
		}
	}

	if reShouldWord.MatchString(block.Body) {
		stripped := stripQuotedText(block.Body)
		if reShouldWord.MatchString(stripped) {
			findings = append(findings, EARSFinding{
				REQ:      block.Code,
				Line:     block.Line,
				Severity: "error",
				Rule:     "ears-banned-word",
				Message:  "use SHALL instead of should",
			})
		}
	}

	return findings
}

func checkPassiveVoice(block reqBlock) []EARSFinding {
	matches := rePassiveVoice.FindAllString(block.Body, -1)
	var findings []EARSFinding
	for _, m := range matches {
		findings = append(findings, EARSFinding{
			REQ:      block.Code,
			Line:     block.Line,
			Severity: "warning",
			Rule:     "ears-passive-voice",
			Message:  fmt.Sprintf("possible passive voice: %q", m),
		})
	}
	return findings
}

func checkMultipleThoughts(block reqBlock) []EARSFinding {
	text := strings.Join(strings.Fields(block.Body), " ")
	if reMultipleThoughts.MatchString(text) {
		return []EARSFinding{{
			REQ:      block.Code,
			Line:     block.Line,
			Severity: "warning",
			Rule:     "ears-multiple-thoughts",
			Message:  "requirement contains multiple SHALL clauses joined by 'and'",
		}}
	}
	return nil
}

func lintEARSFile(path string) *EARSResult {
	content, err := os.ReadFile(path)
	if err != nil {
		return &EARSResult{
			Findings: []EARSFinding{{
				Severity: "error",
				Rule:     "ears-file-error",
				Message:  fmt.Sprintf("cannot read file: %v", err),
			}},
			Errors:  1,
			HasGaps: true,
		}
	}

	blocks := parseREQBlocks(string(content))
	result := &EARSResult{
		Total:    len(blocks),
		Findings: []EARSFinding{},
	}

	for _, block := range blocks {
		var findings []EARSFinding
		findings = append(findings, checkEARSPattern(block)...)
		findings = append(findings, checkBannedWords(block)...)
		findings = append(findings, checkPassiveVoice(block)...)
		findings = append(findings, checkMultipleThoughts(block)...)

		for _, f := range findings {
			if f.Severity == "error" {
				result.Errors++
			} else {
				result.Warnings++
			}
		}
		result.Findings = append(result.Findings, findings...)
	}

	result.HasGaps = result.Errors > 0
	return result
}

func printEARSHuman(result *EARSResult) {
	writeStdout("\n📋 EARS Lint — %s\n\n", result.Epic)

	if len(result.Findings) == 0 {
		writelnStdout("✅ All requirements follow EARS patterns — no issues found.")
		writeStdout("\nTotal: %d requirements\n\n", result.Total)
		return
	}

	for _, f := range result.Findings {
		icon := "❌"
		if f.Severity == "warning" {
			icon = "⚠️"
		}
		writeStdout("%s %s (line %d) [%s]: %s\n", icon, f.REQ, f.Line, f.Rule, f.Message)
	}

	writelnStdout("")
	writelnStdout(strings.Repeat("─", 80))

	statusIcon := "✅"
	if result.HasGaps {
		statusIcon = "❌"
	}

	writeStdout("\n%s Total: %d requirements, %d errors, %d warnings\n\n", statusIcon, result.Total, result.Errors, result.Warnings)
}
