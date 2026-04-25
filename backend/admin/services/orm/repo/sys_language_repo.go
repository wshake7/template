package repo

import (
    "admin/services/orm/models"
    "admin/services/orm/query"
    "go-common/mapper"
    "gorm.io/gen"
    "orm-crud/gormc"
)

type sysLanguageRepo[T, R any] struct {
    *gormc.Repository[T, R]
}

var SysLanguageRepo *sysLanguageRepo[models.SysLanguage, models.SysLanguage]

func init() {
    repository := gormc.NewRepository(mapper.NewCopierMapper[models.SysLanguage, models.SysLanguage]())
    SysLanguageRepo = &sysLanguageRepo[models.SysLanguage, models.SysLanguage]{
        Repository: repository,
    }
}

func (sysLanguageRepo[T, R]) UpdateMap(m map[string]any, conds ...gen.Condition) (gen.ResultInfo, error) {
    if len(m) == 0 {
        return gen.ResultInfo{}, nil
    }
    q := query.SysLanguage
    return q.Where(conds...).Updates(m)
}

func (sysLanguageRepo[T, R]) UpdateNoNilMap(m map[string]any, conds ...gen.Condition) (gen.ResultInfo, error) {
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
    q := query.SysLanguage
    return q.Where(conds...).Updates(d)
}
