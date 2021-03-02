package codegen

import (
	"sort"
	"strings"
)

func firstLetterToUppercase(text string) string {
	return strings.ToUpper(string(text[0:1])) + string(text[1:])
}

func firstLetterToLowercase(text string) string {
	return strings.ToLower(string(text[0:1])) + string(text[1:])
}

func extractReceiverAlias(structType string) string {
	var result = string(structType[0])
	for i := len(structType) - 1; i > 0; i-- {
		aChar := string(structType[i])
		lowerChar := strings.ToLower(aChar)
		if lowerChar != aChar {
			result = lowerChar
			break
		}
	}
	return strings.ToLower(result)
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}
