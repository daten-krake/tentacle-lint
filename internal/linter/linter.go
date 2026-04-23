package linter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/daten-krake/tentacle-lint/internal/model"
	"gopkg.in/yaml.v3"
)

type Severity string

const (
	Error   Severity = "error"
	Warning Severity = "warning"
)

type Issue struct {
	File     string
	Field   string
	Message string
	Sev     Severity
}

func (i Issue) String() string {
	return fmt.Sprintf("%s: %s [%s] %s", i.File, i.Field, i.Sev, i.Message)
}

func (i Issue) EffectiveSev(strict bool) Severity {
	if strict && i.Sev == Warning {
		return Error
	}
	return i.Sev
}

type Config struct {
	Dir       string
	Recursive bool
	Strict    bool
}

func Run(cfg Config) ([]Issue, error) {
	files, err := discover(cfg.Dir, cfg.Recursive)
	if err != nil {
		return nil, fmt.Errorf("discovering yaml files: %w", err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no yaml files found in %s", cfg.Dir)
	}

	var allIssues []Issue
	for _, f := range files {
		issues := lintFile(cfg.Dir, f)
		allIssues = append(allIssues, issues...)
	}

	return allIssues, nil
}

func discover(dir string, recursive bool) ([]string, error) {
	var files []string

	if recursive {
		err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if isYAML(d.Name()) {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("walking directory %s: %w", dir, err)
		}
	} else {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, fmt.Errorf("reading directory %s: %w", dir, err)
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if isYAML(entry.Name()) {
				files = append(files, filepath.Join(dir, entry.Name()))
			}
		}
	}

	return files, nil
}

func isYAML(name string) bool {
	return strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml")
}

func lintFile(baseDir, path string) []Issue {
	relPath := strings.TrimPrefix(path, baseDir+"/")

	data, err := os.ReadFile(path)
	if err != nil {
		return []Issue{{
			File:    relPath,
			Message: fmt.Sprintf("failed to read file: %v", err),
			Sev:     Error,
		}}
	}

	if containsMultiDoc(data) {
		return []Issue{{
			File:    relPath,
			Message: "multi-document yaml files are not supported",
			Sev:     Error,
		}}
	}

	var analytic model.Analytic
	if err := yaml.Unmarshal(data, &analytic); err != nil {
		return []Issue{{
			File:    relPath,
			Message: fmt.Sprintf("failed to parse yaml: %v", err),
			Sev:     Error,
		}}
	}

	var issues []Issue
	issues = append(issues, checkRequired(relPath, &analytic)...)
	issues = append(issues, checkFormat(relPath, &analytic)...)
	issues = append(issues, checkKQL(relPath, &analytic)...)
	issues = append(issues, checkMitre(relPath, &analytic)...)
	issues = append(issues, checkEntityMapping(relPath, &analytic)...)
	issues = append(issues, checkDataSources(relPath, &analytic)...)
	issues = append(issues, checkDuration(relPath, &analytic)...)
	return issues
}

func containsMultiDoc(data []byte) bool {
	seenContent := false
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if trimmed == "---" {
			if seenContent {
				return true
			}
			continue
		}
		seenContent = true
	}
	return false
}