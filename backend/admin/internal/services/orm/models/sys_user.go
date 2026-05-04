package models

import (
	"gorm.io/plugin/soft_delete"
	"orm-crud/gormc/mixin"
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
	mixin.IsEnabled
	DeletedAt    soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;uniqueIndex:idx_sys_user_username_deleted_at,priority:2" json:"deletedAt"`
	Username     string                `gorm:"column:username;type:varchar(64);not null;uniqueIndex:idx_sys_user_username_deleted_at,priority:1;comment:用户名" json:"username"`
	Nickname     string                `gorm:"column:nickname;type:varchar(64);not null;default:'';comment:昵称" json:"nickname"`
	Password     string                `gorm:"column:password;type:varchar(255);not null;default:'';comment:密码" json:"password"`
	LastLoginAt  *time.Time            `gorm:"column:last_login_at;comment:最后一次登录的时间" json:"lastLoginAt"`
	LastLoginIP  string                `gorm:"column:last_login_ip;type:varchar(45);not null;default:'';comment:最后一次登录的IP" json:"lastLoginIP"`
	SysRoles     []SysRole             `gorm:"many2many:sys_user_role;foreignKey:ID;joinForeignKey:UserID;References:ID;joinReferences:RoleID" json:"sysRoles"`
	LanguageCode string                `gorm:"column:language_code;type:varchar(32);not null:default:'';comment:语言代码" json:"languageCode"`
}

// TableName SysUser's table name
func (*SysUser) TableName() string {
	return "sys_user"
}
