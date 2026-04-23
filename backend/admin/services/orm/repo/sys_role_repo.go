package repo

import (
	"admin/services/orm/models"
	"go-common/mapper"
	"orm-crud/gorm"
)

type sysRoleRepo[T, R any] struct {
	*gorm.Repository[T, R]
}

var SysRoleRepo *sysRoleRepo[models.SysRole, models.SysRole]

func init() {
	repository := gorm.NewRepository(mapper.NewCopierMapper[models.SysRole, models.SysRole]())
	SysRoleRepo = &sysRoleRepo[models.SysRole, models.SysRole]{
		Repository: repository,
	}
}
