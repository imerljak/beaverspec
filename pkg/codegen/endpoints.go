package codegen

import "github.com/imerljak/beaverspec/pkg/core"

// FilterEndpointsByTag removes endpoints whose first tag is in the excluded set.
// Endpoints with no tags are never excluded.
func FilterEndpointsByTag(endpoints []core.Endpoint, excludeTags []string) []core.Endpoint {
	if len(excludeTags) == 0 {
		return endpoints
	}
	excluded := make(map[string]bool, len(excludeTags))
	for _, t := range excludeTags {
		excluded[t] = true
	}
	var result []core.Endpoint
	for _, ep := range endpoints {
		if len(ep.Tags) > 0 && excluded[ep.Tags[0]] {
			continue
		}
		result = append(result, ep)
	}
	return result
}
