package models

import (
	"gorm.io/plugin/soft_delete"
	"orm-crud/gormc/mixin"
)

func init() {
	Models = append(Models, &SysLanguageType{})
}

type SysLanguageType struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.CreatedBy
	mixin.UpdatedBy
	mixin.SortOrder
	mixin.IsEnabled
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;uniqueIndex:idx_sys_language_type_code_deleted_at,priority:2" json:"deletedAt"`
	TypeCode  string                `gorm:"column:type_code;type:varchar(128);not null;uniqueIndex:idx_sys_language_type_code_deleted_at,priority:1;comment:标准语言代码" json:"typeCode"`
	TypeName  string                `gorm:"column:type_name;type:varchar(255);not null;comment:语言名称" json:"typeName"`
	IsDefault bool                  `gorm:"column:is_default;default:0;comment:是否为默认语言" json:"isDefault"`
	Entries   []SysLanguageEntry    `gorm:"foreignKey:SysLanguageTypeId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"entries"`
}

// TableName 指定表名
func (SysLanguageType) TableName() string {
	return "sys_language_type"
}
