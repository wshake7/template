package repo

import (
    "admin/services/orm/models"
    "admin/services/orm/query"
    "go-common/mapper"
    "gorm.io/gen"
    "orm-crud/gormc"
)

type sysDictTypeRepo[T, R any] struct {
    *gormc.Repository[T, R]
}

var SysDictTypeRepo *sysDictTypeRepo[models.SysDictType, models.SysDictType]

func init() {
    repository := gormc.NewRepository(mapper.NewCopierMapper[models.SysDictType, models.SysDictType]())
    SysDictTypeRepo = &sysDictTypeRepo[models.SysDictType, models.SysDictType]{
        Repository: repository,
    }
}

func (sysDictTypeRepo[T, R]) UpdateMap(m map[string]any, conds ...gen.Condition) (gen.ResultInfo, error) {
    if len(m) == 0 {
        return gen.ResultInfo{}, nil
    }
    q := query.SysDictType
    return q.Where(conds...).Updates(m)
}

func (sysDictTypeRepo[T, R]) UpdateNoNilMap(m map[string]any, conds ...gen.Condition) (gen.ResultInfo, error) {
    if len(m) == 0 {
        return gen.ResultInfo{}, nil
    }
    d := make(map[string]any, len(m))
    for k, v := range m {
        if v != nil {
            d[k] = v
        }
    }
    if len(d) == 0 {
        return gen.ResultInfo{}, nil
    }
    q := query.SysDictType
    return q.Where(conds...).Updates(d)
}
