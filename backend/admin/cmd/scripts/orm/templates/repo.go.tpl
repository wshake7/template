package repo

import (
    "{{.ModuleName}}/services/orm/models"
    "go-common/mapper"
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