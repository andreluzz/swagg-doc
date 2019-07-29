package parser

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// YAML return the YAML string of this parse
func (p *Parser) YAML() error {
	var builder strings.Builder

	if _, ok := p.files["api"]; ok {
		for _, val := range p.files["api"] {
			builder.WriteString(val.definitions + "\n")
			break
		}
	} else {
		return fmt.Errorf("yaml without swagg-doc:api")
	}

	builder.WriteString("paths:\n")

	for _, gf := range p.files["controller"] {
		for _, endpoint := range gf.controller.endpoints {
			builder.WriteString("  " + endpoint.path + ":\n")
			for key, value := range endpoint.operations {
				builder.WriteString("    " + strings.ToLower(key) + ":\n")
				builder.WriteString(value)
			}
		}
	}

	builder.WriteString("tags:\n")

	for _, gf := range p.files["controller"] {
		builder.WriteString(gf.controller.tags)
	}

	for _, mapFiles := range p.files {
		for _, gf := range mapFiles {
			for k, v := range gf.globalComponents {
				p.globalComponents[k] = p.globalComponents[k] + v
			}
		}
	}

	builder.WriteString("\ncomponents:\n")

	for componentType, componentValue := range p.globalComponents {
		builder.WriteString("  " + componentType + "\n")
		builder.WriteString(componentValue)
	}

	return ioutil.WriteFile("swagger.yaml", []byte(builder.String()), 0644)
}
