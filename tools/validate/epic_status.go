package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var epicStatusRowPattern = regexp.MustCompile(`(?i)\|\s*\*\*Status\*\*\s*\|\s*([^|]+)\|`)

// parseEpicStatusFromScope extracts the Status table cell from ep-scope.md body text.
func parseEpicStatusFromScope(content string) (string, bool) {
	for _, line := range strings.Split(content, "\n") {
		if matches := epicStatusRowPattern.FindStringSubmatch(line); len(matches) > 1 {
			return strings.TrimSpace(matches[1]), true
		}
	}
	return "", false
}

// isSkippedEpicStatus reports whether ears/req all should skip this epic (NEW or CANCEL*).
func isSkippedEpicStatus(status string) bool {
	upper := strings.ToUpper(strings.TrimSpace(status))
	if upper == "NEW" {
		return true
	}
	return strings.HasPrefix(upper, "CANCEL")
}

func readEpicScopeStatus(epicDir string) (string, error) {
	scopePath := filepath.Join(epicDir, "ep-scope.md")
	content, err := os.ReadFile(scopePath)
	if err != nil {
		return "", fmt.Errorf("read ep-scope.md: %w", err)
	}
	status, ok := parseEpicStatusFromScope(string(content))
	if !ok {
		// Fixtures and legacy scopes without a Status row are treated as in-scope.
		return "DONE", nil
	}
	return status, nil
}

// readInScopeEpicNames returns epic IDs that require ears/req validation and those skipped by status.
func readInScopeEpicNames(epicsPath string) (inScope, skipped []string, err error) {
	epics, err := readSortedEpicNames(epicsPath)
	if err != nil {
		return nil, nil, err
	}
	for _, epic := range epics {
		epicDir := filepath.Join(epicsPath, epic)
		scopePath := filepath.Join(epicDir, "ep-scope.md")
		if _, statErr := os.Stat(scopePath); os.IsNotExist(statErr) {
			return nil, nil, fmt.Errorf("%s: ep-scope.md not found", epic)
		}
		status, readErr := readEpicScopeStatus(epicDir)
		if readErr != nil {
			return nil, nil, fmt.Errorf("%s: %w", epic, readErr)
		}
		if isSkippedEpicStatus(status) {
			skipped = append(skipped, epic)
			continue
		}
		inScope = append(inScope, epic)
	}
	return inScope, skipped, nil
}
