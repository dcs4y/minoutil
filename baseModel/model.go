package baseModel

import (
	"gorm.io/gorm"
)

// BaseModel 数据库表的基本字段。根据gorm模型定义。
type BaseModel struct {
	Id        int64 `gorm:"primarykey" json:"id"` // 主键ID
	CreatedAt Time  `json:"createTime"`           // 创建时间
}

type BaseUpdateModel struct {
	Id        int64 `gorm:"primarykey" json:"id"` // 主键ID
	CreatedAt Time  `json:"createTime"`           // 创建时间
	UpdatedAt Time  `json:"updateTime"`           // 更新时间
}

type BaseDeleteModel struct {
	Id        int64          `gorm:"primarykey" json:"id"` // 主键ID
	CreatedAt Time           `json:"createTime"`           // 创建时间
	UpdatedAt Time           `json:"updateTime"`           // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`       // 删除时间
}
