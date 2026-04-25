package models

import "orm-crud/gormc/mixin"

func init() {
	Models = append(Models, &SysCasbinModel{})
}

type SysCasbinModel struct {
	mixin.AutoIncrementID
	mixin.TimeAt
	mixin.OperatorID
	mixin.IsEnabled
	mixin.Remark
	Name    string `gorm:"column:name;type:varchar(255);not null;comment:模型名称" json:"name,omitempty"`
	Content string `gorm:"column:content;type:text;not null;comment:模型内容" json:"content,omitempty"`
}

func (SysCasbinModel) TableName() string {
	return "sys_casbin_model"
}
