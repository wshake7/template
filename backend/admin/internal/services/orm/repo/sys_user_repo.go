package repo

import (
    "admin/internal/services/orm/models"
    "admin/internal/services/orm/query"
    "go-common/mapper"
    "gorm.io/gen"
    "orm-crud/gormc"
)

type sysUserRepo[T, R any] struct {
    *gormc.Repository[T, R]
}

var SysUserRepo *sysUserRepo[models.SysUser, models.SysUser]

func init() {
    repository := gormc.NewRepository(mapper.NewCopierMapper[models.SysUser, models.SysUser]())
    SysUserRepo = &sysUserRepo[models.SysUser, models.SysUser]{
        Repository: repository,
    }
}

func (sysUserRepo[T, R]) UpdateMap(m map[string]any, conds ...gen.Condition) (gen.ResultInfo, error) {
    if len(m) == 0 {
        return gen.ResultInfo{}, nil
    }
    q := query.SysUser
    result, err := q.Where(conds...).Updates(m)
    if err != nil {
        return result, err
    }
    if result.Error !=nil {
        return result, result.Error
    }
    return result, err
}

func (sysUserRepo[T, R]) UpdateNoNilMap(m map[string]any, conds ...gen.Condition) (gen.ResultInfo, error) {
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
    q := query.SysUser
    result, err := q.Where(conds...).Updates(d)
    if err != nil {
        return result, err
    }
    if result.Error !=nil {
        return result, result.Error
    }
    return result, err
}
