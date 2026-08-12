package main

import (
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"
)

const testMainName = "TestMain"

// findTestsMissingACTrace walks tests/, internal/, and cmd/ like findCoverageInCodebase and
// returns sorted "rel/path/to/file_test.go::TestName" entries for top-level Test functions that
// have no AC trace comment bound to them, using the same shared AST index/binding rules as
// coverage (lineDeclaresACCoverage, extractACsFromLine non-empty, parseTestASTIndex/bindTraceLine).
// TestMain is excluded.
func findTestsMissingACTrace(codebasePath string) ([]string, error) {
	root, err := os.OpenRoot(codebasePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = root.Close() }()
	fsys := root.FS()

	var out []string
	for _, dir := range []string{"tests", "internal", "cmd"} {
		if _, err := root.Stat(dir); os.IsNotExist(err) {
			continue
		} else if err != nil {
			return nil, err
		}
		err = fs.WalkDir(fsys, dir, func(path string, d fs.DirEntry, errWalk error) error {
			if errWalk != nil {
				return fmt.Errorf("walk %s: %w", path, errWalk)
			}
			if d.IsDir() {
				return nil
			}
			if !strings.HasSuffix(path, "_test.go") {
				return nil
			}
			content, errRead := root.ReadFile(path)
			if errRead != nil {
				return fmt.Errorf("read test file %s: %w", path, errRead)
			}
			missing, err := testFuncsMissingACTraceInFile(path, string(content))
			if err != nil {
				return err
			}
			out = append(out, missing...)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	sort.Strings(out)
	return out, nil
}

// testFuncsMissingACTraceInFile returns refs relPath::TestName for Test* functions without an AC trace.
func testFuncsMissingACTraceInFile(relPath, fileContent string) ([]string, error) {
	index, err := parseTestASTIndex([]byte(fileContent), relPath)
	if err != nil {
		return nil, fmt.Errorf("parse test file %s: %w", relPath, err)
	}

	lines := strings.Split(fileContent, "\n")
	testNames := index.topLevelTestNames()
	if len(testNames) == 0 {
		return nil, nil
	}
	traced := make(map[string]struct{})
	for i, line := range lines {
		if !index.isActualCommentLine(i + 1) {
			continue
		}
		if !lineDeclaresACCoverage(line) || len(extractACsFromLine(line)) == 0 {
			continue
		}
		name, err := index.bindTraceLine(i + 1)
		if err != nil {
			return nil, err
		}
		traced[name] = struct{}{}
	}
	var out []string
	for _, name := range testNames {
		if name == testMainName {
			continue
		}
		if _, ok := traced[name]; ok {
			continue
		}
		out = append(out, fmt.Sprintf("%s::%s", relPath, name))
	}
	return out, nil
}
