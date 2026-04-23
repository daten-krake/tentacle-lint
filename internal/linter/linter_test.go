package linter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeYAML(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatalf("writing yaml file: %v", err)
	}
}

func validAnalyticYAML() string {
	return `
name: Test Rule
severity: High
description: A test detection rule
query: |
  SecurityAlert
  | where Severity == "High"
  | extend AlertName = AlertName
mitre:
  - tactics:
      - InitialAccess
    techniques:
      - T1190
tags:
  - test
os_family:
  - Windows
fp_rate: Low
permission_required: User
technical_description: Technical details
considerations: Some considerations
false_positives: Legitimate admin activity
blindspots: None known
response_plan: Investigate and escalate
references:
  - https://example.com
query_frequency: PT5M
query_period: PT1H
entity_mapping:
  - entity_type: Account
    field_mapping:
      - identifier: Name
        column_name: AccountName
data_sources:
  - provider: Azure
    event_id: "4688"
    table_name: SecurityAlert
test_block: |
  let test = SecurityAlert | limit 10;
`
}

func TestLintValidFile(t *testing.T) {
	dir := t.TempDir()
	writeYAML(t, dir, "valid.yaml", validAnalyticYAML())

	issues, err := Run(Config{Dir: dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, issue := range issues {
		if issue.Sev == Error {
			t.Errorf("unexpected error: %s", issue)
		}
	}
}

func TestLintMissingRequiredFields(t *testing.T) {
	dir := t.TempDir()
	writeYAML(t, dir, "missing.yaml", `
name: ""
severity: ""
description: ""
query: ""
mitre: []
`)
	issues, err := Run(Config{Dir: dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	errorFields := map[string]bool{}
	for _, issue := range issues {
		if issue.Sev == Error {
			errorFields[issue.Field] = true
		}
	}

	for _, field := range []string{"name", "severity", "query", "description", "mitre"} {
		if !errorFields[field] {
			t.Errorf("expected error for field %s, not found", field)
		}
	}
}

func TestLintInvalidSeverity(t *testing.T) {
	dir := t.TempDir()
	writeYAML(t, dir, "bad-sev.yaml", `
name: Test
severity: Critical
description: test
query: SecurityAlert | where true
mitre:
  - tactics:
      - InitialAccess
`)
	issues, err := Run(Config{Dir: dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, issue := range issues {
		if issue.Field == "severity" && issue.Sev == Error {
			found = true
		}
	}
	if !found {
		t.Error("expected error for invalid severity")
	}
}

func TestLintKQLUnmatchedParens(t *testing.T) {
	dir := t.TempDir()
	writeYAML(t, dir, "unmatched.yaml", `
name: Test
severity: High
description: test
query: |
  SecurityAlert
  | where (Severity == "High"
mitre:
  - tactics:
      - InitialAccess
`)
	issues, err := Run(Config{Dir: dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, issue := range issues {
		if issue.Field == "query" && issue.Sev == Error {
			found = true
		}
	}
	if !found {
		t.Error("expected error for unmatched parentheses")
	}
}

func TestLintKQLUnmatchedBrackets(t *testing.T) {
	dir := t.TempDir()
	writeYAML(t, dir, "bracket.yaml", `
name: Test
severity: High
description: test
query: |
  let arr = dynamic([1, 2, 3
  SecurityAlert | where true
mitre:
  - tactics:
      - InitialAccess
`)
	issues, err := Run(Config{Dir: dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, issue := range issues {
		if issue.Field == "query" && issue.Sev == Error && contains(issue.Message, "brackets") {
			found = true
		}
	}
	if !found {
		t.Error("expected error for unmatched brackets")
	}
}

func TestLintKQLUnmatchedBraces(t *testing.T) {
	dir := t.TempDir()
	writeYAML(t, dir, "braces.yaml", `
name: Test
severity: High
description: test
query: |
  let obj = dynamic({"key": "value"
  SecurityAlert | where true
mitre:
  - tactics:
      - InitialAccess
`)
	issues, err := Run(Config{Dir: dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, issue := range issues {
		if issue.Field == "query" && issue.Sev == Error && contains(issue.Message, "braces") {
			found = true
		}
	}
	if !found {
		t.Error("expected error for unmatched braces")
	}
}

func TestLintPipeSyntaxStandalonePipe(t *testing.T) {
	dir := t.TempDir()
	writeYAML(t, dir, "pipe.yaml", `
name: Test
severity: High
description: test
query: |
  SecurityAlert
  |
  | where true
mitre:
  - tactics:
      - InitialAccess
`)
	issues, err := Run(Config{Dir: dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, issue := range issues {
		if issue.Field == "query" && issue.Sev == Error {
			found = true
		}
	}
	if !found {
		t.Error("expected error for standalone pipe")
	}
}

func TestLintInvalidTactics(t *testing.T) {
	dir := t.TempDir()
	writeYAML(t, dir, "bad-tactic.yaml", `
name: Test
severity: High
description: test
query: SecurityAlert | where true
mitre:
  - tactics:
      - InvalidTactic
`)
	issues, err := Run(Config{Dir: dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, issue := range issues {
		if issue.Field == "mitre[0].tactics" && issue.Sev == Error {
			found = true
		}
	}
	if !found {
		t.Error("expected error for invalid tactic")
	}
}

func TestLintMitreNoValidTactics(t *testing.T) {
	dir := t.TempDir()
	writeYAML(t, dir, "no-tactics.yaml", `
name: Test
severity: High
description: test
query: SecurityAlert | where true
mitre:
  - tactics: []
    techniques:
      - T1190
`)
	issues, err := Run(Config{Dir: dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, issue := range issues {
		if issue.Field == "mitre" && issue.Sev == Error {
			found = true
		}
	}
	if !found {
		t.Error("expected error when no valid tactics exist across all mitre entries")
	}
}

func TestLintInvalidDuration(t *testing.T) {
	dir := t.TempDir()
	writeYAML(t, dir, "bad-duration.yaml", `
name: Test
severity: High
description: test
query: SecurityAlert | where true
mitre:
  - tactics:
      - InitialAccess
query_frequency: "5 minutes"
query_period: "1 hour"
`)
	issues, err := Run(Config{Dir: dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := 0
	for _, issue := range issues {
		if (issue.Field == "query_frequency" || issue.Field == "query_period") && issue.Sev == Error {
			found++
		}
	}
	if found != 2 {
		t.Errorf("expected 2 duration errors, got %d", found)
	}
}

func TestDurationDegenerateForms(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"PT5M", true},
		{"PT1H", true},
		{"P1D", true},
		{"P", false},
		{"PT", false},
		{"", false},
		{"5 minutes", false},
		{"1h", false},
	}
	for _, tt := range tests {
		result := isValidISO8601Duration(tt.input)
		if result != tt.expected {
			t.Errorf("isValidISO8601Duration(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestStrictModePromotion(t *testing.T) {
	dir := t.TempDir()
	writeYAML(t, dir, "strict.yaml", `
name: Test
severity: High
description: test
query: SecurityAlert | where true
mitre:
  - tactics:
      - InitialAccess
fp_rate: ""
`)

	issues, err := Run(Config{Dir: dir, Strict: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, issue := range issues {
		if issue.Field == "fp_rate" && issue.EffectiveSev(true) == Error {
			found = true
		}
	}
	if !found {
		t.Error("expected fp_rate warning to be promoted to error in strict mode")
	}
}

func TestStrictModePreservesOriginal(t *testing.T) {
	dir := t.TempDir()
	writeYAML(t, dir, "preserve.yaml", `
name: Test
severity: High
description: test
query: SecurityAlert | where true
mitre:
  - tactics:
      - InitialAccess
fp_rate: ""
`)

	issues, err := Run(Config{Dir: dir, Strict: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, issue := range issues {
		if issue.Field == "fp_rate" {
			if issue.Sev != Warning {
				t.Error("original severity should be preserved as warning")
			}
			if issue.EffectiveSev(true) != Error {
				t.Error("effective severity should be error in strict mode")
			}
		}
	}
}

func TestEmptyDirectory(t *testing.T) {
	dir := t.TempDir()
	_, err := Run(Config{Dir: dir})
	if err == nil {
		t.Error("expected error for empty directory")
	}
}

func TestLintEntityMappingValidation(t *testing.T) {
	dir := t.TempDir()
	writeYAML(t, dir, "entity.yaml", `
name: Test
severity: High
description: test
query: SecurityAlert | where true
mitre:
  - tactics:
      - InitialAccess
entity_mapping:
  - entity_type: ""
    field_mapping: []
`)
	issues, err := Run(Config{Dir: dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	errorFields := map[string]bool{}
	for _, issue := range issues {
		if issue.Sev == Error {
			errorFields[issue.Field] = true
		}
	}
	if !errorFields["entity_mapping[0].entity_type"] {
		t.Error("expected error for empty entity_type")
	}
}

func TestLintDataSourcesValidation(t *testing.T) {
	dir := t.TempDir()
	writeYAML(t, dir, "datasource.yaml", `
name: Test
severity: High
description: test
query: SecurityAlert | where true
mitre:
  - tactics:
      - InitialAccess
data_sources:
  - provider: ""
    event_id: "1"
    table_name: ""
`)
	issues, err := Run(Config{Dir: dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	errorFields := map[string]bool{}
	for _, issue := range issues {
		if issue.Sev == Error {
			errorFields[issue.Field] = true
		}
	}
	if !errorFields["data_sources[0].provider"] {
		t.Error("expected error for empty provider")
	}
	if !errorFields["data_sources[0].table_name"] {
		t.Error("expected error for empty table_name")
	}
}

func TestRecursiveDirectoryDiscovery(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "subdir")
	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatalf("creating subdir: %v", err)
	}
	writeYAML(t, dir, "root.yaml", validAnalyticYAML())
	writeYAML(t, sub, "nested.yaml", validAnalyticYAML())

	files, err := discover(dir, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d: %v", len(files), files)
	}

	hasRoot := false
	hasNested := false
	for _, f := range files {
		if filepath.Base(f) == "root.yaml" {
			hasRoot = true
		}
		if filepath.Base(f) == "nested.yaml" {
			hasNested = true
		}
	}
	if !hasRoot {
		t.Error("expected root.yaml to be discovered")
	}
	if !hasNested {
		t.Error("expected nested.yaml to be discovered recursively")
	}
}

func TestNonRecursiveDirectoryDiscovery(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "subdir")
	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatalf("creating subdir: %v", err)
	}
	writeYAML(t, dir, "root.yaml", validAnalyticYAML())
	writeYAML(t, sub, "nested.yaml", `
name: Bad
severity: Invalid
description: test
query: SecurityAlert | where true
mitre: []
`)

	issues, err := Run(Config{Dir: dir, Recursive: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, issue := range issues {
		if contains(issue.File, "subdir") {
			t.Error("non-recursive mode should not check subdirectory files")
		}
	}
}

func TestMalformedYAML(t *testing.T) {
	dir := t.TempDir()
	writeYAML(t, dir, "broken.yaml", `{{invalid yaml: [}`)
	issues, err := Run(Config{Dir: dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, issue := range issues {
		if issue.Sev == Error && contains(issue.Message, "parse yaml") {
			found = true
		}
	}
	if !found {
		t.Error("expected parse error for malformed yaml")
	}
}

func TestEffectiveSev(t *testing.T) {
	warning := Issue{Sev: Warning}
	if warning.EffectiveSev(false) != Warning {
		t.Error("non-strict should keep warning as warning")
	}
	if warning.EffectiveSev(true) != Error {
		t.Error("strict should promote warning to error")
	}

	err := Issue{Sev: Error}
	if err.EffectiveSev(false) != Error {
		t.Error("error should stay error in non-strict")
	}
	if err.EffectiveSev(true) != Error {
		t.Error("error should stay error in strict")
	}
}

func TestMultiDocDetection(t *testing.T) {
	dir := t.TempDir()

	writeYAML(t, dir, "leading-sep.yaml", "---\nname: Test\nseverity: High\ndescription: test\nquery: SecurityAlert | where true\nmitre:\n  - tactics:\n      - InitialAccess\n")
	issues, err := Run(Config{Dir: dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, issue := range issues {
		if contains(issue.Message, "multi-document") {
			t.Error("leading --- should not trigger multi-doc detection")
		}
	}

	writeYAML(t, dir, "multi-doc.yaml", "name: First\nseverity: High\n---\nname: Second\nseverity: Low\n")
	for _, issue := range issues {
		if contains(issue.Message, "multi-document") {
			t.Error("previous file issues leaking")
		}
	}
}

func TestLeadingSeparatorAllowed(t *testing.T) {
	data := []byte("---\nname: Test\nseverity: High\n")
	if containsMultiDoc(data) {
		t.Error("leading --- should not be flagged as multi-doc")
	}
}

func TestActualMultiDoc(t *testing.T) {
	data := []byte("name: First\nseverity: High\n---\nname: Second\nseverity: Low\n")
	if !containsMultiDoc(data) {
		t.Error("actual multi-doc should be flagged")
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}