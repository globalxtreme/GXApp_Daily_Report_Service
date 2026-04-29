package migration

import (
	"os"

	xtremedb "github.com/globalxtreme/go-core/v2/database"
	"service/internal/pkg/config"
	"service/internal/pkg/model"
)

type ReportUsers_20260429000001 struct{}

func (ReportUsers_20260429000001) Reference() string {
	return "ReportUsers_20260429000001"
}

func (ReportUsers_20260429000001) Tables() []xtremedb.Table {
	owner := os.Getenv("DB_OWNER")
	return []xtremedb.Table{
		{Connection: config.PgSQL, CreateTable: model.ReportUser{}, Owner: owner},
	}
}

func (ReportUsers_20260429000001) Columns() []xtremedb.Column {
	return []xtremedb.Column{}
}
