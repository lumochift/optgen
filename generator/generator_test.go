package generator_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lumochift/optgen/generator"
	"github.com/stretchr/testify/assert"
)

func getTokenFileset(path string) (*ast.File, error) {
	fset := token.NewFileSet()
	buffer, err := os.ReadFile(path)
	if err != nil {
		return nil, err

	}
	f, err := parser.ParseFile(fset, "", buffer, 0)
	return f, err
}
func TestGetOptTags(t *testing.T) {
	var (
		field1 = generator.Tag{
			Name:     "Field1",
			DataType: "string",
		}

		field2 = generator.Tag{
			Name:     "Field2",
			DataType: "[]*int",
		}

		field3 = generator.Tag{Name: "Field3",
			DataType: "map[byte]float64"}

		field4 = generator.Tag{Name: "Field4",
			DataType: "http.Header"}
	)
	type args struct {
		sourcePath string
		Tag        string
		structName string
	}
	tests := map[string]struct {
		args args
		want []generator.Tag
	}{
		"foo empty tag": {
			args: args{sourcePath: ".././testfile/foo.sample", structName: "Thing"},
			want: []generator.Tag{field1, field2},
		},
		"foo": {
			args: args{sourcePath: ".././testfile/foo.sample", structName: "Thing"},
			want: []generator.Tag{field1, field2},
		},
		"foo with defined tags": {
			args: args{sourcePath: ".././testfile/foo.sample", Tag: "opt", structName: "Thing"},
			want: []generator.Tag{field1, field2},
		},
		"fooJerax": {
			args: args{sourcePath: ".././testfile/fooJerax.sample", Tag: "jerax", structName: "Thing"},
			want: []generator.Tag{field1, field2},
		},
		"fooJerax want dendy": {
			args: args{sourcePath: ".././testfile/fooJerax.sample", Tag: "dendy", structName: "Thing"},
			want: []generator.Tag{field3},
		},
		"missing struct": {
			args: args{sourcePath: ".././testfile/fooJerax.sample", Tag: "jerax", structName: "Thingy"},
			want: []generator.Tag{},
		},
		"All": {
			args: args{sourcePath: ".././testfile/foo.sample", Tag: "_all_", structName: "Thing"},
			want: []generator.Tag{field1, field2, field3, field4},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f, err := getTokenFileset(tt.args.sourcePath)
			assert.NoError(t, err)
			got := generator.GetTags(f, tt.args.structName, tt.args.Tag)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetOptgenerator.Tags() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_getType(t *testing.T) {
	type args struct {
		n ast.Expr
	}
	tests := map[string]struct {
		args args
		want string
	}{
		"ident": {
			args: args{&ast.Ident{Name: "int"}},
			want: "int",
		},
		"header": {
			args: args{&ast.Ident{Name: "http.Header"}},
			want: "http.Header",
		},
		"empty": {
			args: args{nil},
			want: "",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			if got := generator.GetType(tt.args.n); got != tt.want {
				t.Errorf("GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}
