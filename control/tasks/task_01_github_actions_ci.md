# No CI Automation

**Type**: Task
**Status**: Done
**Priority**: High

## Description
The repository lacks automated quality gates. The `.github` directory does not exist, meaning PRs and commits are not gated by automated tests or linters, leading directly to broken builds.
**Impact**: High
**Likelihood**: High

## Acceptance Criteria
- [x] Add basic `.github/workflows/ci.yml` pipeline.
- [x] Pipeline executes `go test ./...`.
- [x] Pipeline executes `golangci-lint run`.
