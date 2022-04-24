package database

import (
	"fmt"
	"gitee.com/snxamdf/http-server/src/config"
	"gitee.com/snxamdf/http-server/src/consts"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

func init() {
	defer func() {
		if err := recover(); err != nil {
			consts.GlobalPanicRecoverString = "初始Sqlite数据库错误 " + (err.(error)).Error()
		}
	}()
	dbConf := config.Cfg.Sqlite3
	if dbConf.Path == "" || dbConf.Name == "" {
		consts.GlobalPanicRecoverString = fmt.Sprintf("Sqlite数据库配置错误 Path:%s Name:%s ", dbConf.Path, dbConf.Name)
		return
	}
	var err error
	var dbName = dbConf.Path + dbConf.Name + ".db"
	db, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		consts.GlobalPanicRecoverString = "打开Sqlite数据库失败:" + err.Error()
	}
}
