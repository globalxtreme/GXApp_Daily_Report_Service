package database

import (
	xtremedb "github.com/globalxtreme/go-core/v2/database"
	"service/internal/app/database/migration"
)

func Migrations() []xtremedb.Migration {
	return []xtremedb.Migration{
		&migration.Activity_1726651211960757{},
		// Daily Report tables
		&migration.ReportUsers_20260429000001{},
		&migration.Reports_20260429000002{},
		&migration.ReportSessions_20260429000003{},
	}
}
