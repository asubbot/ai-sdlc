package main

import (
	"fmt"
	"log"
	"os"
)

// stdout is human-oriented CLI output. writeStdout avoids forbidigo (^fmt\.Print*) and satisfies errcheck.
var stdout = os.Stdout

// errLog is the CLI diagnostic log. Zero flags keep messages byte-identical to
// direct stderr writes: no date, no time, no prefix.
var errLog = log.New(os.Stderr, "", 0)

func writeStdout(format string, args ...any) {
	_, _ = fmt.Fprintf(stdout, format, args...)
}

func writelnStdout(args ...any) {
	_, _ = fmt.Fprintln(stdout, args...)
}

// writeStderr emits non-diagnostic text on stderr - usage hints that accompany
// an error rather than describing one.
func writeStderr(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
}
