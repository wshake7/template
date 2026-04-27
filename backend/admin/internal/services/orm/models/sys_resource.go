package models

import (
	"orm-crud/gormc/mixin"

	"gorm.io/plugin/soft_delete"
)

func init() {
	Models = append(Models, &SysResource{})
}

// SysResource 对应表 sys_resource
type SysResource struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.CreatedBy
	mixin.UpdatedBy
	mixin.IsEnabled
	mixin.Remark
	mixin.Metadata
	DeletedAt    soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;default:0;uniqueIndex:idx_sys_resource_code_deleted_at,priority:2" json:"deletedAt"`
	ResourceType string                `gorm:"column:resource_type;type:varchar(32);not null;comment:资源类型: api / data / menu / component" json:"resourceType"`
	ResourceCode string                `gorm:"column:resource_code;type:varchar(255);not null;uniqueIndex:idx_sys_resource_code_deleted_at,priority:1;comment:资源唯一标识" json:"resourceCode"`
}

func (SysResource) TableName() string {
	return "sys_resource"
}
