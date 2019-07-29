package parser

import (
	"fmt"
	"go/ast"
	"reflect"
	"strings"
)

func parseModel(lines []string, gf *gofile, parser *Parser, typeSpec *ast.TypeSpec, imports map[string]string, packagePath string, externalPackage bool) error {
	var builder strings.Builder

	space := "  "

	builder.WriteString(strings.Repeat(space, 2))
	if externalPackage {
		builder.WriteString(strings.ReplaceAll(packagePath, "/", "-"))
		builder.WriteString("-")
	}
	builder.WriteString(typeSpec.Name.String())
	builder.WriteString(":\n")
	builder.WriteString(strings.Repeat(space, 3))
	builder.WriteString("type: object\n")
	builder.WriteString(strings.Repeat(space, 3))
	builder.WriteString("properties:\n")
	if astStructType, ok := typeSpec.Type.(*ast.StructType); ok {
		for _, field := range astStructType.Fields.List {
			structTag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
			jsonName := strings.Split(structTag.Get("json"), ",")[0]
			builder.WriteString(strings.Repeat(space, 4))
			builder.WriteString(jsonName)
			builder.WriteString(":\n")
			builder.WriteString(strings.Repeat(space, 5))

			fieldTypeStr := getTypeAsString(field.Type)
			fieldType := strings.Split(fieldTypeStr, ".")
			if len(fieldType) > 1 {
				importPath := imports[fieldType[0]]
				builder.WriteString("$ref: '#/components/schemas/")
				builder.WriteString(strings.ReplaceAll(importPath, "/", "-"))
				builder.WriteString("-")
				builder.WriteString(fieldType[1])
				builder.WriteString("'")
				if err := parser.Process(importPath, true); err != nil {
					return err
				}
			} else {
				builder.WriteString("type: ")
				builder.WriteString(fieldTypeStr)
			}
			builder.WriteString("\n")
		}
	}

	gf.globalComponents["schemas:"] += builder.String()

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
			// log.Printf("Get type as string (star expression)! %#v, type: %s\n", astStarExpr.X, fmt.Sprint(astStarExpr.X))
		} else if astSelectorExpr, ok := fieldType.(*ast.SelectorExpr); ok {
			packageNameIdent, _ := astSelectorExpr.X.(*ast.Ident)
			realType = packageNameIdent.Name + "." + astSelectorExpr.Sel.Name
			//log.Printf("Get type as string(selector expression)! X: %#v , Sel: %#v, type %s\n", astSelectorExpr.X, astSelectorExpr.Sel, realType)
			if realType == "time.Time" {
				realType = "string"
			}
		} else {
			// log.Printf("Get type as string(no star expression)! %#v , type: %s\n", fieldType, fmt.Sprint(fieldType))
			realType = fmt.Sprint(fieldType)
			if realType == "int" {
				realType = "integer"
			}
		}
	}
	return realType
}
