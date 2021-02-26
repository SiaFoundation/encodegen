package codegen

import (
	"strings"
)

func countLeft(line string, prefix string) int {
	return len(line) - len(strings.TrimLeft(line, prefix))
}

func countRight(line string, prefix string) int {
	return len(line) - len(strings.TrimRight(line, prefix))
}
