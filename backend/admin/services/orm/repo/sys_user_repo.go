package repo

import (
    "admin/services/orm/models"
    "admin/services/orm/query"
    "go-common/mapper"
    "gorm.io/gen"
    "gorm.io/gen/field"
    "orm-crud/gorm"
)

type sysUserRepo[T, R any] struct {
    *gorm.Repository[T, R]
}

var SysUserRepo *sysUserRepo[models.SysUser, models.SysUser]

func init() {
    repository := gorm.NewRepository(mapper.NewCopierMapper[models.SysUser, models.SysUser]())
    SysUserRepo = &sysUserRepo[models.SysUser, models.SysUser]{
        Repository: repository,
    }
}

func (sysUserRepo[T, R]) UpdateMap(m map[field.Expr]any, conds ...gen.Condition) (gen.ResultInfo, error) {
    d := make(map[string]any, len(m))
    for k, v := range m {
        d[k.ColumnName().String()] = v
    }
    q := query.SysUser
    return q.Where(conds...).Updates(d)
}