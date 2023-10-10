package db

import (
	"goweb/common/logx"
	"goweb/common/stringx"
	"goweb/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func GetDb() *sqlx.DB {
	return db
}

func DbInit() {
	var err error
	dbConf := config.GetGlobalConfig().DbConf

	db, err = sqlx.Open(dbConf.DriverName, dbConf.DataSourceName)
	// db, err = sqlx.Open("sqlite3", "/Users/wujimaster/Desktop/web-go/sqlite.Db")
	if err != nil {
		logx.Loggerx.Infof("%+v\n", err)
	}
	db.SetMaxIdleConns(dbConf.MaxIdleConn)
	db.SetMaxOpenConns(dbConf.MaxOpenConn)
	db.MapperFunc(stringx.SnakeFormat)
}

func DbClose() {
	logx.Loggerx.Info("Db disconnecting...")
	db.Close()
}

//get database version.A second way to check database connection alive.
func GetDbVersion() (sqliteVersion string) {
	db.Get(&sqliteVersion, "select SQLITE_VERSION() as version")
	return sqliteVersion
}
