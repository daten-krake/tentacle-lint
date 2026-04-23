package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/daten-krake/tentacle-lint/internal/linter"
)

func Text(w io.Writer, issues []linter.Issue, strict bool) {
	if len(issues) == 0 {
		fmt.Fprintln(w, "No issues found.")
		return
	}
	for _, issue := range issues {
		sev := issue.EffectiveSev(strict)
		fmt.Fprintf(w, "%s: %s [%s] %s\n", issue.File, issue.Field, sev, issue.Message)
	}
}

type jsonIssue struct {
	File            string `json:"file"`
	Field           string `json:"field"`
	Message         string `json:"message"`
	Severity        string `json:"severity"`
	EffectiveSev    string `json:"effective_severity"`
	Promoted        bool   `json:"promoted"`
}

type jsonOutput struct {
	Issues []jsonIssue `json:"issues"`
	Errors int         `json:"errors"`
	Warns  int         `json:"warnings"`
}

func JSON(w io.Writer, issues []linter.Issue, strict bool) {
	out := jsonOutput{
		Issues: make([]jsonIssue, 0, len(issues)),
	}

	for _, issue := range issues {
		promoted := strict && issue.Sev == linter.Warning
		effectiveSev := linter.Warning
		if promoted || issue.Sev == linter.Error {
			effectiveSev = linter.Error
		}
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
			out.Warns++
		}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(out)
}