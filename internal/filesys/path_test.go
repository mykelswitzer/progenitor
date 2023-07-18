package filesys

import (
	"fmt"
	"reflect"
	"testing"
)

type PathTestData struct {
	TestPath     string
	ExpectDir    string
	ExpectParent string
}

func TestGetDirAndParentFromPath(t *testing.T) {

	testTable := []PathTestData{
		// single directory
		{
			"internal",
			"internal",
			"",
		},
		// single nesting
		{
			"internal/app",
			"app",
			"internal",
		},
		// deeper nesting
		{
			"internal/db/migrations",
			"migrations",
			"db",
		},
	}

	for _, v := range testTable {
		dir, parent := GetDirAndParentFromPath(v.TestPath)
		if dir != v.ExpectDir || parent != v.ExpectParent {
			t.Errorf("GetDirAndParentFromPath(%s) = %s,%s ; want %s,%s", v.TestPath, dir, parent, v.ExpectDir, v.ExpectParent)
		}
	}
}
