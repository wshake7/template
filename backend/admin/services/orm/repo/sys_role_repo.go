package repo

import (
    "admin/services/orm/models"
    "admin/services/orm/query"
    "go-common/mapper"
    "gorm.io/gen"
    "gorm.io/gen/field"
    "orm-crud/gorm"
)

type sysRoleRepo[T, R any] struct {
    *gorm.Repository[T, R]
}

var SysRoleRepo *sysRoleRepo[models.SysRole, models.SysRole]

func init() {
    repository := gorm.NewRepository(mapper.NewCopierMapper[models.SysRole, models.SysRole]())
    SysRoleRepo = &sysRoleRepo[models.SysRole, models.SysRole]{
        Repository: repository,
    }
}

func (sysRoleRepo[T, R]) UpdateMap(m map[field.Expr]any, conds ...gen.Condition) (gen.ResultInfo, error) {
    d := make(map[string]any, len(m))
    for k, v := range m {
        d[k.ColumnName().String()] = v
    }
    q := query.SysRole
    return q.Where(conds...).Updates(d)
}