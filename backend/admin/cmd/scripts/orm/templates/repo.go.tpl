package repo

import (
    "{{.ModuleName}}/services/orm/models"
    "{{.ModuleName}}/services/orm/query"
    "go-common/mapper"
    "gorm.io/gen"
    "orm-crud/gormc"
)

type {{.ModelName | toLower}}Repo[T, R any] struct {
    *gormc.Repository[T, R]
}

var {{.ModelName}}Repo *{{.ModelName | toLower}}Repo[models.{{.ModelName}}, models.{{.ModelName}}]

func init() {
    repository := gormc.NewRepository(mapper.NewCopierMapper[models.{{.ModelName}}, models.{{.ModelName}}]())
    {{.ModelName}}Repo = &{{.ModelName | toLower}}Repo[models.{{.ModelName}}, models.{{.ModelName}}]{
        Repository: repository,
    }
}

func ({{.ModelName | toLower}}Repo[T, R]) UpdateMap(m map[string]any, conds ...gen.Condition) (gen.ResultInfo, error) {
    if len(m) == 0 {
        return gen.ResultInfo{}, nil
    }
    q := query.{{.ModelName}}
    result, err := q.Where(conds...).Updates(m)
    if err != nil {
        return result, err
    }
    if result.Error !=nil {
        return result, result.Error
    }
    return result, err
}

func ({{.ModelName | toLower}}Repo[T, R]) UpdateNoNilMap(m map[string]any, conds ...gen.Condition) (gen.ResultInfo, error) {
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
    q := query.{{.ModelName}}
    result, err := q.Where(conds...).Updates(d)
    if err != nil {
        return result, err
    }
    if result.Error !=nil {
        return result, result.Error
    }
    return result, err
}
