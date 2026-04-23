package linter

import (
	"fmt"
	"strings"

	"github.com/daten-krake/tentacle-lint/internal/model"
)

func checkEntityMapping(file string, a *model.Analytic) []Issue {
	var issues []Issue

	for i, e := range a.EntityMapping {
		if strings.TrimSpace(e.EntityType) == "" {
			issues = append(issues, Issue{
				File:    file,
				Field:   fmt.Sprintf("entity_mapping[%d].entity_type", i),
				Message: "entity_type is empty",
				Sev:     Error,
			})
		}

		if len(e.FieldMapping) == 0 {
			issues = append(issues, Issue{
				File:    file,
				Field:   fmt.Sprintf("entity_mapping[%d].field_mapping", i),
				Message: "no field mappings defined",
				Sev:     Warning,
			})
		}

		for j, fm := range e.FieldMapping {
			if strings.TrimSpace(fm.Identifier) == "" {
				issues = append(issues, Issue{
					File:    file,
					Field:   fmt.Sprintf("entity_mapping[%d].field_mapping[%d].identifier", i, j),
					Message: "identifier is empty",
					Sev:     Error,
				})
			}
			if strings.TrimSpace(fm.ColumnName) == "" {
				issues = append(issues, Issue{
					File:    file,
					Field:   fmt.Sprintf("entity_mapping[%d].field_mapping[%d].column_name", i, j),
					Message: "column_name is empty",
					Sev:     Error,
				})
			}
		}
	}

	return issues
}

func checkDataSources(file string, a *model.Analytic) []Issue {
	var issues []Issue

	for i, ds := range a.DataSources {
		if strings.TrimSpace(ds.Provider) == "" {
			issues = append(issues, Issue{
				File:    file,
				Field:   fmt.Sprintf("data_sources[%d].provider", i),
				Message: "provider is empty",
				Sev:     Error,
			})
		}
		if strings.TrimSpace(ds.TableName) == "" {
			issues = append(issues, Issue{
				File:    file,
				Field:   fmt.Sprintf("data_sources[%d].table_name", i),
				Message: "table_name is empty",
				Sev:     Error,
			})
		}
	}

	return issues
}
