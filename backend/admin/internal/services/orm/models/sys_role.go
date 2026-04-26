package models

import (
	"gorm.io/datatypes"
	"gorm.io/plugin/soft_delete"
	"orm-crud/gormc/mixin"
)

func init() {
	Models = append(Models, &SysRole{})
}

type SysRole struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.Remark
	mixin.CreatedBy
	mixin.UpdatedBy
	mixin.IsEnabled
	DeletedAt     soft_delete.DeletedAt       `gorm:"column:deleted_at;softDelete:milli;default:0;uniqueIndex:idx_sys_role_code_deleted_at,priority:2" json:"deletedAt"`
	Name          string                      `gorm:"column:name;type:varchar(255);not null;comment:角色名称" json:"name"`
	Code          string                      `gorm:"column:code;type:varchar(128);not null;uniqueIndex:idx_sys_role_code_deleted_at,priority:1;comment:角色标识" json:"code"`
	Menus         datatypes.JSONSlice[uint64] `gorm:"column:menus;type:json;comment:分配的菜单列表" json:"menus"`
	Apis          datatypes.JSONSlice[string] `gorm:"column:apis;type:json;comment:分配的API列表" json:"apis"`
	DataScope     string                      `gorm:"column:data_scope;type:varchar(32);not null;default:'';comment:数据权限范围" json:"dataScope"`
	ParentID      uint64                      `gorm:"column:parent_id;type:bigint;not null;default:0;comment:父级ID" json:"parentID"`
	Path          string                      `gorm:"column:path;type:varchar(1024);not null;default:'';comment:节点路径" json:"path"`
	ChildIDs      datatypes.JSON              `gorm:"column:child_ids;type:json;comment:所有子节点ID" json:"childIDs"`
	ParentSysRole *SysRole                    `gorm:"foreignKey:ParentID;references:ID" json:"parentSysRole"`
	Children      []SysRole                   `gorm:"foreignKey:ParentID;references:ID" json:"children"`
}

// TableName SysRole's table name
func (*SysRole) TableName() string {
	return "sys_role"
}
