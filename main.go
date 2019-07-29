package main

import (
	"flag"
	"fmt"

	"github.com/andreluzz/swagg-doc/parser"
)

var pkg = flag.String("package", "github.com/andreluzz/swagg-doc/mock/api", "The path to the application main package")

func main() {
	flag.Parse()
	p, err := parser.New()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := p.Process(*pkg, false); err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := p.YAML(); err != nil {
		fmt.Println(err.Error())
		return
	}
}
