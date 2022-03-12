package utils

import "strings"

func CorrectCase(s string) string {
	if strings.HasSuffix(s, "ID") {
		s = strings.TrimSuffix(s, "ID")
		return s + "Id"
	}
	return s
}

func SnakeToCamel(word string) string {
	words := strings.Split(word, "_")
	words[0] = strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		words[i] = strings.ToUpper(string(words[i][0])) + strings.TrimPrefix(words[i], string(words[i][0]))
	}
	return strings.Join(words, "")
}

func SnakeToPascal(word string) string {
	return "hello"
}
