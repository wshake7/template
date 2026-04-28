package models

import (
	"gorm.io/plugin/soft_delete"
	"orm-crud/gormc/mixin"
)

func init() {
	Models = append(Models, &SysCasbinModel{})
}

type SysCasbinModel struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.OperatorID
	mixin.IsEnabled
	mixin.Remark
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;index" json:"deletedAt"`
	Name      string                `gorm:"column:name;type:varchar(255);not null;uniqueIndex;comment:模型名称" json:"name"`
	Content   string                `gorm:"column:content;type:text;not null;comment:模型内容" json:"content"`
}

func (SysCasbinModel) TableName() string {
	return "sys_casbin_model"
}
