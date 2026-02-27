# No Release Automation

**Type**: Task
**Status**: Done
**Priority**: Low

## Description
The project currently relies on manual SemVer strings (e.g., `v0.1.0-alpha` in `main.go`). Releases should be automated via GitHub.
**Impact**: Low
**Likelihood**: Medium

## Acceptance Criteria
- [x] Use GoReleaser or GitHub Release actions to automate binary distribution.
- [x] Ensure `main.go` version strings are injected correctly.
