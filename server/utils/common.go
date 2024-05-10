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
	// 将所有 <br> 和 <br/> 标签转换为换行符 \n
	brReplacer, err := regexp.Compile(`(?i)<br\s*/?>`)
	if err != nil {
		panic(err) // Handle regexp compilation error for <br> tags.
	}
	input = brReplacer.ReplaceAllString(input, "\n")

	// 编译正则表达式以匹配所有HTML标签
	tagReplacer, err := regexp.Compile("<.*?>")
	if err != nil {
		panic(err) // Handle regexp compilation error for general HTML tags.
	}
	return tagReplacer.ReplaceAllString(input, "")
}
