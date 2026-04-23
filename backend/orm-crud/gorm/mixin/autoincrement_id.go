package mixin

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// AutoIncrementID 是 GORM 可复用的 mixin，包含自增主键字段。
// 在模型中嵌入：
//
//	type User struct {
//	    mixin.AutoIncrementID
//	    Name string
//	}
type AutoIncrementID struct {
	ID uint64 `gorm:"column:id;type:bigint;primaryKey;autoIncrement" json:"id"`
}

func (AutoIncrementID) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	// 保持默认行为或根据 db.Dialector 返回特定类型
	return ""
}
