# tentacle-lint — Core Functionality

## Purpose

tentacle-lint is a linter for the Tentacle DE Framework that validates detection rule YAML files conform to the prodyaml schema and contain well-formed KQL queries. It is designed to run in CI/CD pipelines and be invoked by other repositories.

## Schema Source

The canonical schema is defined in the `tentacle-conv` repository (`internal/model/prodyaml.go`). tentacle-lint mirrors the `Analytic` struct locally in `internal/model/analytic.go` (since Go's `internal` package restriction prevents cross-module imports). A source-tracking comment at the top of the file points to the upstream definition and must be updated when the upstream schema changes.

## Architecture

The linter is a pure-function pipeline:
```
Run(Config) → discover() → lintFile() → check*(file, analytic) → []Issue
```

Each check function is a pure function that takes a file identifier and an `*model.Analytic`, and returns `[]Issue`. There is no mutable shared state. Strict-mode promotion is handled via `Issue.EffectiveSev(strict bool)` which preserves the original severity.

## Validation Rules

### Required Fields

The following fields must be present and non-empty in every YAML file:

| Field         | Type       | Validation                                      |
|---------------|------------|-------------------------------------------------|
| `name`        | string     | Non-empty                                       |
| `severity`    | string     | Must be one of: Informational, Low, Medium, High|
| `query`       | string     | Non-empty, passes basic KQL structural checks   |
| `description` | string     | Non-empty                                       |
| `mitre`      | []Mitre    | At least one entry with valid tactics            |

### Format Checks

- **Severity values**: Must be one of `Informational`, `Low`, `Medium`, `High`
- **query_frequency / query_period**: Must be valid ISO 8601 durations (e.g., `PT5M`, `P1D`)
- **entity_mapping**: If present, each entry must have `entity_type` and at least one `field_mapping` with both `identifier` and `column_name`
- **data_sources**: If present, each entry must have `provider` and `table_name`
- **references**: If present, empty entries are warned
- **tags / os_family**: Warned when empty

### KQL Structural Validation

Basic structural checks on the `query` field:

1. **Non-empty**: Query must not be empty or whitespace-only (covered by required check)
2. **Unmatched delimiters**: Detect unmatched `()`, `[]`, `{}`
3. **Pipe syntax**: Detect standalone pipe operators with no content

### MITRE ATT&CK Validation

- Tactic names must be valid MITRE ATT&CK tactics (14 enterprise tactics)
- Empty tactic/technique entries are warned
- At least one mitre entry must have valid tactics

### Optional Fields (warned when empty, not errored)

`fp_rate`, `permission_required`, `tags`, `os_family`, `technical_description`, `considerations`, `false_positives`, `blindspots`, `response_plan`, `test_block`

## Invocation

### CLI Binary

```bash
# Lint all YAML files recursively in current directory
tentacle-lint

# Lint specific directory (recursive by default)
tentacle-lint --dir ./detection-rules

# Non-recursive (flat directory only)
tentacle-lint --dir ./rules --recursive=false

# Strict mode (warnings = errors)
tentacle-lint --strict

# JSON output
tentacle-lint --format json
```

### Exit Codes

| Code | Meaning          |
|------|------------------|
| 0    | All files pass   |
| 1    | Lint failures found |
| 2    | Runtime error (bad flags, no files, I/O errors) |

### GitHub Action

```yaml
- uses: daten-krake/tentacle-lint@v1
  with:
    directory: './yaml'
    recursive: true
    strict: false
    format: text
```

## JSON Output Format

```json
{
  "issues": [
    {
      "file": "rule.yaml",
      "field": "severity",
      "message": "must be one of: Informational, Low, Medium, High",
      "severity": "error",
      "effective_severity": "error",
      "promoted": false
    }
  ],
  "errors": 1,
  "warnings": 0
}
```

## Release Pipeline

- GoReleaser builds cross-platform binaries (linux/darwin, amd64/arm64)
- GitHub Action wraps the CLI for workflow use
- CI pipeline runs on PR/push (test + vet)
- Release pipeline triggered on `v*` tag push