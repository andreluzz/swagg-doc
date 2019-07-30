package parser

import (
	"fmt"
	"go/ast"
	"reflect"
	"strings"
)

func parseStructDefinition(lines []string, gf *gofile, parser *Parser, typeSpec *ast.TypeSpec, imports map[string]string, packagePath string, externalPackage bool) error {

	// process comment attribute definitions
	attributes := make(map[string]string, len(lines))
	ignore := make(map[string]bool, len(lines))
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 1 && strings.Contains(fields[1], "swagg-doc:attribute") {
			params := strings.Split(fields[1], ":")
			if len(params) < 3 {
				return fmt.Errorf("invalid swagg-doc:attribute tag | %s", gf.filepath)
			}
			if params[2] == "ignore_write" {
				ignore[fields[0]] = true
			} else {
				attributes[fields[0]] = params[2]
			}
		}
	}

	fields := strings.Fields(lines[0])
	isModuleInterface := len(fields) > 1 && strings.HasPrefix(fields[1], "swagg-doc:model:interface")

	if err := createSchemaModule("", lines, gf, parser, typeSpec, imports, packagePath, externalPackage, ignore, attributes, isModuleInterface); err != nil {
		return err
	}
	if !isModuleInterface {
		if err := createSchemaModule("Update", lines, gf, parser, typeSpec, imports, packagePath, externalPackage, ignore, attributes, isModuleInterface); err != nil {
			return err
		}
		if err := createSchemaModule("Create", lines, gf, parser, typeSpec, imports, packagePath, externalPackage, ignore, attributes, isModuleInterface); err != nil {
			return err
		}
	}

	return nil
}

func createSchemaModule(sufix string, lines []string, gf *gofile, parser *Parser, typeSpec *ast.TypeSpec, imports map[string]string, packagePath string, externalPackage bool, ignore map[string]bool, attributes map[string]string, isModelInterface bool) error {
	var builder strings.Builder
	modelPath := typeSpec.Name.String()
	if externalPackage {
		modelPath = fmt.Sprintf("%s-%s", strings.ReplaceAll(packagePath, "/", "-"), modelPath)
	}
	parser.modelsAlias[modelPath] = getAlias(lines[0])

	space := "  "

	builder.WriteString(strings.Repeat(space, 2))
	builder.WriteString(modelPath)
	builder.WriteString(sufix)
	builder.WriteString(":\n")
	builder.WriteString(strings.Repeat(space, 3))
	builder.WriteString("type: object\n")
	builder.WriteString(strings.Repeat(space, 3))
	builder.WriteString("properties:\n")
	attRequired := []string{}
	if astStructType, ok := typeSpec.Type.(*ast.StructType); ok {
		for _, field := range astStructType.Fields.List {
			structTag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
			jsonName := strings.Split(structTag.Get("json"), ",")[0]
			if strings.Contains(structTag.Get("validate"), "required") {
				attRequired = append(attRequired, jsonName)
			}
			if _, match := ignore[jsonName]; match && sufix != "" {
				continue
			}
			if sufix == "Update" && structTag.Get("updatable") == "false" {
				continue
			}
			if structTag.Get("pk") == "true" && sufix != "" {
				continue
			}
			builder.WriteString(strings.Repeat(space, 4))
			builder.WriteString(jsonName)
			builder.WriteString(":\n")
			builder.WriteString(strings.Repeat(space, 5))

			if val, ok := attributes[jsonName]; ok {
				builder.WriteString("type: ")
				builder.WriteString(val)
				builder.WriteString("\n")
				continue
			}

			fieldTypeStr := getTypeAsString(field.Type)
			fieldType := strings.Split(fieldTypeStr, ".")
			if len(fieldType) > 1 {
				importPath := imports[fieldType[0]]
				builder.WriteString("$ref: '#/components/schemas/")
				builder.WriteString(strings.ReplaceAll(importPath, "/", "-"))
				builder.WriteString("-")
				builder.WriteString(fieldType[1])
				builder.WriteString(sufix)
				builder.WriteString("'")
				if err := parser.Process(importPath, true); err != nil {
					return err
				}
			} else {
				if isStruct(fieldTypeStr) {
					builder.WriteString("$ref: '#/components/schemas/")
					if externalPackage {
						builder.WriteString(strings.ReplaceAll(packagePath, "/", "-"))
						builder.WriteString("-")
					}
					builder.WriteString(fieldTypeStr)
					builder.WriteString(sufix)
					builder.WriteString("'")
				} else {
					builder.WriteString("type: ")
					builder.WriteString(fieldTypeStr)
				}
			}
			builder.WriteString("\n")
		}
	}
	if len(attRequired) > 0 && sufix != "Update" {
		builder.WriteString(strings.Repeat(space, 3))
		builder.WriteString("required:\n")
		for _, att := range attRequired {
			builder.WriteString(strings.Repeat(space, 4))
			builder.WriteString("- ")
			builder.WriteString(att)
			builder.WriteString("\n")
		}
	}

	if isModelInterface {
		alias := parser.modelsAlias[modelPath]
		if _, ok := parser.interfaces[alias]; ok {
			return fmt.Errorf("interface alias (%s) already exists", alias)
		}
		parser.interfaces[alias] = builder.String()
	} else {
		gf.globalComponents["schemas:"] += builder.String()
	}

	return nil
}

func getTypeAsString(fieldType interface{}) string {
	var realType string
	if astArrayType, ok := fieldType.(*ast.ArrayType); ok {
		// log.Printf("arrayType: %#v\n", astArrayType)
		realType = fmt.Sprintf("[]%v", getTypeAsString(astArrayType.Elt))
	} else if astMapType, ok := fieldType.(*ast.MapType); ok {
		// log.Printf("arrayType: %#v\n", astArrayType)
		realType = fmt.Sprintf("[]%v", getTypeAsString(astMapType.Value))
	} else if _, ok := fieldType.(*ast.InterfaceType); ok {
		realType = "interface"
	} else {
		if astStarExpr, ok := fieldType.(*ast.StarExpr); ok {
			realType = fmt.Sprint(astStarExpr.X)
			//log.Printf("Get type as string (star expression)! %#v, type: %s\n", astStarExpr.X, fmt.Sprint(astStarExpr.X))
		} else if astSelectorExpr, ok := fieldType.(*ast.SelectorExpr); ok {
			packageNameIdent, _ := astSelectorExpr.X.(*ast.Ident)
			realType = packageNameIdent.Name + "." + astSelectorExpr.Sel.Name
			//log.Printf("Get type as string(selector expression)! X: %#v , Sel: %#v, type %s\n", astSelectorExpr.X, astSelectorExpr.Sel, realType)
			if realType == "time.Time" {
				realType = "string"
			}
		} else {
			//log.Printf("Get type as string(no star expression)! %#v , type: %s\n", fieldType, fmt.Sprint(fieldType))
			realType = fmt.Sprint(fieldType)
			if realType == "int" {
				realType = "integer"
			}
		}
	}
	return realType
}

func isStruct(t string) bool {
	return t != "string" && t != "integer" && t != "bool"
}

func getAlias(line string) string {
	fields := strings.Fields(line)
	defs := strings.Split(fields[1], ":")
	if len(defs) > 3 {
		return defs[3]
	}
	if len(defs) > 2 && strings.ToLower(defs[2]) != "interface" {
		return defs[2]
	}
	return fields[0]
}
