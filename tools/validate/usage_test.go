package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestHelpRequested(t *testing.T) {
	cases := []struct {
		name string
		argv []string
		want bool
	}{
		{"no args", nil, false},
		{"short flag", []string{"-h"}, true},
		{"single dash long", []string{"-help"}, true},
		{"double dash long", []string{"--help"}, true},
		{"help subcommand", []string{"help"}, true},
		{"help after subcommand", []string{"ears", "EP-009", "-h"}, true},
		{"json flag only", []string{"--json"}, false},
		{"epic id", []string{"EP-009"}, false},
		{"subcommand and epic", []string{"req", "EP-009", "--json"}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := helpRequested(tc.argv); got != tc.want {
				t.Fatalf("helpRequested(%q) = %v, want %v", tc.argv, got, tc.want)
			}
		})
	}
}

func TestPrintUsageToWritesFullUsageToGivenWriter(t *testing.T) {
	var buf bytes.Buffer

	printUsageTo(&buf)

	out := buf.String()
	for _, want := range []string{
		"Usage: validate [subcommand] [EP-XXX] [--json]",
		"Subcommands (single check):",
		"validate req EP-009 --json  # JSON for explicit subcommand",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("usage output missing %q; got:\n%s", want, out)
		}
	}
}
