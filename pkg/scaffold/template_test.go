package scaffold

import (
	"testing"
)

type GetParentDirTestData struct {
	Test   string
	Result string
}

func TestGetParentDirFromPath(t *testing.T) {

	testTable := []GetParentDirTestData{
		{
			"apple",
			"apple",
		},
		{
			"terraform/rds.tf.tmpl",
			"terraform",
		},
		{
			"internal/db/migrations",
			"db",
		},
		{
			"internal/db/migrations/010000_create_initial_tables.down.sql",
			"migrations",
		},
	}

	for _, v := range testTable {
		res := getParentDirFromPath(v.Test)
		if res != v.Result {
			t.Errorf("GetParentDir(%s) = %s; want %s", v.Test, res, v.Result)
		}
	}
}
