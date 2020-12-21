package strings

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

type stringBuilder = strings.Builder

func ToCamel(s string) string {
	if s == "" {
		return s
	}

	s = ToPascal(s)
	r, size := utf8.DecodeRuneInString(s)

	buf := &stringBuilder{}
	buf.WriteRune(unicode.ToLower(r))
	buf.WriteString(s[size:])
	return buf.String()
}

// Formats string into acceptable go package name
func ToPackage(s string) string {
	if s == "" {
		return s
	}

	return strings.ReplaceAll(s, "-", "")
}

// Converts a string to Pascal case
func ToPascal(s string) string {
	if s == "" {
		return s
	}

	a := regexp.MustCompile(`-`)

	buf := &stringBuilder{}
	words := a.Split(s, -1)
	for _, word := range words {
		buf.WriteString(strings.Title(word))
	}
	return buf.String()
}

// Formats string into acceptable go package name
func ToSnakeCase(s string) string {
	if s == "" {
		return s
	}

	return strings.ReplaceAll(s, "-", "_")
}
