package parser

import (
	"strings"
)

func parseAPI(lines []string, gf *gofile) error {
	var builder strings.Builder
	index := 0
	for i := 1; i < len(lines); i++ {
		line := lines[i]
		if line == "" {
			index = i
			break
		}
		if _, err := builder.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	gf.definitions = builder.String()

	for index < len(lines)-1 {
		builder.Reset()
		componentLabel := lines[index+1]
		for i := index + 2; i < len(lines); i++ {
			line := lines[i]
			index = i
			if line == "" {
				break
			}
			if _, err := builder.WriteString("  " + line + "\n"); err != nil {
				return err
			}
		}
		gf.globalComponents[componentLabel] = builder.String()
	}

	return nil
}
