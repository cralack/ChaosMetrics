package utils

import (
	"regexp"
)

func RemoveExtraLF(str string) string {
	for i, char := range str {
		if char == 10 {
			return str[:i]
		}
	}
	return str
}

func RemoveHTMLTags(input string) string {
	r, err := regexp.Compile("<.*?>")
	if err != nil {
		// Handle regexp compilation error.
		panic(err)
	}
	return r.ReplaceAllString(input, "")
}
