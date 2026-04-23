package repo

import (
    "admin/services/orm/models"
    "go-common/mapper"
    "orm-crud/gorm"
)

type sysUserRoleRepo[T, R any] struct {
    *gorm.Repository[T, R]
}

var SysUserRoleRepo *sysUserRoleRepo[models.SysUserRole, models.SysUserRole]

func init() {
    repository := gorm.NewRepository(mapper.NewCopierMapper[models.SysUserRole, models.SysUserRole]())
    SysUserRoleRepo = &sysUserRoleRepo[models.SysUserRole, models.SysUserRole]{
        Repository: repository,
    }
}