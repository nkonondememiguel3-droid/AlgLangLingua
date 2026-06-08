package err_test

import (
	errs "alg/err"
	"strings"
	"testing"
)

func TestEmptyListHasNoErrors(t *testing.T) {
	var dl errs.DiagnosticList
	if dl.HasErrors() {
		t.Error("empty list should have no errors")
	}
	if dl.Count() != 0 {
		t.Errorf("expected count 0, got %d", dl.Count())
	}
}

func TestErrorfIncreasesCount(t *testing.T) {
	var dl errs.DiagnosticList
	dl.Errorf("main.al", 1, 5, "unexpected character %q", '@')
	if !dl.HasErrors() {
		t.Error("expected HasErrors() == true after Errorf")
	}
	if dl.Count() != 1 {
		t.Errorf("expected count 1, got %d", dl.Count())
	}
}

func TestWarnfDoesNotSetHasErrors(t *testing.T) {
	var dl errs.DiagnosticList
	dl.Warnf("main.al", 2, 1, "unused variable %q", "x")
	if dl.HasErrors() {
		t.Error("a warning alone must not set HasErrors")
	}
	if len(dl.Warnings()) != 1 {
		t.Errorf("expected 1 warning, got %d", len(dl.Warnings()))
	}
}

func TestSeverityFiltering(t *testing.T) {
	var dl errs.DiagnosticList
	dl.Errorf("a.al", 1, 1, "error one")
	dl.Warnf("a.al", 2, 1, "warning one")
	dl.Errorf("a.al", 3, 1, "error two")

	if len(dl.Errors()) != 2 {
		t.Errorf("expected 2 errors, got %d", len(dl.Errors()))
	}
	if len(dl.Warnings()) != 1 {
		t.Errorf("expected 1 warning, got %d", len(dl.Warnings()))
	}
	if dl.Count() != 3 {
		t.Errorf("expected total 3, got %d", dl.Count())
	}
}

func TestMerge(t *testing.T) {
	var a, b errs.DiagnosticList
	a.Errorf("a.al", 1, 1, "from a")
	b.Errorf("b.al", 2, 2, "from b")
	b.Warnf("b.al", 3, 1, "warn from b")

	a.Merge(&b)

	if a.Count() != 3 {
		t.Errorf("expected 3 after merge, got %d", a.Count())
	}
	if len(a.Errors()) != 2 {
		t.Errorf("expected 2 errors after merge, got %d", len(a.Errors()))
	}
}

func TestFormatContainsAllMessages(t *testing.T) {
	var dl errs.DiagnosticList
	dl.Errorf("f.al", 1, 5, "bad char")
	dl.Warnf("f.al", 2, 1, "unused var")

	out := dl.Format()
	if !strings.Contains(out, "bad char") {
		t.Errorf("Format missing error message; got:\n%s", out)
	}
	if !strings.Contains(out, "unused var") {
		t.Errorf("Format missing warning message; got:\n%s", out)
	}
	// each diagnostic on its own line
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines in Format output, got %d", len(lines))
	}
}

func TestDiagnosticStringFormat(t *testing.T) {
	d := errs.Diagnostic{
		Severity: errs.Error,
		File:     "test.al",
		Line:     10,
		Column:   4,
		Message:  "unexpected character '@'",
	}
	s := d.String()
	// must contain file, line, column, severity, message
	for _, want := range []string{"test.al", "10", "4", "error", "@"} {
		if !strings.Contains(s, want) {
			t.Errorf("diagnostic string %q missing %q", s, want)
		}
	}
}

func TestAllReturnsInsertionOrder(t *testing.T) {
	var dl errs.DiagnosticList
	dl.Errorf("f.al", 3, 1, "third")
	dl.Errorf("f.al", 1, 1, "first")
	dl.Errorf("f.al", 2, 1, "second")

	all := dl.All()
	if all[0].Message != "third" || all[1].Message != "first" || all[2].Message != "second" {
		t.Errorf("All() must preserve insertion order, got: %v", all)
	}
}
