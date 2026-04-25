package models

import (
	"gorm.io/gorm"
	"orm-crud/gormc/mixin"
)

func init() {
	Models = append(Models, &SysDictType{})
}

type SysDictType struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.OperatorID
	mixin.IsEnabled
	mixin.SortOrder
	mixin.Description
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at;softDelete:milli;uniqueIndex:idx_sys_dict_type_type_code_deleted_at,priority:2" json:"deletedAt,omitempty"`
	TypeCode  string          `gorm:"column:type_code;type:varchar(128);not null;uniqueIndex:idx_sys_dict_type_type_code_deleted_at,priority:1;comment:字典类型唯一代码" json:"typeCode,omitempty"`
	TypeName  string          `gorm:"column:type_name;type:varchar(255);not null;comment:字典类型名称" json:"typeName,omitempty"`
	Entries   []SysDictEntry  `gorm:"foreignKey:SysDictTypeId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"entries,omitempty"`
}

// TableName 指定表名
func (SysDictType) TableName() string {
	return "sys_dict_types"
}
