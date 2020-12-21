package strings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToCamel(t *testing.T) {
	str1 := "test-string"
	str2 := "testString"

	strout := ToCamel(str1)
	assert.Equal(t, strout, str2, "Expected string to be formatted in camelCase")

}

func TestToPackage(t *testing.T) {
	str1 := "test-string"
	str2 := "teststring"

	strout := ToPackage(str1)
	assert.Equal(t, strout, str2, "Expected string to be formatted as a valid go package name")

}

func TestToPascal(t *testing.T) {
	str1 := "test-string"
	str2 := "TestString"

	strout := ToPascal(str1)
	assert.Equal(t, strout, str2, "Expected string to be formatted in PascalCase")

}

func TestSnakeCase(t *testing.T) {
	str1 := "test-string"
	str2 := "test_string"

	strout := ToSnakeCase(str1)
	assert.Equal(t, strout, str2, "Expected string to be formatted in snake_case")

}
