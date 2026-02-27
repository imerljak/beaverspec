# Zero Test Coverage on Codegen Package

**Type**: Task
**Status**: Done
**Priority**: Medium

## Description
The new `pkg/codegen/naming.go` has complex pluralization logic with 0% test coverage, making regressions probable.
**Impact**: Medium
**Likelihood**: Medium

## Acceptance Criteria
- [x] Write unit tests for the `pkg/codegen` package.
- [x] Ensure edge cases for string pluralization and casing conversions are covered.
