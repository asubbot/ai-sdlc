package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func mergeACExclusion(prev, next acExclusionKind) acExclusionKind {
	if next == acExclusionObsolete || prev == acExclusionObsolete {
		return acExclusionObsolete
	}
	if next == acExclusionDeferred || prev == acExclusionDeferred {
		return acExclusionDeferred
	}
	return acExclusionNone
}

func acLineMentionsTarget(lineUpper, targetUpper string) bool {
	return strings.Contains(lineUpper, targetUpper)
}

func lineImpliesObsolete(lineUpper, targetUpper string) bool {
	if acLineMentionsTarget(lineUpper, targetUpper) && strings.Contains(lineUpper, "OBSOLETE") {
		return true
	}
	return strings.Contains(lineUpper, "STATUS:") && strings.Contains(lineUpper, "OBSOLETE")
}

func lineImpliesDeferred(lineUpper, targetUpper string) bool {
	if !acLineMentionsTarget(lineUpper, targetUpper) {
		return strings.Contains(lineUpper, "STATUS:") && strings.Contains(lineUpper, "DEFERRED")
	}
	if strings.Contains(lineUpper, "DEFERRED") || strings.Contains(lineUpper, "MANUAL ONLY") {
		return true
	}
	return false
}

// detectACExclusion returns whether an AC is excluded from mandatory automated traces
// because the epic marks it **Deferred:** / **Obsolete:** / MANUAL ONLY, or **Status:** … DEFERRED/OBSOLETE,
// within a few lines of the AC mention (same heuristic as historical "deferred" parsing).
func detectACExclusion(lines []string, idx int, code string) acExclusionKind {
	targetUpper := strings.ToUpper(code)
	start := idx - 3
	if start < 0 {
		start = 0
	}
	end := idx + 6
	if end >= len(lines) {
		end = len(lines) - 1
	}
	obsolete := false
	deferred := false
	for i := start; i <= end; i++ {
		lineUpper := strings.ToUpper(lines[i])
		if lineImpliesObsolete(lineUpper, targetUpper) {
			obsolete = true
		}
		if lineImpliesDeferred(lineUpper, targetUpper) {
			deferred = true
		}
	}
	if obsolete {
		return acExclusionObsolete
	}
	if deferred {
		return acExclusionDeferred
	}
	return acExclusionNone
}

// parseACsFromFile extracts all AC-EE.NNN codes from an acceptance criteria markdown file.
// The second map marks ACs excluded from automated test requirements (**Deferred:**, **Obsolete:**, etc.).
func parseACsFromFile(path string) (map[ACCode]string, map[ACCode]acExclusionKind, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	acs := make(map[ACCode]string)
	excluded := make(map[ACCode]acExclusionKind)
	lines := strings.Split(string(content), "\n")

	acCodePattern := regexp.MustCompile(`AC-(\d{2,3})\.(\d{3})`)

	for i, line := range lines {
		matches := acCodePattern.FindAllStringSubmatch(line, -1)
		if len(matches) == 0 {
			continue
		}

		for _, match := range matches {
			code := fmt.Sprintf("AC-%s.%s", match[1], match[2])
			ac := ACCode(code)

			if _, exists := acs[ac]; !exists {
				criterion := strings.TrimSpace(acCodePattern.ReplaceAllString(line, ""))
				criterion = strings.Trim(criterion, "*[]()#:- \t")

				if criterion == "" && i+1 < len(lines) {
					criterion = strings.TrimSpace(lines[i+1])
				}
				acs[ac] = criterion
			}

			if k := detectACExclusion(lines, i, code); k != acExclusionNone {
				excluded[ac] = mergeACExclusion(excluded[ac], k)
			}
		}
	}

	return acs, excluded, nil
}

// parseREQCountFromFile extracts a count of unique REQ-EE.NNN codes from ep-requirements.md.
func parseREQCountFromFile(path string) (int, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	re := regexp.MustCompile(`REQ-(\d{2,3})\.(\d{3})`)
	matches := re.FindAllStringSubmatch(string(content), -1)
	seen := make(map[string]struct{}, len(matches))
	for _, m := range matches {
		seen["REQ-"+m[1]+"."+m[2]] = struct{}{}
	}
	return len(seen), nil
}

// normalizeCriterionPreview drops table/index lines from ep-acceptance-criteria.md that are not useful as a one-line hint.
func normalizeCriterionPreview(c string) string {
	c = strings.TrimSpace(c)
	if c == "" {
		return ""
	}
	if strings.HasPrefix(c, "|") {
		return ""
	}
	if strings.Count(c, "|") >= 3 {
		return ""
	}
	if len(c) > 160 {
		return c[:157] + "..."
	}
	return c
}

// Regexes for common traceability comment shapes (see project test style).
var (
	reEpicACLine = regexp.MustCompile(`EP-\d+\s+AC-\d{2,3}\.\d{3}`)
	reACLabel    = regexp.MustCompile(`\bAC-\d{2,3}\.\d{3}\s*:`)
	reACSlashReq = regexp.MustCompile(`\bAC-\d{2,3}\.\d{3}\s*/\s*REQ-`)
	reACCode     = regexp.MustCompile(`AC-\d{2,3}\.\d{3}`)
	reManualWord = regexp.MustCompile(`(?i)\bmanual\b`)
	reSingleAC   = regexp.MustCompile(`AC-(\d{2,3})\.(\d{3})`)
	reACRange    = regexp.MustCompile(`AC-(\d{2,3})\.(\d{3})[–-](\d{3})`)
)

// lineDeclaresManualTrace is true when the traceability line explicitly marks manual-only intent.
func lineDeclaresManualTrace(line string) bool {
	return reManualWord.MatchString(line)
}

// lineDeclaresACCoverage returns true if a line in a *_test.go file should be scanned for AC codes.
func lineDeclaresACCoverage(line string) bool {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "//") {
		return false
	}
	lineLower := strings.ToLower(line)
	if strings.Contains(lineLower, "covers") || strings.Contains(lineLower, "supporting") {
		return true
	}
	if reEpicACLine.MatchString(line) {
		return true
	}
	if reACLabel.MatchString(line) {
		return true
	}
	if reACSlashReq.MatchString(line) {
		return true
	}
	if strings.Contains(line, "REQ-") && reACCode.MatchString(line) {
		return true
	}
	return false
}

// extractACsFromLine extracts all AC codes from a traceability line.
func extractACsFromLine(line string) []ACCode {
	var result []ACCode

	matches := reSingleAC.FindAllStringSubmatch(line, -1)
	acMap := make(map[ACCode]bool)

	for _, match := range matches {
		epic := match[1]
		num := match[2]
		code := fmt.Sprintf("AC-%s.%s", epic, num)
		acMap[ACCode(code)] = true
	}

	rangeMatches := reACRange.FindAllStringSubmatch(line, -1)

	for _, match := range rangeMatches {
		epic := match[1]
		startStr := match[2]
		endStr := match[3]

		start, err1 := strconv.Atoi(startStr)
		end, err2 := strconv.Atoi(endStr)
		if err1 != nil || err2 != nil {
			continue
		}

		for i := start; i <= end; i++ {
			code := fmt.Sprintf("AC-%s.%03d", epic, i)
			acMap[ACCode(code)] = true
		}
	}

	for ac := range acMap {
		result = append(result, ac)
	}
	sort.Slice(result, func(i, j int) bool {
		return string(result[i]) < string(result[j])
	})

	return result
}
