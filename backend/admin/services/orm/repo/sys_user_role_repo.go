package repo

import (
    "admin/services/orm/models"
    "admin/services/orm/query"
    "go-common/mapper"
    "gorm.io/gen"
    "orm-crud/gormc"
)

type sysUserRoleRepo[T, R any] struct {
    *gormc.Repository[T, R]
}

var SysUserRoleRepo *sysUserRoleRepo[models.SysUserRole, models.SysUserRole]

func init() {
    repository := gormc.NewRepository(mapper.NewCopierMapper[models.SysUserRole, models.SysUserRole]())
    SysUserRoleRepo = &sysUserRoleRepo[models.SysUserRole, models.SysUserRole]{
        Repository: repository,
    }
}

func (sysUserRoleRepo[T, R]) UpdateMap(m map[string]any, conds ...gen.Condition) (gen.ResultInfo, error) {
    if len(m) == 0 {
        return gen.ResultInfo{}, nil
    }
    q := query.SysUserRole
    return q.Where(conds...).Updates(m)
}

func (sysUserRoleRepo[T, R]) UpdateNoNilMap(m map[string]any, conds ...gen.Condition) (gen.ResultInfo, error) {
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
    q := query.SysUserRole
    return q.Where(conds...).Updates(d)
}
