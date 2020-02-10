// package optgen provides functionality to generate functional option approach by given struct
package optgen

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

func getType(n ast.Expr) string {
	switch x := n.(type) {
	case *ast.ArrayType:
		return "[]" + getType(x.Elt)
	case *ast.Ident:
		return x.Name
	case *ast.StarExpr:
		return "*" + getType(x.X)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", getType(x.Key), getType(x.Value))
	default:
		return ""
	}
}

// GetTags get tag from existing tags
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
				fmt.Println(getType(field.Type))
				if tag == "_all_" || (field.Tag != nil && strings.Contains(field.Tag.Value, tag)) {
					tags = append(tags, Tag{
						Name:     field.Names[0].Name,
						DataType: getType(field.Type),
					})
				}
			}
		}
		return true
	})
	return tags
}

const CodeTemplate string = `// Option is a {{.StructName}} configurator to be supplied to New{{.StructName}}() function.
type Option func(*{{.StructName}})


// New{{.StructName}} returns a new {{.StructName}}.
func New{{.StructName}}(options ...Option) (*{{.StructName}}, error) {

	// Prepare a {{.StructName}} with default host.
	{{toLower .StructName}} := &{{.StructName}}{}

	// Apply options.
	for _, option := range options {
		option({{toLower .StructName}})
	}

	// Do anything here

	return {{toLower .StructName}}, nil
}
{{ $structName:=.StructName}}
{{ range .Tags }}
// Set{{title .Name}} sets the {{title .Name}}
func Set{{title .Name}}({{toLower .Name}} {{.DataType}}) Option {
	return func(c *{{$structName}}) {
		c.{{.Name}} = {{toLower .Name}}
	}
}
{{end}}`
