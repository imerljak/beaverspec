# M2.6: Request Validation

**Type**: Epic
**Status**: Todo

## Description
**Duration**: 5-7 days
**Goal**: Generate validation code for requests and responses

### Core Work
- [ ] Add format validation (email, uri, date-time, uuid, etc.)
- [ ] Create `pkg/codegen/validation.go` - Validation helpers

### Go Generator Work
- [ ] Generate Validate() method for each model
- [ ] Check required fields, constraints, formats, enums
- [ ] Client-side validation (optional pre-request)
- [ ] Server-side validation (auto-validate, return 400)

### Testing
- [ ] Test validation with valid/invalid data
- [ ] Test all constraint types
