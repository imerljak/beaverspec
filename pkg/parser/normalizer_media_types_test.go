package parser

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/imerljak/beaverspec/pkg/codegen"
)

func TestExtractMediaTypes(t *testing.T) {
	n := NewNormalizer()

	desc := "OK"
	doc := &openapi3.T{
		Paths: openapi3.NewPaths(
			openapi3.WithPath("/data", &openapi3.PathItem{
				Post: &openapi3.Operation{
					OperationID: "postData",
					RequestBody: &openapi3.RequestBodyRef{
						Value: &openapi3.RequestBody{
							Content: openapi3.Content{
								"application/xml": &openapi3.MediaType{},
								"multipart/form-data": &openapi3.MediaType{
									Encoding: map[string]*openapi3.Encoding{
										"file": {ContentType: "image/png"},
									},
								},
							},
						},
					},
					Responses: openapi3.NewResponses(
						openapi3.WithStatus(200, &openapi3.ResponseRef{
							Value: &openapi3.Response{
								Description: &desc,
								Content: openapi3.Content{
									"application/xml": &openapi3.MediaType{},
								},
							},
						}),
					),
				},
			}),
		),
	}

	endpoints := n.extractEndpoints(doc.Paths)
	if len(endpoints) != 1 {
		t.Fatalf("expected 1 endpoint, got %d", len(endpoints))
	}

	ep := endpoints[0]

	// Verify request bodies
	if contentType, ok := ep.RequestBody.Content[codegen.MediaTypeXML]; !ok {
		t.Errorf("expected application/xml in RequestBody, not found")
	} else {
		if len(contentType.Encoding) > 0 {
			t.Errorf("XML should not have encodings, found %d", len(contentType.Encoding))
		}
	}

	if contentType, ok := ep.RequestBody.Content[codegen.MediaTypeMultipartForm]; !ok {
		t.Errorf("expected multipart/form-data in RequestBody, not found")
	} else {
		if len(contentType.Encoding) != 1 {
			t.Errorf("expected 1 encoding for multipart, found %d", len(contentType.Encoding))
		}
		if contentType.Encoding["file"].ContentType != "image/png" {
			t.Errorf("expected encoding content type image/png, got %s", contentType.Encoding["file"].ContentType)
		}
	}

	// Verify response body
	if resp, ok := ep.Responses[0].Content[codegen.MediaTypeXML]; !ok {
		t.Errorf("expected application/xml in Response, not found")
	} else {
		if len(resp.Encoding) > 0 {
			t.Errorf("XML should not have encodings in response, found %d", len(resp.Encoding))
		}
	}
}
