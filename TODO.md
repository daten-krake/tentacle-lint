# tentacle-lint TODO

## Project Setup
- [x] Initialize Go module (`go.mod`) 
- [x] Create project directory structure (cmd/, internal/)
- [x] Set up Makefile (build, test, vet, fmt, lint, clean)

## Core Linter Implementation
- [x] Implement YAML file discovery (recursive + flat, find all .yaml/.yml)
- [x] Implement YAML parsing using local mirror of `model.Analytic` struct
- [x] Implement required field validation (name, severity, query, description, mitre)
- [x] Implement format/structure checks (valid severity values, missing tags/os_family)
- [x] Implement basic KQL structural validation:
  - [x] Non-empty query check
  - [x] Unmatched parentheses/brackets/braces detection
  - [x] Pipe syntax validation (standalone pipe detection)
- [x] Implement MITRE ATT&CK format validation (valid tactics/techniques)
- [x] Implement entity_mapping structure validation
- [x] Implement data_sources structure validation
- [x] Implement ISO 8601 duration validation for query_frequency/query_period
- [x] Output results in standard linter format (stdout with exit codes)
- [x] Multi-document YAML detection
- [x] Strict mode (warnings promoted to errors)

## CLI Interface
- [x] Implement CLI entry point (`cmd/tentacle-lint/main.go`)
- [x] Support flags: `--dir`, `--recursive`, `--strict`, `--format` (text/json)
- [x] Format flag validation
- [x] Exit code 0 on pass, 1 on lint failure, 2 on runtime error

## Testing
- [x] Unit tests for all validation rules (23 tests passing)
- [ ] Integration tests with sample YAML files (from tentacle-conv testdata)
- [ ] Test CLI flag parsing and exit codes
- [ ] Test JSON output format

## GitHub Action
- [ ] Create `action.yml` for GitHub Action wrapper
- [ ] Support `directory`, `recursive`, `strict`, and `format` inputs
- [ ] Support `result` output

## Release Pipeline
- [x] Create `.github/workflows/ci.yml` (lint + test on PR/push)
- [x] Create `.github/workflows/release.yml` (plain GitHub Actions on tag push)
- [x] Cross-platform builds via Go cross-compilation (linux/darwin, amd64/arm64)
- [x] SHA256 checksums generated per binary and aggregated
- [x] Auto-generated changelog from git log
- [x] GitHub release with artifacts and checksums

## Documentation
- [ ] Update README.md with usage instructions
- [ ] Document all CLI flags and exit codes
- [ ] Document GitHub Action usage example