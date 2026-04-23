package repo

import (
    "admin/services/orm/models"
    "go-common/mapper"
    "orm-crud/gorm"
)

type sysDictTypeRepo[T, R any] struct {
    *gorm.Repository[T, R]
}

var SysDictTypeRepo *sysDictTypeRepo[models.SysDictType, models.SysDictType]

func init() {
    repository := gorm.NewRepository(mapper.NewCopierMapper[models.SysDictType, models.SysDictType]())
    SysDictTypeRepo = &sysDictTypeRepo[models.SysDictType, models.SysDictType]{
        Repository: repository,
    }
}