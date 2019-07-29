package parser

import (
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/andreluzz/sandbox/swagg-doc/utils"
)

// Parser process the main package and generate the doc
type Parser struct {
	files            map[string]map[string]*gofile
	globalComponents map[string]string
	goPath           string
	processedFiles   map[string]bool
}

type gofile struct {
	filepath         string
	definitions      string
	controller       *controller
	globalComponents map[string]string
}

type controller struct {
	endpoints []endpoint
	tags      string
}

type endpoint struct {
	path       string
	operations map[string]string
}

// New initializate and returns a new Parser
func New() (*Parser, error) {
	gopath, err := utils.GetGOPATH()
	if err != nil {
		return nil, err
	}

	p := &Parser{
		goPath:           gopath,
		files:            make(map[string]map[string]*gofile),
		globalComponents: make(map[string]string),
		processedFiles:   make(map[string]bool),
	}

	return p, nil
}

// Process analyse go code for comments.
// Only commnets where the first line has swagg-doc willl be processed.
func (p *Parser) Process(packagePath string, externalPackage bool) error {
	packageFullPath := fmt.Sprintf("%s/src/%s", p.goPath, packagePath)
	err := filepath.Walk(packageFullPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if _, processed := p.processedFiles[path]; !processed && !info.IsDir() && filepath.Ext(info.Name()) == ".go" {
			p.processedFiles[path] = true
			if err := parseFileComments(path, p, packagePath, externalPackage); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func parseFileComments(filePath string, parser *Parser, packagePath string, externalPackage bool) error {
	fileSet := token.NewFileSet()
	fileTree, err := goparser.ParseFile(fileSet, filePath, nil, goparser.ParseComments)
	if err != nil {
		return err
	}

	// Get model definitions
	imports := make(map[string]string)

	for _, astImport := range fileTree.Imports {
		importedPackageName := strings.Trim(astImport.Path.Value, "\"")
		fmt.Println(importedPackageName)

		var importedPackageAlias string
		if astImport.Name != nil && astImport.Name.Name != "." && astImport.Name.Name != "_" {
			importedPackageAlias = astImport.Name.Name
		} else {
			importPath := strings.Split(importedPackageName, "/")
			importedPackageAlias = importPath[len(importPath)-1]
		}
		fmt.Println(importedPackageAlias)

		imports[importedPackageAlias] = importedPackageName
	}

	// Get model definitions
	models := make(map[string]*ast.TypeSpec)

	for _, astDeclaration := range fileTree.Decls {
		if generalDeclaration, ok := astDeclaration.(*ast.GenDecl); ok && generalDeclaration.Tok == token.TYPE {
			for _, astSpec := range generalDeclaration.Specs {
				if typeSpec, ok := astSpec.(*ast.TypeSpec); ok {
					models[typeSpec.Name.String()] = typeSpec
				}
			}
		}
	}

	gf := &gofile{
		filepath:         filePath,
		globalComponents: make(map[string]string),
	}
	gfAction := "undefined"

	for _, comment := range fileTree.Comments {
		lines := strings.Split(comment.Text(), "\n")
		valid, action := parseSwaggDocCommand(lines[0])
		if valid {
			if gfAction == "undefined" {
				gfAction = action
			}
			switch action {
			case "api":
				if err := parseAPI(lines, gf); err != nil {
					return err
				}
			case "controller":
				if gf.controller != nil {
					return fmt.Errorf("only one swagg-doc:controller can be defined per page")
				}
				gf.controller = &controller{}
				if err := parseController(lines, gf); err != nil {
					return err
				}
			case "endpoint":
				if gf.controller == nil {
					return fmt.Errorf("swagg-doc:controller should be defined before swagg-doc:endpoint")
				}
				if err := parseEndpoint(lines, gf); err != nil {
					return err
				}
			case "model":
				if parser.files["model"][filePath] != nil {
					return nil
				}
				fields := strings.Fields(lines[0])
				if typeSpec, ok := models[fields[0]]; ok {
					if err := parseModel(lines, gf, parser, typeSpec, imports, packagePath, externalPackage); err != nil {
						return err
					}
				}
			}
		}
	}

	if parser.files[gfAction] == nil {
		parser.files[gfAction] = make(map[string]*gofile)
	}
	parser.files[gfAction][filePath] = gf

	return nil
}

func parseSwaggDocCommand(line string) (bool, string) {
	fields := strings.Fields(line)
	for _, field := range fields {
		split := strings.Split(field, ":")
		if len(split) >= 2 && split[0] == "swagg-doc" {
			return true, split[1]
		}
	}
	return false, ""
}
