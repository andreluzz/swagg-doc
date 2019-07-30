package parser

import (
	"fmt"
	"io/ioutil"
	"regexp"
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

	schemaInterfaces := make(map[string]string)

	for _, gf := range p.files["controller"] {
		for _, endpoint := range gf.controller.endpoints {
			builder.WriteString("  " + endpoint.path + ":\n")
			for key, value := range endpoint.operations {
				builder.WriteString("    " + strings.ToLower(key) + ":\n")
				endpoint, endpointInterfaces := parseInterfaces(value)
				for k, v := range endpointInterfaces {
					schemaInterfaces[k] = v
				}
				builder.WriteString(endpoint)
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

	for k, v := range schemaInterfaces {
		if interfaceModel := generateInterfaceModel(k, v, p.interfaces); interfaceModel != "" {
			p.globalComponents["schemas:"] = p.globalComponents["schemas:"] + interfaceModel
			continue
		}
		return fmt.Errorf("invalid controller interface alias | %s", strings.Split(k, ":")[0])
	}

	processedModelAlias := make(map[string]bool)

	for path, alias := range p.modelsAlias {
		if processedModelAlias[alias] {
			return fmt.Errorf("duplicated model alias | %s - %s", path, alias)
		}
		p.globalComponents["schemas:"] = strings.ReplaceAll(p.globalComponents["schemas:"], path, alias)
		processedModelAlias[alias] = true
	}

	builder.WriteString("\ncomponents:\n")

	for componentType, componentValue := range p.globalComponents {
		builder.WriteString("  " + componentType + "\n")
		builder.WriteString(componentValue)
	}

	return ioutil.WriteFile("swagger.yaml", []byte(builder.String()), 0644)
}

func parseInterfaces(endpoint string) (string, map[string]string) {
	re := regexp.MustCompile("swagg-doc:interface:(.+)")
	interfaces := re.FindAllString(endpoint, -1)
	schemaModels := make(map[string]string)
	for _, match := range interfaces {
		parsed := strings.Replace(match[20:], ":", "", -1)
		endpoint = strings.ReplaceAll(endpoint, match, parsed)
		schemaModels[match[20:len(match)-1]] = parsed[:len(parsed)-1]
	}
	return endpoint, schemaModels
}

func generateInterfaceModel(definition, name string, interfaces map[string]string) string {
	def := strings.Split(definition, ":")
	model, ok := interfaces[def[0]]
	if !ok {
		return ""
	}
	lines := strings.Split(model, "\n")
	var builder strings.Builder
	space := "  "
	builder.WriteString(strings.Repeat(space, 2))
	builder.WriteString(name)
	builder.WriteString(":\n")
	for index := 1; index < len(lines); index++ {
		line := lines[index]
		if strings.Contains(line, "type: payload") {
			modelRef := fmt.Sprintf("$ref: '#/components/schemas/%s'", def[1])
			if len(def) > 2 && def[2] == "Array" {
				builder.WriteString(strings.Replace(line, "type: payload", "type: array\n", -1))
				builder.WriteString(strings.Repeat(space, 5))
				builder.WriteString("items:\n")
				builder.WriteString(strings.Repeat(space, 6))
				builder.WriteString(modelRef)
			} else {
				builder.WriteString(strings.Replace(line, "type: payload", modelRef, -1))
			}
		} else {
			builder.WriteString(line)
		}
		if index < len(lines)-1 {
			builder.WriteString("\n")
		}
	}
	return builder.String()
}
