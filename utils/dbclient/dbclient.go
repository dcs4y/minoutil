package dbclient

import (
	"fmt"
	"game/common"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"sync"
)

var (
	DB    *gorm.DB            // 默认数据库连接
	dbMap map[string]*gorm.DB // 多数据源
	lock  sync.RWMutex
)

func init() {
	// 初始化数据库
	dbMap = make(map[string]*gorm.DB)
	DB = NewClient(common.WS.DBConfig)
}

func NewClient(dbConfig *common.DatabaseModel) (DB *gorm.DB) {
	lock.RLock()
	defer lock.RUnlock()
	db, b := dbMap[dbConfig.Url]
	if b {
		return db
	}
	var dialect gorm.Dialector
	if dbConfig.BbType == "mysql" {
		dialect = mysql.Open(dbConfig.Url)
	} else {
		return nil
	}
	db, err := gorm.Open(dialect, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		fmt.Println("数据库连接失败：", err)
		return
	}
	if dbConfig.ShowSql {
		db.Logger.LogMode(logger.Info)
	}
	dbMap[dbConfig.Url] = db
	return db
}
