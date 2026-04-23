# Changelog

All notable changes to tentacle-lint are tracked here.

## [Unreleased] — 2026-04-23

### Added
- Core linter with pure-function pipeline architecture
- YAML file discovery (recursive with `--recursive` flag, defaults to true)
- YAML parsing against prodyaml `Analytic` schema (mirrored from tentacle-conv)
- Required field validation (name, severity, query, description, mitre)
- Format checks (valid severity values, missing tags/os_family warnings)
- KQL structural validation (unmatched delimiters, pipe syntax)
- MITRE ATT&CK validation (14 enterprise tactics, technique presence)
- Entity mapping structure validation (entity_type, field_mapping)
- Data source structure validation (provider, table_name)
- ISO 8601 duration validation for query_frequency/query_period
- Multi-document YAML detection
- Strict mode (`--strict` promotes warnings to errors, preserves original severity)
- CLI with `--dir`, `--recursive`, `--strict`, `--format` (text/json), `--version` flags
- Release pipeline using plain GitHub Actions (no GoReleaser dependency)
- Cross-platform builds via Go cross-compilation (linux/darwin, amd64/arm64)
- SHA256 checksums per binary with aggregated checksums.txt in releases
- Exit codes: 0 (pass), 1 (lint failure), 2 (runtime error)
- JSON output with `severity`, `effective_severity`, and `promoted` fields
- GitHub Action (action.yml) with authenticated latest-version resolution, checksum verification before extraction, and proper exit code propagation
- CI workflow (test, vet, format check, build on push/PR)
- Release workflow (3-job pipeline: test → build matrix → release with changelog)
- 23 unit tests covering all checkers