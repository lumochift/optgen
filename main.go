package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/lumochift/optgen/generator"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func titleCase(s string) string {
	return cases.Title(language.Und).String(s)
}

var (
	sourceFile, tagName, structName   string
	writeMode, allFields, showVersion bool

	version string = "dev"

	funcMap = template.FuncMap{
		"title":   titleCase,
		"toLower": strings.ToLower,
	}
)

func initCLI() {
	flag.StringVar(&sourceFile, "file", "", "path file")
	flag.StringVar(&tagName, "tag", "opt", "custom tag")
	flag.StringVar(&structName, "name", "", "struct name")
	flag.BoolVar(&writeMode, "w", false, "enable write mode")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&allFields, "all", false, "generate all fields")
	flag.Parse()
}

func main() {
	initCLI()

	if showVersion {
		fmt.Printf("version: %s", version)

		return
	}

	const exampleMessage = "e.g: optgen -file sample-file.go -name Thing"
	if sourceFile == "" || structName == "" {
		log.Fatalf("Source file and struct name must be provided.\n%s", exampleMessage)
	}

	source, err := os.ReadFile(sourceFile)
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
		Tags       []generator.Tag
	}

	if allFields {
		tagName = "_all_"
	}

	tags := generator.GetTags(f, structName, tagName)

	tmpl := template.Must(template.New("generator").Funcs(funcMap).Parse(generator.CodeTemplate))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, TemplateData{
		Tags:       tags,
		StructName: structName,
	}); err != nil {
		log.Panic(err)
	}

	if !writeMode {
		fmt.Println(buf.String())
		return
	}
	var sourceBuf bytes.Buffer
	printer.Fprint(&sourceBuf, fset, f)
	if err := os.WriteFile(sourceFile, append(sourceBuf.Bytes(), buf.Bytes()...), 0644); err != nil {
		log.Println(err)
	}
	// format generated code
	exec.Command("gofmt", "-s", "-w", sourceFile)
}
