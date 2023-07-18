package scaffold

import (
	"fmt"
	"reflect"
	"testing"
)

type MapDirsTestData struct {
	TestPath  string
	ExpectRes map[string][]string
}

func TestMapDirs(t *testing.T) {

	testTable := []MapDirsTestData{
		{
			TestPath: "internal",
			ExpectRes: map[string][]string{
				"":         []string{"internal"},
				"internal": []string{},
			},
		},
		{
			TestPath: "internal/app",
			ExpectRes: map[string][]string{
				"internal": []string{"app"},
				"app":      []string{},
			},
		},
		{
			TestPath: "internal/db/migrations",
			ExpectRes: map[string][]string{
				"db":         []string{"migrations"},
				"migrations": []string{},
			},
		},
	}

	for _, test := range testTable {
		res := mapDirs(map[string][]string{}, test.TestPath)
		if !reflect.DeepEqual(res, test.ExpectRes) {
			for k, v := range res {
				fmt.Println("Result", k, "value is", v)
			}
			for k, v := range test.ExpectRes {
				fmt.Println("Expected Result", k, "value is", v)
			}
			t.Errorf("mapDirs failed")
		}
	}
}

func TestPopulateStructureFromMap(t *testing.T) {

	dbDir := Dir{Name: "db"}
	dbDir.AddSubDirs(Dir{Name: "migrations"})

	internalDir := Dir{Name: "internal"}
	internalDir.AddSubDirs(dbDir, Dir{Name: "app"})

	baseDir := Dir{Name: ""}
	baseDir.AddSubDirs(internalDir)

	testData := map[string][]string{
		"":           []string{"internal"},
		"internal":   []string{"db", "app"},
		"app":        []string{},
		"db":         []string{"migrations"},
		"migrations": []string{},
	}

	res, err := populateStructureFromMap(testData, "")
	if err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(res, baseDir) {
		t.Errorf("populateStructureFromMap failed")
	}

}

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
		dir, parent := getDirAndParentFromPath(v.TestPath)
		if dir != v.ExpectDir || parent != v.ExpectParent {
			t.Errorf("getDirAndParentFromPath(%s) = %s,%s ; want %s,%s", v.TestPath, dir, parent, v.ExpectDir, v.ExpectParent)
		}
	}
}
