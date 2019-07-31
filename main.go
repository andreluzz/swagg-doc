package main

import (
	"flag"
	"fmt"

	"github.com/andreluzz/swagg-doc/parser"
)

var pkg = flag.String("package", "github.com/andreluzz/swagg-doc/mock/api", "The path to the application main package")
var scope = flag.String("scope", "swagg-doc/mock", "The scope to the imports that shoul be parsed. Use ',' to define multiple values.")

func main() {
	flag.Parse()
	p, err := parser.New(*pkg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := p.Process(*pkg, *scope, false); err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := p.YAML(); err != nil {
		fmt.Println(err.Error())
		return
	}
}
