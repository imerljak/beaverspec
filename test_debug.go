package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/imerljak/beaverspec/generators/golang"
	"github.com/imerljak/beaverspec/pkg/core"
)

func main() {
	gen := golang.NewGenerator()

	ep := core.Endpoint{
		OperationID: "auth",
		Path:        "/auth",
		Method:      "POST",
		Responses: []core.Response{
			{
				StatusCode: "200",
				Content: map[string]core.MediaType{
					"application/json": {Schema: &core.Property{RefType: "Session"}},
				},
			},
			{
				StatusCode: "401",
				Content: map[string]core.MediaType{
					"application/json": {Schema: &core.Property{RefType: "AuthError"}},
				},
			},
		},
	}
	spec := &core.Spec{
		Endpoints: []core.Endpoint{ep},
	}

	result, err := gen.Generate(spec, &core.Config{
		OutputDir: "generated",
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	for _, f := range result.Files {
		if f.Path == "client/client.go" {
			lines := strings.Split(string(f.Content), "\n")
			for i, line := range lines {
				if strings.Contains(line, "switch resp.StatusCode") {
					for j := 0; j < 20 && i+j < len(lines); j++ {
						fmt.Println(lines[i+j])
					}
					break
				}
			}
		}
	}
}
