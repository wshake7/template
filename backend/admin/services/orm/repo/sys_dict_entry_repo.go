package repo

import (
    "admin/services/orm/models"
    "go-common/mapper"
    "orm-crud/gorm"
)

type sysDictEntryRepo[T, R any] struct {
    *gorm.Repository[T, R]
}

var SysDictEntryRepo *sysDictEntryRepo[models.SysDictEntry, models.SysDictEntry]

func init() {
    repository := gorm.NewRepository(mapper.NewCopierMapper[models.SysDictEntry, models.SysDictEntry]())
    SysDictEntryRepo = &sysDictEntryRepo[models.SysDictEntry, models.SysDictEntry]{
        Repository: repository,
    }
}