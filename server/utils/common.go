package utils

func RemoveExtraLF(str string) string {
	for i, char := range str {
		if char == 10 {
			return str[:i]
		}
	}
	return str
}
