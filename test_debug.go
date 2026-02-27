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
	minLen := 3
	spec := &core.Spec{
		Models: []core.Model{
			{
				Name: "User",
				Properties: []core.Property{
					{Name: "name", Type: "string", Required: true, MinLength: &minLen},
				},
			},
		},
	}
	result, err := gen.Generate(spec, &core.Config{OutputDir: "generated"})
	if err != nil {
		fmt.Println("err:", err)
		os.Exit(1)
	}
	for _, f := range result.Files {
		if f.Path == "models/models.go" {
			content := string(f.Content)
			lines := strings.Split(content, "\n")
			for i, l := range lines {
				fmt.Printf("%d: %s\n", i+1, l)
			}
		}
	}
}
