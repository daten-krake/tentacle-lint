package linter

import (
	"regexp"

	"github.com/daten-krake/tentacle-lint/internal/model"
)

var iso8601Pattern = regexp.MustCompile(`^P(?:(\d+)Y)?(?:(\d+)M)?(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+(?:\.\d+)?)S)?)?$`)

func checkDuration(file string, a *model.Analytic) []Issue {
	var issues []Issue

	if a.QueryFrequency != "" && !isValidISO8601Duration(a.QueryFrequency) {
		issues = append(issues, Issue{
			File:    file,
			Field:   "query_frequency",
			Message: "must be a valid ISO 8601 duration (e.g., PT5M, PT1H, P1D)",
			Sev:     Error,
		})
	}
	if a.QueryPeriod != "" && !isValidISO8601Duration(a.QueryPeriod) {
		issues = append(issues, Issue{
			File:    file,
			Field:   "query_period",
			Message: "must be a valid ISO 8601 duration (e.g., PT5M, PT1H, P1D)",
			Sev:     Error,
		})
	}

	return issues
}

func isValidISO8601Duration(s string) bool {
	if !iso8601Pattern.MatchString(s) {
		return false
	}
	matches := iso8601Pattern.FindStringSubmatch(s)
	for i := 1; i < len(matches); i++ {
		if matches[i] != "" {
			return true
		}
	}
	return false
}