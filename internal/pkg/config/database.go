package config

import (
	xtremedb "github.com/globalxtreme/go-core/v2/database"
	"gorm.io/gorm"
	"os"
	"time"
)

var (
	PgSQL *gorm.DB
)

const CONNECTION_DB_PGSQL = "PgSQL"

type connectionConf struct {
	SQL             *gorm.DB
	DBConf          xtremedb.DBConf
	MaxOpenCons     int
	MaxIdleCons     int
	MaxLifetimeCons time.Duration
}

func InitDB(connections ...string) func() {
	configurations := map[string]connectionConf{
		CONNECTION_DB_PGSQL: {
			DBConf: xtremedb.DBConf{
				Driver:    xtremedb.POSTGRESQL_DRIVER,
				Host:      os.Getenv("DB_HOST"),
				Port:      os.Getenv("DB_PORT"),
				Username:  os.Getenv("DB_USERNAME"),
				Password:  os.Getenv("DB_PASSWORD"),
				Database:  os.Getenv("DB_DATABASE"),
				ParseTime: true,
				TimeZone:  "Asia/Makassar",
			},
		},
	}

	configurationSelected := map[string]connectionConf{}
	if len(connections) > 0 {
		for _, conn := range connections {
			configurationSelected[conn] = configurations[conn]
		}
	} else {
		configurationSelected = configurations
	}

	var dbCloses []func()
	for connName, conf := range configurationSelected {
		dbSQL, DBClose := xtremedb.Connect(conf.DBConf)

		dbCloses = append(dbCloses, DBClose)

		switch connName {
		case CONNECTION_DB_PGSQL:
			PgSQL = dbSQL
		}
	}

	return func() {
		if len(dbCloses) > 0 {
			for _, DBClose := range dbCloses {
				DBClose()
			}
		}
	}
}
