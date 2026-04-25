package models

import (
	"gorm.io/plugin/soft_delete"
	"orm-crud/gormc/mixin"
)

func init() {
	Models = append(Models, &SysLanguage{})
}

type SysLanguage struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.OperatorID
	mixin.SortOrder
	mixin.IsEnabled

	DeletedAt    soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;default:0;uniqueIndex:idx_sys_language_language_code_deleted_at,priority:2" json:"deletedAt"`
	LanguageCode string                `gorm:"column:language_code;type:varchar(128);not null;uniqueIndex:idx_sys_language_language_code_deleted_at,priority:1;comment:标准语言代码" json:"languageCode"`
	LanguageName string                `gorm:"column:language_name;type:varchar(255);not null;comment:语言名称" json:"languageName"`
	NativeName   string                `gorm:"column:native_name;type:varchar(255);not null;default:'';comment:本地语言名称" json:"nativeName"`
	IsDefault    bool                  `gorm:"column:is_default;default:0;comment:是否为默认语言" json:"isDefault"`
}

// TableName 指定表名
func (SysLanguage) TableName() string {
	return "sys_languages"
}
