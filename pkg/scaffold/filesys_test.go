package scaffold

import (
	"testing"
)

type PathTestData struct {
	TestPath  string
	ExpectRes string
}

func TestGetFinalElementFromPath(t *testing.T) {

	testTable := []PathTestData{
		// single directory
		{
			"internal",
			"internal",
		},
		// single nesting
		{
			"internal/app",
			"app",
		},
		// deeper nesting
		{
			"internal/db/migrations",
			"migrations",
		},
	}

	for _, v := range testTable {
		res := getFinalElementFromPath(v.TestPath)
		if res != v.ExpectRes {
			t.Errorf("GetFinalElementFromPath(%s) = %s; want %s", v.TestPath, res, v.ExpectRes)
		}
	}
}

func TestGetParentDirFromPath(t *testing.T) {

	testTable := []PathTestData{
		// single directory
		{
			"internal",
			"",
		},
		// single nesting
		{
			"internal/app",
			"internal",
		},
		// deeper nesting
		{
			"internal/db/migrations",
			"db",
		},
	}

	for _, v := range testTable {
		res := getParentDirFromPath(v.TestPath)
		if res != v.ExpectRes {
			t.Errorf("GetParentDirFromPath(%s) = %s; want %s", v.TestPath, res, v.ExpectRes)
		}
	}
}
