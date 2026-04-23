package repo

import (
    "admin/services/orm/models"
    "admin/services/orm/query"
    "go-common/mapper"
    "gorm.io/gen"
    "gorm.io/gen/field"
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

func (sysOperationLogRepo[T, R]) UpdateMap(m map[field.Expr]any, conds ...gen.Condition) (gen.ResultInfo, error) {
    d := make(map[string]any, len(m))
    for k, v := range m {
        d[k.ColumnName().String()] = v
    }
    q := query.SysOperationLog
    return q.Where(conds...).Updates(d)
}