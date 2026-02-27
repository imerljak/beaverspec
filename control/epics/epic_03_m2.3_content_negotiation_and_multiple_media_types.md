# M2.3: Content Negotiation & Multiple Media Types

**Type**: Epic
**Status**: Todo

## Description
**Duration**: 4-6 days
**Goal**: Support multiple content types beyond application/json

### Core Work
- [ ] Update `core.MediaType` to include encoding info
- [ ] Update normalizer to extract all content types
- [ ] Create `pkg/codegen/content.go` - Media type helpers

### Go Generator Work
- [ ] Update client to support multiple content types
- [ ] Update server to support multiple content types
- [ ] Add support for:
  - [ ] `application/json` (done)
  - [ ] `application/xml`
  - [ ] `multipart/form-data` (file uploads)
  - [ ] `application/x-www-form-urlencoded`
  - [ ] `text/plain`

### Testing
- [ ] Test client with multiple Accept types
- [ ] Test server with different Content-Type requests
