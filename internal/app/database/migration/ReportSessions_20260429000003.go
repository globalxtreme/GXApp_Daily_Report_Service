package migration

import (
	"os"

	xtremedb "github.com/globalxtreme/go-core/v2/database"
	"service/internal/pkg/config"
	"service/internal/pkg/model"
)

type ReportSessions_20260429000003 struct{}

func (ReportSessions_20260429000003) Reference() string {
	return "ReportSessions_20260429000003"
}

func (ReportSessions_20260429000003) Tables() []xtremedb.Table {
	owner := os.Getenv("DB_OWNER")
	return []xtremedb.Table{
		{Connection: config.PgSQL, CreateTable: model.ReportSession{}, Owner: owner},
	}
}

func (ReportSessions_20260429000003) Columns() []xtremedb.Column {
	return []xtremedb.Column{}
}
