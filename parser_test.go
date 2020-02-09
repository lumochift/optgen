package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
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
		field1 = Tag{
			Name:     "Field1",
			DataType: "string",
		}

		field2 = Tag{
			Name:     "Field2",
			DataType: "[]*int",
		}

		field3 = Tag{Name: "Field3",
			DataType: "map[byte]float64"}
	)
	type args struct {
		sourcePath string
		tag        string
		structName string
	}
	tests := map[string]struct {
		args args
		want []Tag
	}{
		"foo": {
			args: args{sourcePath: "./testfile/foo.sample", structName: "Thing"},
			want: []Tag{field1, field2},
		},
		"foo with defined tags": {
			args: args{sourcePath: "./testfile/foo.sample", tag: "opt", structName: "Thing"},
			want: []Tag{field1, field2},
		},
		"fooJerax": {
			args: args{sourcePath: "./testfile/fooJerax.sample", tag: "jerax", structName: "Thing"},
			want: []Tag{field1, field2},
		},
		"fooJerax want dendy": {
			args: args{sourcePath: "./testfile/fooJerax.sample", tag: "dendy", structName: "Thing"},
			want: []Tag{field3},
		},
		"missing struct": {
			args: args{sourcePath: "./testfile/fooJerax.sample", tag: "jerax", structName: "Thingy"},
			want: []Tag{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f, err := getTokenFileset(tt.args.sourcePath)
			assert.NoError(t, err)
			got := GetTags(f, tt.args.structName, tt.args.tag)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetOptTags() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
