# M2.2: Server Generation (Basic)

**Type**: Epic
**Status**: Done

## Description
**Duration**: 5-7 days
**Goal**: Generate basic HTTP server interfaces and handlers

### Core Work
- [ ] Create `pkg/codegen/http.go` - HTTP constants and helpers
- [ ] Enhance `core.Endpoint` with server-specific metadata
- [ ] Group endpoints by tag for route organization

### Go Generator Work
- [ ] Create server template data structures (HandlerData, RouterData)
- [ ] Implement `convertEndpointsToServer()` function
- [ ] Create server templates:
  - [ ] `server/interface.tmpl`
  - [ ] `server/handlers.tmpl`
  - [ ] `server/routes.tmpl`
- [ ] Update Generate() to include server generation

### Testing
- [ ] Generate server from petstore.yaml
- [ ] Verify server compiles
- [ ] Test basic HTTP server startup
