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
	mixin.OperatorID
	mixin.IsEnabled
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;uniqueIndex:idx_sys_role_code_deleted_at,priority:2" json:"deletedAt"`
	Name      string                `gorm:"column:name;type:varchar(255);not null;comment:角色名称" json:"name"`
	Code      string                `gorm:"column:code;type:varchar(128);not null;uniqueIndex:idx_sys_role_code_deleted_at,priority:1;comment:角色标识" json:"code"`
	//Menus         datatypes.JSONSlice[uint64] `gorm:"column:menus;type:json;comment:分配的菜单列表" json:"menus"`
	//Apis          datatypes.JSONSlice[string] `gorm:"column:apis;type:json;comment:分配的API列表" json:"apis"`
	ParentID      *uint64        `gorm:"column:parent_id;type:bigint;comment:父级ID" json:"parentID"`
	ChildIDs      datatypes.JSON `gorm:"column:child_ids;default:'[]';comment:所有子节点ID" json:"childIDs"`
	ParentSysRole *SysRole       `gorm:"foreignKey:ParentID;references:ID" json:"parentSysRole"`
	Children      []SysRole      `gorm:"foreignKey:ParentID;references:ID" json:"children"`
}

// TableName SysRole's table name
func (*SysRole) TableName() string {
	return "sys_role"
}
