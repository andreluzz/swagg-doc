package parser

import (
	"strings"
)

func parseController(lines []string, gf *gofile) error {
	var builder strings.Builder
	index := 0
	for index < len(lines)-1 {
		componentLabel := lines[index+1]
		for i := index + 2; i < len(lines); i++ {
			line := lines[i]
			index = i
			if line == "" {
				break
			}
			space := ""
			if componentLabel != "tags:" {
				space = "  "
			}
			if _, err := builder.WriteString(space + line + "\n"); err != nil {
				return err
			}
		}
		if componentLabel == "tags:" {
			gf.controller.tags = builder.String()
		} else {
			gf.globalComponents[componentLabel] = builder.String()
		}
		builder.Reset()
	}

	return nil
}
