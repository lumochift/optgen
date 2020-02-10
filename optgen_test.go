package optgen_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lumochift/optgen"
	"github.com/stretchr/testify/assert"
)

func getTokenFileset(path string) (*ast.File, error) {
	fset := token.NewFileSet()
	buffer, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err

	}
	f, err := parser.ParseFile(fset, "", buffer, 0)
	return f, err
}
func TestGetOptTags(t *testing.T) {
	var (
		field1 = optgen.Tag{
			Name:     "Field1",
			DataType: "string",
		}

		field2 = optgen.Tag{
			Name:     "Field2",
			DataType: "[]*int",
		}

		field3 = optgen.Tag{Name: "Field3",
			DataType: "map[byte]float64"}
	)
	type args struct {
		sourcePath string
		Tag        string
		structName string
	}
	tests := map[string]struct {
		args args
		want []optgen.Tag
	}{
		// "foo": {
		// 	args: args{sourcePath: "./testfile/foo.sample", structName: "Thing"},
		// 	want: []optgen.Tag{field1, field2},
		// },
		// "foo with defined tags": {
		// 	args: args{sourcePath: "./testfile/foo.sample", optgen.Tag: "opt", structName: "Thing"},
		// 	want: []optgen.Tag{field1, field2},
		// },
		// "fooJerax": {
		// 	args: args{sourcePath: "./testfile/fooJerax.sample", optgen.Tag: "jerax", structName: "Thing"},
		// 	want: []optgen.Tag{field1, field2},
		// },
		// "fooJerax want dendy": {
		// 	args: args{sourcePath: "./testfile/fooJerax.sample", optgen.Tag: "dendy", structName: "Thing"},
		// 	want: []optgen.Tag{field3},
		// },
		// "missing struct": {
		// 	args: args{sourcePath: "./testfile/fooJerax.sample", optgen.Tag: "jerax", structName: "Thingy"},
		// 	want: []optgen.Tag{},
		// },
		"All": {
			args: args{sourcePath: "./testfile/foo.sample", Tag: "_all_", structName: "Thing"},
			want: []optgen.Tag{field1, field2, field3},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f, err := getTokenFileset(tt.args.sourcePath)
			assert.NoError(t, err)
			got := optgen.GetTags(f, tt.args.structName, tt.args.Tag)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetOptTags() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}