package repo

import (
    "admin/services/orm/models"
    "admin/services/orm/query"
    "go-common/mapper"
    "gorm.io/gen"
    "gorm.io/gen/field"
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

func (sysDictEntryRepo[T, R]) UpdateMap(m map[field.Expr]any, conds ...gen.Condition) (gen.ResultInfo, error) {
    d := make(map[string]any, len(m))
    for k, v := range m {
        d[k.ColumnName().String()] = v
    }
    q := query.SysDictEntry
    return q.Where(conds...).Updates(d)
}