package models

import "orm-crud/gorm/mixin"

func init() {
	Models = append(Models, &SysCasbinModel{})
}

type SysCasbinModel struct {
	mixin.AutoIncrementID
	mixin.TimeAt
	mixin.OperatorID
	mixin.Status
	mixin.Remark
	Name    string `gorm:"column:name;type:varchar(255);not null;comment:模型名称"`
	Content string `gorm:"column:content;type:text;not null;comment:模型内容"`
}

func (SysCasbinModel) TableName() string {
	return "sys_casbin_model"
}
