# M2.4: Response Headers & Typed Errors

**Type**: Epic
**Status**: Todo

## Description
**Duration**: 4-6 days
**Goal**: Generate typed error responses and handle response headers

### Core Work
- [ ] Update normalizer to extract response headers
- [ ] Map status codes to response schemas
- [ ] Create `pkg/codegen/errors.go` - Error modeling helpers

### Go Generator Work
- [ ] Generate typed error structs
- [ ] Update client error handling (parse errors, return typed errors)
- [ ] Update server error handling (error serialization helpers)
- [ ] Add response header support (client reading, server writing)

### Testing
- [ ] Test typed error returns
- [ ] Test response header handling
