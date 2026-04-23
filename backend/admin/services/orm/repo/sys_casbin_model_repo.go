package repo

import (
    "admin/services/orm/models"
    "go-common/mapper"
    "orm-crud/gorm"
)

type sysCasbinModelRepo[T, R any] struct {
    *gorm.Repository[T, R]
}

var SysCasbinModelRepo *sysCasbinModelRepo[models.SysCasbinModel, models.SysCasbinModel]

func init() {
    repository := gorm.NewRepository(mapper.NewCopierMapper[models.SysCasbinModel, models.SysCasbinModel]())
    SysCasbinModelRepo = &sysCasbinModelRepo[models.SysCasbinModel, models.SysCasbinModel]{
        Repository: repository,
    }
}