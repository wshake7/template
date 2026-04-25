package models

import (
	"gorm.io/gorm"
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

	DeletedAt    *gorm.DeletedAt `gorm:"column:deleted_at;softDelete:milli;uniqueIndex:idx_sys_language_language_code_deleted_at,priority:2" json:"deletedAt,omitempty"`
	LanguageCode string          `gorm:"column:language_code;type:varchar(128);not null;uniqueIndex:idx_sys_language_language_code_deleted_at,priority:1;comment:标准语言代码" json:"languageCode,omitempty"`
	LanguageName string          `gorm:"column:language_name;type:varchar(255);not null;comment:语言名称" json:"languageName,omitempty"`
	NativeName   string          `gorm:"column:native_name;type:varchar(255);not null;default:'';comment:本地语言名称" json:"nativeName,omitempty"`
	IsDefault    bool            `gorm:"column:is_default;default:0;comment:是否为默认语言" json:"isDefault,omitempty"`
}

// TableName 指定表名
func (SysLanguage) TableName() string {
	return "sys_languages"
}
