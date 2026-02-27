# Lack of Structured Logging

**Type**: Task
**Status**: Done
**Priority**: Low

## Description
The CLI relies on standard `fmt.Printf` / `log.Fatalf` resulting in lack of observability making debugging broken plugins/templates difficult.
**Impact**: Low
**Likelihood**: High

## Acceptance Criteria
- [x] Integrate a structured logger (e.g., `log/slog`).
- [x] Ensure formatting handles complex errors during the generation process cleanly.
