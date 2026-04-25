package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
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
	mixin.OperatorID
	mixin.IsEnabled
	DeletedAt     *gorm.DeletedAt             `gorm:"column:deleted_at;softDelete:milli;uniqueIndex:idx_sys_role_code_deleted_at,priority:2" json:"deletedAt,omitempty"`
	Name          string                      `gorm:"column:name;type:varchar(255);not null;comment:角色名称" json:"name,omitempty"`
	Code          string                      `gorm:"column:code;type:varchar(128);not null;uniqueIndex:idx_sys_role_code_deleted_at,priority:1;comment:角色标识" json:"code,omitempty"`
	Menus         datatypes.JSONSlice[uint64] `gorm:"column:menus;type:json;comment:分配的菜单列表" json:"menus,omitempty"`
	Apis          datatypes.JSONSlice[string] `gorm:"column:apis;type:json;comment:分配的API列表" json:"apis,omitempty"`
	DataScope     string                      `gorm:"column:data_scope;type:varchar(32);not null;default:'';comment:数据权限范围" json:"dataScope,omitempty"`
	ParentID      uint64                      `gorm:"column:parent_id;type:bigint;not null;default:0;comment:父级ID" json:"parentID,omitempty"`
	Path          string                      `gorm:"column:path;type:varchar(1024);not null;default:'';comment:节点路径" json:"path,omitempty"`
	ChildIDs      datatypes.JSON              `gorm:"column:child_ids;type:json;comment:所有子节点ID" json:"childIDs,omitempty"`
	ParentSysRole *SysRole                    `gorm:"foreignKey:ParentID;references:ID" json:"parentSysRole,omitempty"`
	Children      []SysRole                   `gorm:"foreignKey:ParentID;references:ID" json:"children,omitempty"`
}

// TableName SysRole's table name
func (*SysRole) TableName() string {
	return "sys_role"
}
