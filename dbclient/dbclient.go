package dbclient

import (
	"fmt"
	"github.com/dcs4y/minoutil/logutil"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

var (
	clients = make(map[string]*gorm.DB) // 多数据源
	lock    sync.RWMutex
)

type DBConfig struct {
	Url         string `yaml:"url"`
	BbType      string `yaml:"dbType"`
	ShowSql     bool   `yaml:"showSql"`
	MaxOpenConn int    `yaml:"maxOpenConn"`
	MaxIdleConn int    `yaml:"maxIdleConn"`
	ConsoleLog  bool   //是否打印控制台日志
}

func NewClient(name string, config DBConfig) (DB *gorm.DB) {
	lock.RLock()
	defer lock.RUnlock()
	db, b := clients[name]
	if b {
		return db
	}
	var dialect gorm.Dialector
	if config.BbType == "mysql" {
		dialect = mysql.Open(config.Url)
	} else {
		return
	}
	// 是否显示日志
	var newLogger logger.Interface
	if config.ShowSql {
		var w io.Writer
		if config.ConsoleLog {
			w = io.MultiWriter(os.Stdout, logutil.GetLog("db").Writer())
		} else {
			w = logutil.GetLog("db").Writer()
		}
		newLogger = logger.New(
			log.New(w, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
			logger.Config{
				SlowThreshold:             time.Millisecond, // 慢 SQL 阈值
				LogLevel:                  logger.Info,      // 日志级别
				IgnoreRecordNotFoundError: true,             // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  true,             // 禁用彩色打印
			},
		)
	}
	db, err := gorm.Open(dialect, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		fmt.Println("数据库连接失败：", err)
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("数据库SQL连接失败：", err)
		return
	} else {
		sqlDB.SetMaxOpenConns(config.MaxOpenConn)
		sqlDB.SetMaxIdleConns(config.MaxIdleConn)
		sqlDB.SetConnMaxLifetime(time.Minute)
	}
	clients[name] = db
	return db
}

func GetClientByName(name string) *gorm.DB {
	return clients[name]
}

func GetClient() *gorm.DB {
	return clients[""]
}
