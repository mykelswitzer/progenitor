package scaffold

import (
	"fmt"
	"reflect"
	"testing"
)

type CollectDirsTestData struct {
	TestPath  string
	ExpectRes map[string][]string
}

func TestCollectDirs(t *testing.T) {

	testTable := []CollectDirsTestData{
		{
			TestPath: "internal",
			ExpectRes: map[string][]string{
				"": []string{"internal"},
			},
		},
		{
			TestPath: "internal/app",
			ExpectRes: map[string][]string{
				"internal": []string{"app"},
			},
		},
		{
			TestPath: "internal/db/migrations",
			ExpectRes: map[string][]string{
				"db": []string{"migrations"},
			},
		},
	}

	for _, test := range testTable {
		res := collectDirs(map[string][]string{}, test.TestPath)
		if !reflect.DeepEqual(res, test.ExpectRes) {
			for k, v := range res {
				fmt.Println("Result", k, "value is", v)
			}
			for k, v := range test.ExpectRes {
				fmt.Println("Expected Result", k, "value is", v)
			}
			t.Errorf("ParseDir failed")
		}
	}
}

type PathTestData struct {
	TestPath  string
	ExpectDir string
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
		dir, parent := getDirAndParentFromPath(v.TestPath)
		if dir != v.ExpectDir || parent != v.ExpectParent {
			t.Errorf("GetDirAndParentFromPath(%s) = %s,%s ; want %s,%s", v.TestPath, dir, parent, v.ExpectDir, v.ExpectParent)
		}
	}
}
