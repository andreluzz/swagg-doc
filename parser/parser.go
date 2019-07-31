package parser

import (
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/andreluzz/swagg-doc/utils"
)

// Parser process the main package and generate the doc
type Parser struct {
	apiPackage       string
	files            map[string]map[string]*gofile
	globalComponents map[string]string
	goPath           string
	processedFiles   map[string]bool
	interfaces       map[string]string
	modelsAlias      map[string]string
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
func New(mainPackage string) (*Parser, error) {
	gopath, err := utils.GetGOPATH()
	if err != nil {
		return nil, err
	}

	p := &Parser{
		goPath:           gopath,
		apiPackage:       mainPackage,
		files:            make(map[string]map[string]*gofile),
		globalComponents: make(map[string]string),
		interfaces:       make(map[string]string),
		processedFiles:   make(map[string]bool),
		modelsAlias:      make(map[string]string),
	}

	return p, nil
}

// Process analyse go code for comments.
// Only commnets where the first line has swagg-doc willl be processed.
func (p *Parser) Process(packagePath, importScope string, externalPackage bool) error {
	packageFullPath := fmt.Sprintf("%s/src/%s", p.goPath, packagePath)
	err := filepath.Walk(packageFullPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if _, processed := p.processedFiles[path]; !processed && !info.IsDir() && filepath.Ext(info.Name()) == ".go" {
			p.processedFiles[path] = true
			if err := parseFileComments(path, p, packagePath, importScope, externalPackage); err != nil {
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

func parseFileComments(filePath string, parser *Parser, packagePath, importScope string, externalPackage bool) error {
	fileSet := token.NewFileSet()
	astFile, err := goparser.ParseFile(fileSet, filePath, nil, goparser.ParseComments)
	if err != nil {
		return err
	}

	// Get model imports
	imports := make(map[string]string)

	for _, astImport := range astFile.Imports {
		importedPackageName := strings.Trim(astImport.Path.Value, "\"")

		var importedPackageAlias string
		if astImport.Name != nil && astImport.Name.Name != "." && astImport.Name.Name != "_" {
			importedPackageAlias = astImport.Name.Name
		} else {
			importPath := strings.Split(importedPackageName, "/")
			importedPackageAlias = importPath[len(importPath)-1]
		}
		if !strings.Contains(importedPackageName, parser.apiPackage) {
			imports[importedPackageAlias] = importedPackageName
		}
	}

	// Get structs definitions
	structsDefinitions := make(map[string]*ast.TypeSpec)

	for _, astDeclaration := range astFile.Decls {
		if generalDeclaration, ok := astDeclaration.(*ast.GenDecl); ok && generalDeclaration.Tok == token.TYPE {
			for _, astSpec := range generalDeclaration.Specs {
				if typeSpec, ok := astSpec.(*ast.TypeSpec); ok {
					structsDefinitions[typeSpec.Name.String()] = typeSpec
				}
			}
		}
	}

	gf := &gofile{
		filepath:         filePath,
		globalComponents: make(map[string]string),
	}
	gfAction := "undefined"

	for _, comment := range astFile.Comments {
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
				for _, path := range imports {
					scopes := strings.Split(importScope, ",")
					for _, scope := range scopes {
						if strings.Contains(path, scope) {
							if err := parser.Process(path, importScope, true); err != nil {
								return err
							}
						}
					}
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
				fields := strings.Fields(lines[0])
				if typeSpec, ok := structsDefinitions[fields[0]]; ok {
					if err := parseStructDefinition(lines, gf, parser, typeSpec, imports, packagePath, importScope, externalPackage); err != nil {
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
