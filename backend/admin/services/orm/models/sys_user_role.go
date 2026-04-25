package models

import (
	"gorm.io/gorm"
	"orm-crud/gormc/mixin"
)

func init() {
	Models = append(Models, &SysUserRole{})
}

type SysUserRole struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.OperatorID
	mixin.IsEnabled
	UserID    uint64          `gorm:"column:user_id;type:bigint;comment:用户ID;index:idx_sys_user_role_user_id;uniqueIndex:idx_sys_user_role_user_id_role_id_delete_at,priority:1;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"userID,omitempty"`
	RoleID    uint64          `gorm:"column:role_id;type:bigint;comment:角色ID;index:idx_sys_user_role_role_id;uniqueIndex:idx_sys_user_role_user_id_role_id_delete_at,priority:2;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"roleID,omitempty"`
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at;softDelete:milli;uniqueIndex:idx_sys_user_role_user_id_role_id_delete_at,priority:3" json:"deletedAt,omitempty"`
	SysUser   *SysUser        `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"sysUser,omitempty"`
	SysRole   *SysRole        `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"sysRole,omitempty"`
}

// TableName 指定表名
func (SysUserRole) TableName() string {
	return "sys_user_role"
}
