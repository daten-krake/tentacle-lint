package linter

import (
	"strings"

	"github.com/daten-krake/tentacle-lint/internal/model"
)

func checkRequired(file string, a *model.Analytic) []Issue {
	var issues []Issue

	required := []struct {
		field string
		value string
	}{
		{"name", a.Name},
		{"severity", a.Severity},
		{"query", a.Query},
		{"description", a.Description},
	}

	for _, r := range required {
		if strings.TrimSpace(r.value) == "" {
			issues = append(issues, Issue{
				File:    file,
				Field:   r.field,
				Message: "required field is empty",
				Sev:     Error,
			})
		}
	}

	if len(a.Mitre) == 0 {
		issues = append(issues, Issue{
			File:    file,
			Field:   "mitre",
			Message: "at least one mitre entry is required",
			Sev:     Error,
		})
	}

	optional := []struct {
		field string
		value string
	}{
		{"fp_rate", a.FPRate},
		{"permission_required", a.PermissionRequired},
		{"technical_description", a.TechnicalDescription},
		{"considerations", a.Considerations},
		{"false_positives", a.FalsePositives},
		{"blindspots", a.Blindspots},
		{"response_plan", a.ResponsePlan},
		{"test_block", a.TestBlock},
	}

	for _, o := range optional {
		if strings.TrimSpace(o.value) == "" {
			issues = append(issues, Issue{
				File:    file,
				Field:   o.field,
				Message: "optional field is empty",
				Sev:     Warning,
			})
		}
	}

	return issues
}
