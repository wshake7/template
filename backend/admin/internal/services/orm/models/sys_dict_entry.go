package models

import (
	"gorm.io/plugin/soft_delete"
	"orm-crud/gormc/mixin"
)

func init() {
	Models = append(Models, &SysDictEntry{})
}

// SysDictEntry 对应表 sys_dict_entries
type SysDictEntry struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.OperatorID
	mixin.SortOrder
	mixin.IsEnabled
	mixin.Remark
	DeletedAt      soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;index" json:"deletedAt"`
	LabelComponent string                `gorm:"column:label_component;type:varchar(255);not null;default:'';comment:字典项的显示标签组件" json:"labelComponent"`
	EntryLabel     string                `gorm:"column:entry_label;type:varchar(255);not null;comment:字典项的显示标签" json:"entryLabel"`
	EntryValue     string                `gorm:"column:entry_value;type:varchar(255);not null;comment:字典项的实际值" json:"entryValue"`
	LanguageCode   string                `gorm:"column:language_code;type:varchar(32);not null:default:'';comment:语言代码" json:"languageCode"`
	SysDictTypeId  uint64                `gorm:"column:sys_dict_type_id;type:bigint;not null;comment:字典类型ID;index:idx_sys_dict_entry_dict_type_id;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"sysDictTypeId"`
	SysDictType    *SysDictType          `gorm:"foreignKey:SysDictTypeId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"sysDictType"`
}

func (SysDictEntry) TableName() string {
	return "sys_dict_entry"
}
