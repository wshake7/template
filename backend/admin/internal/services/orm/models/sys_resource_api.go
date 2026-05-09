package models

import (
	"gorm.io/plugin/soft_delete"
	"orm-crud/gormc/mixin"
)

func init() {
	//Models = append(Models, &SysResourceApi{})
}

// SysResourceApi 对应表 sys_resource_api
type SysResourceApi struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.SortOrder
	mixin.OperatorID
	mixin.Remark
	mixin.IsEnabled
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;index" json:"deletedAt"`
	Module    string                `gorm:"column:module;type:varchar(128);not null;default:'';comment:所属业务模块" json:"module"`
	Path      string                `gorm:"column:path;type:varchar(255);not null;comment:接口路径" json:"path"`
	Method    string                `gorm:"column:method;type:varchar(16);not null;comment:请求方法" json:"method"`
}

func (SysResourceApi) TableName() string {
	return "sys_resource_api"
}
