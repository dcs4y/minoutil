package model

import (
	"game/common"
)

// Manager 系统管理员表
type Manager struct {
	common.BaseDeleteModel
	Name     string // 姓名
	UserName string // 用户名
	Password string // 密码
	State    int8   // 状态。0.停用；1.启用；2.删除；
}

func (t Manager) TableName() string {
	return "sys_manager"
}

// ManagerRole 管理员角色
type ManagerRole struct {
	common.BaseModel
	ManagerId uint64 // 管理员ID
	RoleId    uint64 // 角色ID
}

func (t ManagerRole) TableName() string {
	return "sys_manager_role"
}

// Menu 菜单
type Menu struct {
	common.BaseModel
	ParentId  uint64 // 上级菜单ID
	Path      string // 路由path
	Name      string // 路由name
	Hidden    bool   // 是否在列表隐藏
	Component string // 对应前端文件路径
	Type      string // 类型。menu：菜单，button：按钮，function：功能(不显示)。
	Sort      uint   // 排序
	Remark    string // 备注
	Children  []Menu `json:"children"`
}

func (t Menu) TableName() string {
	return "sys_menu"
}

// MenuResource 菜单资源
type MenuResource struct {
	common.BaseModel
	MenuId     uint64 // 菜单ID
	ResourceId uint64 // 资源ID
}

func (t MenuResource) TableName() string {
	return "sys_menu_resource"
}

// Resource 资源
type Resource struct {
	common.BaseModel
	Name   string // 名称
	Uri    string // 路径。http_uri。
	Weight int    // 权重，值越大优先级超高。
	Remark string // 备注
}

func (t Resource) TableName() string {
	return "sys_resource"
}

// Role 角色
type Role struct {
	common.BaseModel
	Name   string // 名称
	Code   string // 角色编码
	Remark string // 备注
}

func (t Role) TableName() string {
	return "sys_role"
}

// RoleMenu 角色菜单
type RoleMenu struct {
	common.BaseModel
	RoleId uint64 // 角色ID
	MenuId uint64 // 菜单ID
}

func (t RoleMenu) TableName() string {
	return "sys_role_menu"
}
