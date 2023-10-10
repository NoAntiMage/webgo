package common

import (
	"goweb/common/db"
	"goweb/common/logx"
	"goweb/common/singleflight"
	"goweb/common/validx"
	"goweb/config"
)

func GlobalInit() {
	validx.ValidatorInit()
	config.ConfigInit()
	logx.LoggerInit()
	db.DbInit()
	singleflight.SingleFlightInit()
}
