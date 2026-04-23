package repo

import (
    "admin/services/orm/models"
    "admin/services/orm/query"
    "go-common/mapper"
    "gorm.io/gen"
    "gorm.io/gen/field"
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

func (sysLanguageRepo[T, R]) UpdateMap(m map[field.Expr]any, conds ...gen.Condition) (gen.ResultInfo, error) {
    d := make(map[string]any, len(m))
    for k, v := range m {
        d[k.ColumnName().String()] = v
    }
    q := query.SysLanguage
    return q.Where(conds...).Updates(d)
}