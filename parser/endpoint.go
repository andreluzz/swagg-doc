package parser

import "strings"

func parseEndpoint(lines []string, gf *gofile) error {
	var builder strings.Builder
	fields := strings.Fields(lines[0])

	for i := 1; i < len(lines); i++ {
		if _, err := builder.WriteString("      " + lines[i] + "\n"); err != nil {
			return err
		}
	}

	endpointPath := fields[3]
	index := getEndpointIndexByPath(endpointPath, gf.controller.endpoints)
	if index != -1 {
		gf.controller.endpoints[index].operations[fields[2]] = builder.String()
	} else {
		ep := endpoint{
			path:       endpointPath,
			operations: make(map[string]string),
		}
		ep.operations[fields[2]] = builder.String()
		gf.controller.endpoints = append(gf.controller.endpoints, ep)
	}

	return nil
}

func getEndpointIndexByPath(path string, endpoints []endpoint) int {
	for i, endpoint := range endpoints {
		if endpoint.path == path {
			return i
		}
	}
	return -1
}
