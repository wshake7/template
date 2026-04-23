package repo

import (
	"admin/services/orm/models"
	"go-common/mapper"
	"orm-crud/gorm"
)

var RoleRepo *gorm.Repository[models.SysRole, models.SysRole]

func init() {
	RoleRepo = gorm.NewRepository(mapper.NewCopierMapper[models.SysRole, models.SysRole]())
}
