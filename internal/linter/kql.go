package linter

import (
	"fmt"
	"strings"

	"github.com/daten-krake/tentacle-lint/internal/model"
)

func checkKQL(file string, a *model.Analytic) []Issue {
	query := strings.TrimSpace(a.Query)
	if query == "" {
		return nil
	}

	var issues []Issue
	issues = append(issues, checkUnmatchedDelimiters(file, query)...)
	issues = append(issues, checkPipeSyntax(file, query)...)
	return issues
}

func checkUnmatchedDelimiters(file, query string) []Issue {
	var issues []Issue

	pairs := []struct {
		open  rune
		close rune
		name  string
	}{{'(', ')', "parentheses"}, {'[', ']', "brackets"}, {'{', '}', "braces"}}

	for _, p := range pairs {
		depth := 0
		inString := false
		escaped := false

		for _, ch := range query {
			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == '"' {
				inString = !inString
				continue
			}
			if inString {
				continue
			}
			if ch == p.open {
				depth++
			} else if ch == p.close {
				depth--
				if depth < 0 {
					issues = append(issues, Issue{
						File:    file,
						Field:   "query",
						Message: fmt.Sprintf("unmatched closing %s", p.name),
						Sev:     Error,
					})
					break
				}
			}
		}

		if depth > 0 {
			issues = append(issues, Issue{
				File:    file,
				Field:   "query",
				Message: fmt.Sprintf("unmatched opening %s (%d unclosed)", p.name, depth),
				Sev:     Error,
			})
		}
	}

	return issues
}

func checkPipeSyntax(file, query string) []Issue {
	var issues []Issue

	lines := strings.Split(query, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if trimmed == "|" {
			issues = append(issues, Issue{
				File:    file,
				Field:   "query",
				Message: fmt.Sprintf("line %d: standalone pipe operator with no content", i+1),
				Sev:     Error,
			})
		}
	}

	return issues
}