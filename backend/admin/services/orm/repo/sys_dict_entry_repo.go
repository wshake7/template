package repo

import (
    "admin/services/orm/models"
    "admin/services/orm/query"
    "go-common/mapper"
    "gorm.io/gen"
    "orm-crud/gormc"
)

type sysDictEntryRepo[T, R any] struct {
    *gormc.Repository[T, R]
}

var SysDictEntryRepo *sysDictEntryRepo[models.SysDictEntry, models.SysDictEntry]

func init() {
    repository := gormc.NewRepository(mapper.NewCopierMapper[models.SysDictEntry, models.SysDictEntry]())
    SysDictEntryRepo = &sysDictEntryRepo[models.SysDictEntry, models.SysDictEntry]{
        Repository: repository,
    }
}

func (sysDictEntryRepo[T, R]) UpdateMap(m map[string]any, conds ...gen.Condition) (gen.ResultInfo, error) {
    if len(m) == 0 {
        return gen.ResultInfo{}, nil
    }
    q := query.SysDictEntry
    result, err := q.Where(conds...).Updates(m)
    if err != nil {
        return result, err
    }
    if result.Error !=nil {
        return result, result.Error
    }
    return result, err
}

func (sysDictEntryRepo[T, R]) UpdateNoNilMap(m map[string]any, conds ...gen.Condition) (gen.ResultInfo, error) {
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
    q := query.SysDictEntry
    result, err := q.Where(conds...).Updates(d)
    if err != nil {
        return result, err
    }
    if result.Error !=nil {
        return result, result.Error
    }
    return result, err
}
