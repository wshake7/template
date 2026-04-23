package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"orm-crud/gorm/mixin"
)

func init() {
	Models = append(Models, &SysRole{})
}

type SysRole struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.Remark
	mixin.OperatorID
	mixin.Status
	DeletedAt     *gorm.DeletedAt             `gorm:"column:deleted_at;softDelete:milli;uniqueIndex:idx_sys_role_code_deleted_at,priority:2"`
	Name          string                      `gorm:"column:name;type:varchar(255);not null;comment:角色名称"`
	Code          string                      `gorm:"column:code;type:varchar(128);not null;uniqueIndex:idx_sys_role_code_deleted_at,priority:1;comment:角色标识"`
	Menus         datatypes.JSONSlice[uint64] `gorm:"column:menus;type:json;comment:分配的菜单列表"`
	Apis          datatypes.JSONSlice[string] `gorm:"column:apis;type:json;comment:分配的API列表"`
	DataScope     string                      `gorm:"column:data_scope;type:varchar(32);not null;default:'';comment:数据权限范围"`
	ParentID      uint64                      `gorm:"column:parent_id;type:bigint;not null;default:0;comment:父级ID"`
	Path          string                      `gorm:"column:path;type:varchar(1024);not null;default:'';comment:节点路径"`
	ChildIDs      datatypes.JSON              `gorm:"column:child_ids;type:json;comment:所有子节点ID"`
	ParentSysRole *SysRole                    `gorm:"foreignKey:ParentID;references:ID"`
	Children      []SysRole                   `gorm:"foreignKey:ParentID;references:ID"`
}

// TableName SysRole's table name
func (*SysRole) TableName() string {
	return "sys_role"
}
