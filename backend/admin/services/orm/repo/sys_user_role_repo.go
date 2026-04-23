package repo

import (
    "admin/services/orm/models"
    "admin/services/orm/query"
    "go-common/mapper"
    "gorm.io/gen"
    "gorm.io/gen/field"
    "orm-crud/gorm"
)

type sysUserRoleRepo[T, R any] struct {
    *gorm.Repository[T, R]
}

var SysUserRoleRepo *sysUserRoleRepo[models.SysUserRole, models.SysUserRole]

func init() {
    repository := gorm.NewRepository(mapper.NewCopierMapper[models.SysUserRole, models.SysUserRole]())
    SysUserRoleRepo = &sysUserRoleRepo[models.SysUserRole, models.SysUserRole]{
        Repository: repository,
    }
}

func (sysUserRoleRepo[T, R]) UpdateMap(m map[field.Expr]any, conds ...gen.Condition) (gen.ResultInfo, error) {
    d := make(map[string]any, len(m))
    for k, v := range m {
        d[k.ColumnName().String()] = v
    }
    q := query.SysUserRole
    return q.Where(conds...).Updates(d)
}