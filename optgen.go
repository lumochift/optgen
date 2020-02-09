package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
)

func main() {
	var sourceFile = flag.String("file", "", "path file")
	flag.Parse()
	source, err := ioutil.ReadFile(*sourceFile)
	if err != nil {
		log.Fatal(err)
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", source, 0)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(GetTags(f, ""))

	GetGeneratedExpression(nil, "")
}
