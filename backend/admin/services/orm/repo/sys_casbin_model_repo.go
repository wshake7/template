package repo

import (
    "admin/services/orm/models"
    "admin/services/orm/query"
    "go-common/mapper"
    "gorm.io/gen"
    "gorm.io/gen/field"
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

func (sysCasbinModelRepo[T, R]) UpdateMap(m map[field.Expr]any, conds ...gen.Condition) (gen.ResultInfo, error) {
    d := make(map[string]any, len(m))
    for k, v := range m {
        d[k.ColumnName().String()] = v
    }
    q := query.SysCasbinModel
    return q.Where(conds...).Updates(d)
}