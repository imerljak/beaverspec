# M2.7: Framework Support

**Type**: Epic
**Status**: Todo

## Description
**Duration**: 6-8 days
**Goal**: Support multiple Go HTTP frameworks

### Core Work
- [ ] Create `pkg/codegen/framework.go` - Framework abstraction

### Go Generator Work
- [ ] Framework configuration in config.yaml
- [ ] Implement framework support:
  - [ ] net/http (standard library - done)
  - [ ] Chi router
  - [ ] Echo framework
  - [ ] Gin framework
- [ ] Template organization (frameworks/<name>/)
- [ ] Update template loading for framework overrides

### Testing
- [ ] Generate server for each framework
- [ ] Verify compilation
