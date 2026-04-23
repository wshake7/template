package repo

import (
    "admin/services/orm/models"
    "go-common/mapper"
    "orm-crud/gorm"
)

type sysLanguageRepo[T, R any] struct {
    *gorm.Repository[T, R]
}

var SysLanguageRepo *sysLanguageRepo[models.SysLanguage, models.SysLanguage]

func init() {
    repository := gorm.NewRepository(mapper.NewCopierMapper[models.SysLanguage, models.SysLanguage]())
    SysLanguageRepo = &sysLanguageRepo[models.SysLanguage, models.SysLanguage]{
        Repository: repository,
    }
}