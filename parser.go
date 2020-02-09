package main

import (
	"go/ast"
	"strings"
)

// GetTags get tag from existing tags
func GetTags(f *ast.File, tag string) []string {
	if tag == "" {
		tag = `opt`
	}
	// Inspect the AST and print all identifiers and literals.
	tags := []string{}
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Field:
			if x.Tag != nil && strings.Contains(x.Tag.Value, tag) {
				tags = append(tags, x.Names[0].Name)
			}
		}
		return true
	})
	return tags
}
