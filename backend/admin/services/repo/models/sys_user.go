package models

import (
	"gorm.io/plugin/soft_delete"
	"orm-crud/gorm/mixin"
	"time"
)

func init() {
	Models = append(Models, &SysUser{})
}

type SysUser struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.Remark
	mixin.OperatorID
	mixin.Status
	DeletedAt   soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;uniqueIndex:idx_sys_user_username_deleted_at,priority:2"`
	Username    string                `gorm:"column:username;type:varchar(64);not null;uniqueIndex:idx_sys_user_username_deleted_at,priority:1;comment:用户名"`
	Nickname    string                `gorm:"column:nickname;type:varchar(64);not null;default:'';comment:昵称"`
	Password    string                `gorm:"column:password;type:varchar(255);not null;default:'';comment:密码"`
	LastLoginAt *time.Time            `gorm:"column:last_login_at;comment:最后一次登录的时间"`
	LastLoginIP string                `gorm:"column:last_login_ip;type:varchar(45);not null;default:'';comment:最后一次登录的IP"`
	SysRoles    []SysRole             `gorm:"many2many:sys_user_role;foreignKey:ID;joinForeignKey:UserID;References:ID;joinReferences:RoleID"`
}

// TableName SysUser's table name
func (*SysUser) TableName() string {
	return "sys_user"
}
