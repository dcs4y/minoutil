package customer

import (
	"game/utils/dbclient"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB // 数据库连接
)

// 管理端业务初始化
func init() {
	// 初始化数据库
	DB = dbclient.DB
}
