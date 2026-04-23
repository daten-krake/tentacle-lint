package model

// Analytic and related types are mirrored from tentacle-conv/internal/model/prodyaml.go.
// When updating, sync from the upstream repository to keep the schema in sync.
// Source: https://github.com/daten-krake/tentacle-conv/blob/main/internal/model/prodyaml.go

type Analytic struct {
	ID                   string       `yaml:"id"`
	Name                 string       `yaml:"name"`
	Severity             string       `yaml:"severity"`
	FPRate               string       `yaml:"fp_rate"`
	PermissionRequired   string       `yaml:"permission_required"`
	Mitre                []Mitre      `yaml:"mitre"`
	EntityMapping        []Entities   `yaml:"entity_mapping"`
	DataSources          []DataSource `yaml:"data_sources"`
	Tags                 []string     `yaml:"tags"`
	OSFamily             []string     `yaml:"os_family"`
	Description          string       `yaml:"description"`
	TechnicalDescription string       `yaml:"technical_description"`
	Considerations       string       `yaml:"considerations"`
	FalsePositives       string       `yaml:"false_positives"`
	Blindspots           string       `yaml:"blindspots"`
	ResponsePlan         string       `yaml:"response_plan"`
	References           []string     `yaml:"references"`
	Query                string       `yaml:"query"`
	TestBlock            string       `yaml:"test_block"`
	QueryFrequency       string       `yaml:"query_frequency"`
	QueryPeriod          string       `yaml:"query_period"`
}

type Mitre struct {
	Tactics    []string `yaml:"tactics"`
	Techniques []string `yaml:"techniques"`
}

type Entities struct {
	EntityType   string         `yaml:"entity_type"`
	FieldMapping []FieldMapping `yaml:"field_mapping"`
}

type FieldMapping struct {
	Identifier string `yaml:"identifier"`
	ColumnName string `yaml:"column_name"`
}

type DataSource struct {
	Provider  string `yaml:"provider" json:"provider"`
	EventID   string `yaml:"event_id" json:"event_id"`
	TableName string `yaml:"table_name" json:"table_name"`
}
