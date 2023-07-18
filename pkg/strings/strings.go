package strings

import (
	_ "errors"
	_ "reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	pl "github.com/gertd/go-pluralize"
	_ "github.com/spf13/cast"
)

type stringBuilder = strings.Builder

// func String(s any) (string, error) {
// 	v := reflect.ValueOf(s)
// 	switch v.Kind() {
// 	case reflect.String:
// 		return cast.ToString(s), nil
// 	default:
// 		return "", errors.New("parse error")
// 	}
// }

func ToCamel(s string) string {

	// s, err := String(i)
	// if err != nil {
	// 	return err.Error()
	// }

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

	// s, err := String(i)
	// if err != nil {
	// 	return err.Error()
	// }

	if s == "" {
		return s
	}

	return strings.ToLower(strings.ReplaceAll(s, "-", ""))
}

// Converts a string to Pascal case
func ToPascal(s string) string {

	// s, err := String(i)
	// if err != nil {
	// 	return err.Error()
	// }

	if s == "" {
		return s
	}

	if strings.Contains(s, "-") {
		buf := &stringBuilder{}
		words := strings.Split(s, "-")
		for _, word := range words {
			buf.WriteString(strings.Title(word))
		}
		return buf.String()
	}

	return strings.Title(s)
}

func ToPlural(s string) string {

	// s, err := String(i)
	// if err != nil {
	// 	return err.Error()
	// }

	var singular string = s
	var prefix string = ""
	if strings.Contains(s, "-") {
		parts := strings.SplitAfterN(s, "-", -1)
		max := len(parts) - 1
		singular = parts[max]
		for i, v := range parts {
			if i < max {
				prefix = prefix + v
			}
		}
	}

	plural := pl.NewClient().Plural(singular)

	return prefix + plural
}

// Formats string into acceptable go package name
func ToSnakeCase(s string) string {

	// s, err := String(i)
	// if err != nil {
	// 	return err.Error()
	// }

	if s == "" {
		return s
	}

	return strings.ReplaceAll(s, "-", "_")
}
