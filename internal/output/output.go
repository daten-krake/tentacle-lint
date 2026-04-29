package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/daten-krake/tentacle-lint/internal/linter"
)

const (
	ColorRed     = "\033[31m"
	ColorYellow  = "\033[33m"
	ColorCyan    = "\033[36m"
	ColorMagenta = "\033[35m"
	ColorBold    = "\033[1m"
	ColorReset   = "\033[0m"
)

type Options struct {
	Color bool
}

func (o Options) color(s string, c string) string {
	if !o.Color {
		return s
	}
	return c + s + ColorReset
}

func Text(w io.Writer, issues []linter.Issue, strict bool, opts Options) {
	if len(issues) == 0 {
		fmt.Fprintln(w, "No issues found.")
		return
	}
	for _, issue := range issues {
		sev := issue.EffectiveSev(strict)
		sevColor := ColorYellow
		if sev == linter.Error {
			sevColor = ColorRed
		}
		field := issue.Field
		if field == "" {
			field = "-"
		}
		fmt.Fprintf(w, "%s: %s [%s] %s\n",
			opts.color(issue.File, ColorCyan),
			opts.color(field, ColorMagenta),
			opts.color(string(sev), sevColor+ColorBold),
			issue.Message,
		)
	}
}

type jsonIssue struct {
	File         string `json:"file"`
	Field        string `json:"field"`
	Message      string `json:"message"`
	Severity     string `json:"severity"`
	EffectiveSev string `json:"effective_severity"`
	Promoted     bool   `json:"promoted"`
}

type jsonOutput struct {
	Issues   []jsonIssue `json:"issues"`
	Errors   int         `json:"errors"`
	Warnings int         `json:"warnings"`
}

func JSON(w io.Writer, issues []linter.Issue, strict bool) error {
	out := jsonOutput{
		Issues: make([]jsonIssue, 0, len(issues)),
	}

	for _, issue := range issues {
		effectiveSev := issue.EffectiveSev(strict)
		promoted := issue.Sev != effectiveSev
		out.Issues = append(out.Issues, jsonIssue{
			File:         issue.File,
			Field:        issue.Field,
			Message:      issue.Message,
			Severity:     string(issue.Sev),
			EffectiveSev: string(effectiveSev),
			Promoted:     promoted,
		})
		if effectiveSev == linter.Error {
			out.Errors++
		} else {
			out.Warnings++
		}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
