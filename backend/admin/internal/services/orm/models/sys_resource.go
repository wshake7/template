package models

import (
	"orm-crud/gormc/mixin"

	"gorm.io/plugin/soft_delete"
)

func init() {
	Models = append(Models, &SysResource{})
}

type SysResource struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.OperatorID
	mixin.IsEnabled
	mixin.Remark
	mixin.Metadata
	DeletedAt    soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;index:idx_sys_resource_deleted_at" json:"deletedAt"`
	ResourceType string                `gorm:"column:resource_type;type:varchar(32);not null;comment:资源类型: api / data / menu" json:"resource_type"`
	Code         string                `gorm:"column:code;type:varchar(255);not null;uniqueIndex:idx_sys_resource_code_active,where:deleted_at = 0;comment:资源唯一标识" json:"code"`
	Name         string                `gorm:"column:name;type:varchar(255);not null;comment:资源名称" json:"name"`
}

func (SysResource) TableName() string {
	return "sys_resource"
}
