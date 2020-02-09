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
	type args struct {
		sourcePath string
		tag        string
	}
	tests := map[string]struct {
		args args
		want []string
	}{
		"foo": {
			args: args{sourcePath: "./testfile/foo.sample"},
			want: []string{"Field1", "Field2"},
		},
		"foo with defined tags": {
			args: args{sourcePath: "./testfile/foo.sample", tag: "opt"},
			want: []string{"Field1", "Field2"},
		},
		"fooJerax": {
			args: args{sourcePath: "./testfile/fooJerax.sample", tag: "jerax"},
			want: []string{"Field1", "Field2"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f, err := getTokenFileset(tt.args.sourcePath)
			assert.NoError(t, err)
			got := GetTags(f, tt.args.tag)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetOptTags() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
