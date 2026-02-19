package output

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

// DateTimeFormat is the standard human-friendly date format for CLI output.
const DateTimeFormat = "Jan 2, 2006 at 15:04"

var (
	stdOut io.Writer = os.Stdout
	stdErr io.Writer = os.Stderr
)

// SetWriters overrides the output writers (useful for tests).
func SetWriters(stdout, stderr io.Writer) {
	if stdout != nil {
		stdOut = stdout
	}
	if stderr != nil {
		stdErr = stderr
	}
}

// Writer returns the shared stdout writer for user-facing output.
func Writer() io.Writer {
	return stdOut
}

// ErrorWriter returns the stderr writer for error output.
func ErrorWriter() io.Writer {
	return stdErr
}

// Printf writes formatted output to stdout.
func Printf(format string, args ...any) {
	fmt.Fprintf(stdOut, format, args...)
}

// Println writes a line to stdout.
func Println(args ...any) {
	fmt.Fprintln(stdOut, args...)
}

// Blank writes a single blank line to stdout.
func Blank() {
	fmt.Fprintln(stdOut)
}

// NewTableWriter returns a tabwriter for stdout with consistent spacing.
func NewTableWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(stdOut, 0, 0, 2, ' ', 0)
}

// FormatVersion normalizes a version string with a single leading 'v'.
func FormatVersion(version string) string {
	trimmed := strings.TrimSpace(version)
	trimmed = strings.TrimPrefix(trimmed, "v")
	if trimmed == "" {
		return "v0.0.0"
	}
	return "v" + trimmed
}
