# Broken Main Branch Tests

**Type**: Bug
**Status**: Done
**Priority**: High

## Description
The `client_test.go` fails compilation due to mismatched types from the recent `codegen` imports refactor.
**Impact**: High
**Likelihood**: High

## Acceptance Criteria
- [x] Fix compiling errors in `generators/golang/client_test.go`.
- [x] Test suite passes successfully locally.
