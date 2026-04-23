package repo

import (
    "{{.ModuleName}}/services/orm/models"
    "{{.ModuleName}}/services/orm/query"
    "go-common/mapper"
    "gorm.io/gen"
    "gorm.io/gen/field"
    "orm-crud/gorm"
)

type {{.ModelName | toLower}}Repo[T, R any] struct {
    *gorm.Repository[T, R]
}

var {{.ModelName}}Repo *{{.ModelName | toLower}}Repo[models.{{.ModelName}}, models.{{.ModelName}}]

func init() {
    repository := gorm.NewRepository(mapper.NewCopierMapper[models.{{.ModelName}}, models.{{.ModelName}}]())
    {{.ModelName}}Repo = &{{.ModelName | toLower}}Repo[models.{{.ModelName}}, models.{{.ModelName}}]{
        Repository: repository,
    }
}

func ({{.ModelName | toLower}}Repo[T, R]) UpdateMap(m map[field.Expr]any, conds ...gen.Condition) (gen.ResultInfo, error) {
    d := make(map[string]any, len(m))
    for k, v := range m {
        d[k.ColumnName().String()] = v
    }
    q := query.{{.ModelName}}
    return q.Where(conds...).Updates(d)
}