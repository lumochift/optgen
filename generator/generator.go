// package generator provides functionality to generate functional option approach by given struct
package generator

import (
	"fmt"
	"go/ast"
	"strings"
)

// Tag represent tag data
type Tag struct {
	Name     string
	DataType string
}

// GetType determine type from expression
func GetType(n ast.Expr) string {
	switch x := n.(type) {
	case *ast.ArrayType:
		return "[]" + GetType(x.Elt)
	case *ast.Ident:
		return x.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", GetType(x.X), GetType(x.Sel))
	case *ast.StarExpr:
		return "*" + GetType(x.X)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", GetType(x.Key), GetType(x.Value))
	default:
		return ""
	}
}

// GetTags from existing tags
func GetTags(f *ast.File, structName, tag string) []Tag {
	if tag == "" {
		tag = `opt`
	}
	tags := []Tag{}
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if x.Name.String() != structName {
				return false
			}
			for _, field := range x.Type.(*ast.StructType).Fields.List {
				if tag == "_all_" || (field.Tag != nil && strings.Contains(field.Tag.Value, tag)) {
					tags = append(tags, Tag{
						Name:     field.Names[0].Name,
						DataType: GetType(field.Type),
					})
				}
			}
		}
		return true
	})
	return tags
}

const CodeTemplate string = `// {{.OptName}} is a {{.StructName}} configurator to be supplied to New{{.StructName}}() function.
type {{.OptName}} func(*{{.StructName}})


// New{{.StructName}} returns a new {{.StructName}}.
func New{{.StructName}}(opt ...{{.OptName}}) (*{{.StructName}}, error) {

	// Prepare a {{.StructName}} with default host.
	{{toLower .StructName}} := &{{.StructName}}{}

	// Apply options.
	for _, o := range opt {
		o({{toLower .StructName}})
	}

	// Do anything here

	return {{toLower .StructName}}, nil
}
{{ $structName:=.StructName}}
{{ range .Tags }}
// Set{{title .Name}} sets the {{title .Name}}
func Set{{title .Name}}({{toLower .Name}} {{.DataType}}) {{.OptName}} {
	return func(c *{{$structName}}) {
		c.{{.Name}} = {{toLower .Name}}
	}
}
{{end}}`
