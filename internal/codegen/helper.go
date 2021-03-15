package codegen

import (
	"go.sia.tech/encodegen/internal/toolbox"
	"sort"
	"strconv"
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

func getNextIterator(currentIdentifier string) string {
	// this function allows the generated the code to iterate over slices of structs that have slices within them without having iteration identifiers conflict (i.e., there'd be multiple "range i := r.Fields"s)
	idSplit := strings.Split(currentIdentifier, "i")
	if len(idSplit) != 2 {
		return "i"
	}
	if idSplit[1] != "" {
		num, err := strconv.Atoi(idSplit[1])
		if err != nil {
			return "i"
		}
		return "i" + strconv.Itoa(num+1)
	} else {
		return "i1"
	}
}

func fieldsHaveSlice(fields []*toolbox.FieldInfo) bool {
	for _, field := range fields {
		if field.IsSlice {
			return true
		}
		if len(field.AnonymousChildFields) > 0 {
			hasSlice := fieldsHaveSlice(field.AnonymousChildFields)
			if hasSlice {
				return hasSlice
			}
		}
	}
	return false
}
