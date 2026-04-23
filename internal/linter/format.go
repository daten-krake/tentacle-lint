package linter

import (
	"fmt"
	"strings"

	"github.com/daten-krake/tentacle-lint/internal/model"
)

var validSeverities = map[string]bool{
	"Informational": true,
	"Low":           true,
	"Medium":        true,
	"High":          true,
}

func checkFormat(file string, a *model.Analytic) []Issue {
	var issues []Issue

	if a.Severity != "" && !validSeverities[strings.TrimSpace(a.Severity)] {
		issues = append(issues, Issue{
			File:    file,
			Field:   "severity",
			Message: "must be one of: Informational, Low, Medium, High",
			Sev:     Error,
		})
	}

	if len(a.Tags) == 0 {
		issues = append(issues, Issue{
			File:    file,
			Field:   "tags",
			Message: "no tags defined",
			Sev:     Warning,
		})
	}

	if len(a.OSFamily) == 0 {
		issues = append(issues, Issue{
			File:    file,
			Field:   "os_family",
			Message: "no os_family defined",
			Sev:     Warning,
		})
	}

	for i, ref := range a.References {
		if strings.TrimSpace(ref) == "" {
			issues = append(issues, Issue{
				File:    file,
				Field:   fmt.Sprintf("references[%d]", i),
				Message: "reference is empty",
				Sev:     Warning,
			})
		}
	}

	return issues
}
