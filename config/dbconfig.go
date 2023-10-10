package config

type DbConfig struct {
	DriverName     string `default:"sqlite3"`
	DataSourceName string `default:"./sqlite.Db"`
	MaxIdleConn    int    `default:"20"`
	MaxOpenConn    int    `default:"10"`
}
