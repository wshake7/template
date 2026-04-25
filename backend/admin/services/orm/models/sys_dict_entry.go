package models

import (
	"orm-crud/gormc/mixin"
)

func init() {
	Models = append(Models, &SysDictEntry{})
}

// SysDictEntry 对应表 sys_dict_entries
type SysDictEntry struct {
	mixin.AutoIncrementID
	mixin.TimeAt
	mixin.OperatorID
	mixin.SortOrder
	mixin.IsEnabled
	mixin.Remark
	EntryLabel    string       `gorm:"column:entry_label;type:varchar(255);not null;comment:字典项的显示标签" json:"entryLabel,omitempty"`
	EntryValue    string       `gorm:"column:entry_value;type:varchar(255);not null;comment:字典项的实际值" json:"entryValue,omitempty"`
	NumericValue  int32        `gorm:"column:numeric_value;type:int;not null;default:0;comment:数值型值" json:"numericValue,omitempty"`
	LanguageCode  string       `gorm:"column:language_code;type:varchar(32);not null:default:'';comment:语言代码" json:"languageCode,omitempty"`
	SysDictTypeId uint64       `gorm:"column:sys_dict_type_id;type:bigint;not null;comment:字典类型ID;index:idx_sys_dict_entry_dict_type_id;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"sysDictTypeId,omitempty"`
	SysDictType   *SysDictType `gorm:"foreignKey:SysDictTypeId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"sysDictType,omitempty"`
}

func (SysDictEntry) TableName() string {
	return "sys_dict_entries"
}
