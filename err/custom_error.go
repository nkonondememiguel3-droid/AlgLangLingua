package err

import (
	"fmt"
	"strings"
)

// Severity classifies how serious a diagnostic is.
type Severity int

const (
	Warning Severity = iota
	Error
)

func (s Severity) String() string {
	switch s {
	case Warning:
		return "warning"
	case Error:
		return "error"
	default:
		return "unknown"
	}
}

// Diagnostic is a single compiler message tied to a source location.
type Diagnostic struct {
	Severity Severity
	File     string
	Line     int
	Column   int
	Message  string
}

func (d Diagnostic) String() string {
	loc := ""
	if d.File != "" {
		loc = fmt.Sprintf("%s:", d.File)
	}
	if d.Line > 0 {
		loc += fmt.Sprintf("%d:%d: ", d.Line, d.Column)
	}
	return fmt.Sprintf("%s%s: %s", loc, d.Severity, d.Message)
}

// DiagnosticList accumulates diagnostics across a compilation phase.
// The zero value is ready to use.
type DiagnosticList struct {
	items []Diagnostic
}

// Add appends a diagnostic.
func (dl *DiagnosticList) Add(d Diagnostic) {
	dl.items = append(dl.items, d)
}

// Errorf appends an error-severity diagnostic at the given location.
func (dl *DiagnosticList) Errorf(file string, line, col int, format string, args ...any) {
	dl.Add(Diagnostic{
		Severity: Error,
		File:     file,
		Line:     line,
		Column:   col,
		Message:  fmt.Sprintf(format, args...),
	})
}

// Warnf appends a warning-severity diagnostic.
func (dl *DiagnosticList) Warnf(file string, line, col int, format string, args ...any) {
	dl.Add(Diagnostic{
		Severity: Warning,
		File:     file,
		Line:     line,
		Column:   col,
		Message:  fmt.Sprintf(format, args...),
	})
}

// HasErrors reports whether any error-severity diagnostic was recorded.
func (dl *DiagnosticList) HasErrors() bool {
	for _, d := range dl.items {
		if d.Severity == Error {
			return true
		}
	}
	return false
}

// All returns all diagnostics in insertion order.
func (dl *DiagnosticList) All() []Diagnostic {
	return dl.items
}

// Errors returns only error-severity diagnostics.
func (dl *DiagnosticList) Errors() []Diagnostic {
	var out []Diagnostic
	for _, d := range dl.items {
		if d.Severity == Error {
			out = append(out, d)
		}
	}
	return out
}

// Warnings returns only warning-severity diagnostics.
func (dl *DiagnosticList) Warnings() []Diagnostic {
	var out []Diagnostic
	for _, d := range dl.items {
		if d.Severity == Warning {
			out = append(out, d)
		}
	}
	return out
}

// Merge appends all diagnostics from another list into this one.
// Used to bubble phase diagnostics up to the top-level reporter.
func (dl *DiagnosticList) Merge(other *DiagnosticList) {
	dl.items = append(dl.items, other.items...)
}

// Format returns all diagnostics as a single human-readable string,
// one per line, sorted by insertion order (which is source order
// as long as the lexer/parser scan left-to-right).
func (dl *DiagnosticList) Format() string {
	if len(dl.items) == 0 {
		return ""
	}
	var sb strings.Builder
	for _, d := range dl.items {
		sb.WriteString(d.String())
		sb.WriteByte('\n')
	}
	return sb.String()
}

// Count returns the total number of diagnostics.
func (dl *DiagnosticList) Count() int {
	return len(dl.items)
}
