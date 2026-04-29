package migration

import (
	"os"

	xtremedb "github.com/globalxtreme/go-core/v2/database"
	"service/internal/pkg/config"
	"service/internal/pkg/model"
)

type Reports_20260429000002 struct{}

func (Reports_20260429000002) Reference() string {
	return "Reports_20260429000002"
}

func (Reports_20260429000002) Tables() []xtremedb.Table {
	owner := os.Getenv("DB_OWNER")
	return []xtremedb.Table{
		{Connection: config.PgSQL, CreateTable: model.Report{}, Owner: owner},
	}
}

func (Reports_20260429000002) Columns() []xtremedb.Column {
	return []xtremedb.Column{}
}
