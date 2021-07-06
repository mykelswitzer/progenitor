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
	assert.Equal(t, strout, str2, "Expected hyphenated string to be formatted in PascalCase")

	str1 = "teststring"
	str2 = "Teststring"

	strout = ToPascal(str1)
	assert.Equal(t, strout, str2, "Expected non-hyphenated string to be formatted in PascalCase")

}

func TestToPlural(t *testing.T) {
	str1 := "box"
	str2 := "boxes"

	strout := ToPlural(str1)
	assert.Equal(t, strout, str2, "Expected string to be Pluralized")

	str1 = "partner-community"
	str2 = "partner-communities"

	strout = ToPlural(str1)
	assert.Equal(t, strout, str2, "Expected string to be Pluralized")

	str1 = "addresses"
	str2 = "addresses"

	strout = ToPlural(str1)
	assert.Equal(t, strout, str2, "Expected string to be Pluralized")

	str1 = "weekday"
	str2 = "weekdays"

	strout = ToPlural(str1)
	assert.Equal(t, strout, str2, "Expected string to be Pluralized")

}

func TestSnakeCase(t *testing.T) {
	str1 := "test-string"
	str2 := "test_string"

	strout := ToSnakeCase(str1)
	assert.Equal(t, strout, str2, "Expected string to be formatted in snake_case")

}
