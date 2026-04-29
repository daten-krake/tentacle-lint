package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/daten-krake/tentacle-lint/internal/linter"
)

func TestTextNoIssues(t *testing.T) {
	var buf bytes.Buffer
	Text(&buf, nil, false, Options{})
	if buf.String() != "No issues found.\n" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestTextError(t *testing.T) {
	issues := []linter.Issue{
		{File: "test.yaml", Field: "severity", Message: "invalid value", Sev: linter.Error},
	}
	var buf bytes.Buffer
	Text(&buf, issues, false, Options{})
	want := "test.yaml: severity [error] invalid value\n"
	if buf.String() != want {
		t.Errorf("got  %q\nwant %q", buf.String(), want)
	}
}

func TestTextWarning(t *testing.T) {
	issues := []linter.Issue{
		{File: "test.yaml", Field: "fp_rate", Message: "empty", Sev: linter.Warning},
	}
	var buf bytes.Buffer
	Text(&buf, issues, false, Options{})
	want := "test.yaml: fp_rate [warning] empty\n"
	if buf.String() != want {
		t.Errorf("got  %q\nwant %q", buf.String(), want)
	}
}

func TestTextStrictPromotesWarning(t *testing.T) {
	issues := []linter.Issue{
		{File: "test.yaml", Field: "fp_rate", Message: "empty", Sev: linter.Warning},
	}
	var buf bytes.Buffer
	Text(&buf, issues, true, Options{})
	want := "test.yaml: fp_rate [error] empty\n"
	if buf.String() != want {
		t.Errorf("got  %q\nwant %q", buf.String(), want)
	}
}

func TestTextEmptyField(t *testing.T) {
	issues := []linter.Issue{
		{File: "test.yaml", Field: "", Message: "parse error", Sev: linter.Error},
	}
	var buf bytes.Buffer
	Text(&buf, issues, false, Options{})
	want := "test.yaml: - [error] parse error\n"
	if buf.String() != want {
		t.Errorf("got  %q\nwant %q", buf.String(), want)
	}
}

func TestTextMultipleIssues(t *testing.T) {
	issues := []linter.Issue{
		{File: "a.yaml", Field: "name", Message: "empty", Sev: linter.Error},
		{File: "a.yaml", Field: "query", Message: "invalid", Sev: linter.Error},
		{File: "b.yaml", Field: "severity", Message: "bad", Sev: linter.Warning},
	}
	var buf bytes.Buffer
	Text(&buf, issues, false, Options{})
	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
}

func TestTextColorEnabled(t *testing.T) {
	issues := []linter.Issue{
		{File: "f.yaml", Field: "x", Message: "err", Sev: linter.Error},
	}
	var buf bytes.Buffer
	Text(&buf, issues, false, Options{Color: true})
	got := buf.String()
	if !strings.Contains(got, "\033[") {
		t.Error("expected ANSI escape codes when Color is true")
	}
}

func TestTextColorDisabled(t *testing.T) {
	issues := []linter.Issue{
		{File: "f.yaml", Field: "x", Message: "err", Sev: linter.Error},
	}
	var buf bytes.Buffer
	Text(&buf, issues, false, Options{Color: false})
	got := buf.String()
	if strings.Contains(got, "\033[") {
		t.Error("expected no ANSI escape codes when Color is false")
	}
}

func TestJSONNoIssues(t *testing.T) {
	var buf bytes.Buffer
	if err := JSON(&buf, nil, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out jsonOutput
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(out.Issues) != 0 {
		t.Errorf("expected 0 issues, got %d", len(out.Issues))
	}
}

func TestJSONError(t *testing.T) {
	issues := []linter.Issue{
		{File: "f.yaml", Field: "severity", Message: "bad", Sev: linter.Error},
	}
	var buf bytes.Buffer
	if err := JSON(&buf, issues, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out jsonOutput
	json.Unmarshal(buf.Bytes(), &out)
	if len(out.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(out.Issues))
	}
	ji := out.Issues[0]
	if ji.File != "f.yaml" || ji.Field != "severity" || ji.Severity != "error" || ji.EffectiveSev != "error" {
		t.Errorf("unexpected issue fields: %+v", ji)
	}
	if ji.Promoted {
		t.Error("error should not be marked as promoted")
	}
	if out.Errors != 1 || out.Warnings != 0 {
		t.Errorf("expected 1 error, 0 warnings, got %d errors, %d warnings", out.Errors, out.Warnings)
	}
}

func TestJSONWarning(t *testing.T) {
	issues := []linter.Issue{
		{File: "f.yaml", Field: "fp_rate", Message: "empty", Sev: linter.Warning},
	}
	var buf bytes.Buffer
	if err := JSON(&buf, issues, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out jsonOutput
	json.Unmarshal(buf.Bytes(), &out)
	ji := out.Issues[0]
	if ji.Severity != "warning" || ji.EffectiveSev != "warning" {
		t.Errorf("expected severity warning, effective warning; got %q / %q", ji.Severity, ji.EffectiveSev)
	}
	if ji.Promoted {
		t.Error("warning should not be promoted in non-strict mode")
	}
	if out.Errors != 0 || out.Warnings != 1 {
		t.Errorf("expected 0 errors, 1 warning, got %d errors, %d warnings", out.Errors, out.Warnings)
	}
}

func TestJSONStrictPromotesWarning(t *testing.T) {
	issues := []linter.Issue{
		{File: "f.yaml", Field: "fp_rate", Message: "empty", Sev: linter.Warning},
	}
	var buf bytes.Buffer
	if err := JSON(&buf, issues, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out jsonOutput
	json.Unmarshal(buf.Bytes(), &out)
	ji := out.Issues[0]
	if ji.Severity != "warning" || ji.EffectiveSev != "error" {
		t.Errorf("expected severity warning, effective error; got %q / %q", ji.Severity, ji.EffectiveSev)
	}
	if !ji.Promoted {
		t.Error("warning should be promoted in strict mode")
	}
	if out.Errors != 1 || out.Warnings != 0 {
		t.Errorf("expected 1 error, 0 warnings, got %d errors, %d warnings", out.Errors, out.Warnings)
	}
}

func TestJSONMultipleIssues(t *testing.T) {
	issues := []linter.Issue{
		{File: "a.yaml", Field: "name", Sev: linter.Error},
		{File: "a.yaml", Field: "query", Sev: linter.Error},
		{File: "b.yaml", Field: "fp_rate", Sev: linter.Warning},
	}
	var buf bytes.Buffer
	if err := JSON(&buf, issues, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out jsonOutput
	json.Unmarshal(buf.Bytes(), &out)
	if len(out.Issues) != 3 {
		t.Fatalf("expected 3 issues, got %d", len(out.Issues))
	}
	if out.Errors != 2 || out.Warnings != 1 {
		t.Errorf("expected 2 errors, 1 warning, got %d errors, %d warnings", out.Errors, out.Warnings)
	}
}

func TestOptionsColor(t *testing.T) {
	tests := []struct {
		opts Options
		want bool
	}{
		{Options{Color: true}, true},
		{Options{Color: false}, false},
	}
	for _, tt := range tests {
		got := tt.opts.color("x", ColorRed)
		if tt.want && got != ColorRed+"x"+ColorReset {
			t.Errorf("color=true: got %q, want colored", got)
		}
		if !tt.want && got != "x" {
			t.Errorf("color=false: got %q, want plain", got)
		}
	}
}

func TestJSONWriteError(t *testing.T) {
	issues := []linter.Issue{
		{File: "f.yaml", Field: "x", Sev: linter.Error},
	}
	err := JSON(&failWriter{}, issues, false)
	if err == nil {
		t.Error("expected error from writer, got nil")
	}
}

type failWriter struct{}

func (f *failWriter) Write(p []byte) (int, error) {
	return 0, errWriteFail
}

type errSentinel struct{}

func (e errSentinel) Error() string { return "write failure" }

var errWriteFail errSentinel
