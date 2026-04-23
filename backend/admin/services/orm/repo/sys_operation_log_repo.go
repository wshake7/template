package repo

import (
    "admin/services/orm/models"
    "go-common/mapper"
    "orm-crud/gorm"
)

type sysOperationLogRepo[T, R any] struct {
    *gorm.Repository[T, R]
}

var SysOperationLogRepo *sysOperationLogRepo[models.SysOperationLog, models.SysOperationLog]

func init() {
    repository := gorm.NewRepository(mapper.NewCopierMapper[models.SysOperationLog, models.SysOperationLog]())
    SysOperationLogRepo = &sysOperationLogRepo[models.SysOperationLog, models.SysOperationLog]{
        Repository: repository,
    }
}