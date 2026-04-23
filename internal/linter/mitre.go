package linter

import (
	"fmt"
	"strings"

	"github.com/daten-krake/tentacle-lint/internal/model"
)

var validTactics = map[string]bool{
	"Reconnaissance":      true,
	"ResourceDevelopment": true,
	"InitialAccess":       true,
	"Execution":           true,
	"Persistence":         true,
	"PrivilegeEscalation": true,
	"DefenseEvasion":      true,
	"CredentialAccess":    true,
	"Discovery":           true,
	"LateralMovement":     true,
	"Collection":          true,
	"CommandAndControl":   true,
	"Exfiltration":        true,
	"Impact":              true,
}

func checkMitre(file string, a *model.Analytic) []Issue {
	var issues []Issue

	if len(a.Mitre) == 0 {
		return issues
	}

	hasTactic := false
	for i, m := range a.Mitre {
		for _, t := range m.Tactics {
			t = strings.TrimSpace(t)
			if t == "" {
				issues = append(issues, Issue{
					File:    file,
					Field:   fmt.Sprintf("mitre[%d].tactics", i),
					Message: "empty tactic",
					Sev:     Error,
				})
			} else if !validTactics[t] {
				issues = append(issues, Issue{
					File:    file,
					Field:   fmt.Sprintf("mitre[%d].tactics", i),
					Message: fmt.Sprintf("invalid tactic: %s", t),
					Sev:     Error,
				})
			} else {
				hasTactic = true
			}
		}

		if len(m.Tactics) == 0 {
			issues = append(issues, Issue{
				File:    file,
				Field:   fmt.Sprintf("mitre[%d].tactics", i),
				Message: "no tactics defined",
				Sev:     Warning,
			})
		}

		for j, tech := range m.Techniques {
			if strings.TrimSpace(tech) == "" {
				issues = append(issues, Issue{
					File:    file,
					Field:   fmt.Sprintf("mitre[%d].techniques[%d]", i, j),
					Message: "empty technique",
					Sev:     Warning,
				})
			}
		}

		if len(m.Techniques) == 0 {
			issues = append(issues, Issue{
				File:    file,
				Field:   fmt.Sprintf("mitre[%d].techniques", i),
				Message: "no techniques defined",
				Sev:     Warning,
			})
		}
	}

	if !hasTactic {
		issues = append(issues, Issue{
			File:    file,
			Field:   "mitre",
			Message: "at least one mitre entry must have valid tactics",
			Sev:     Error,
		})
	}

	return issues
}
