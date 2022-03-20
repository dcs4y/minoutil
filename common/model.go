package common

import (
	"gorm.io/gorm"
	"time"
)

// BaseModel 数据库表的基本字段。根据gorm模型定义。
type BaseModel struct {
	ID        uint64    `gorm:"primarykey"` // 主键ID
	CreatedAt time.Time // 创建时间
}

type BaseUpdateModel struct {
	ID        uint64    `gorm:"primarykey"` // 主键ID
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
}

type BaseDeleteModel struct {
	ID        uint64         `gorm:"primarykey"` // 主键ID
	CreatedAt time.Time      // 创建时间
	UpdatedAt time.Time      // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}
