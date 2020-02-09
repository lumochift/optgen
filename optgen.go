package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"strings"
	"text/template"
)

var (
	sourceFile, tagName, structName string
	writeMode                       bool

	funcMap = template.FuncMap{
		"title":   strings.Title,
		"toLower": strings.ToLower,
	}
)

func initCLI() {
	flag.StringVar(&sourceFile, "file", "", "path file")
	flag.StringVar(&tagName, "tag", "opt", "custom tag")
	flag.StringVar(&structName, "name", "", "struct name")
	flag.BoolVar(&writeMode, "w", false, "enable write mode")
	flag.Parse()
	const exampleMessage = "e.g: optgen -file sample-file.go -name Thing"
	if sourceFile == "" || structName == "" {
		log.Fatalf("Source file and struct name must be provided.\n%s", exampleMessage)
	}
}

func main() {
	initCLI()
	source, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		log.Fatal(err)
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", source, 0)

	if err != nil {
		log.Fatal(err)
	}

	type TemplateData struct {
		StructName string
		Tags       []Tag
	}

	tags := GetTags(f, structName, tagName)

	tmpl := template.Must(template.New("optgen").Funcs(funcMap).Parse(optTemplate))
	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, TemplateData{
		Tags:       tags,
		StructName: structName,
	}); err != nil {
		log.Fatal(err)
	}
	if !writeMode {
		fmt.Println(buf.String())
		return
	}

}
