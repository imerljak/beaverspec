# M2.5: Security Scheme Implementation

**Type**: Epic
**Status**: Todo

## Description
**Duration**: 5-7 days
**Goal**: Generate authentication/authorization code

### Core Work
- [ ] Update normalizer to apply security requirements to endpoints
- [ ] Handle global vs operation-level security
- [ ] Create `pkg/codegen/security.go` - Security utilities

### Go Generator Work
- [ ] Client authentication:
  - [ ] API Key (header, query, cookie)
  - [ ] HTTP Basic Auth
  - [ ] HTTP Bearer Token
  - [ ] OAuth2 (basic token injection)
- [ ] Server authentication middleware:
  - [ ] Generate auth middleware interfaces
  - [ ] Generate auth verification stubs
  - [ ] Add auth context to handlers

### Testing
- [ ] Test client with different auth schemes
- [ ] Test server auth middleware integration
