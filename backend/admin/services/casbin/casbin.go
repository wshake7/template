package casbin

import (
	"admin/services/orm"
	"admin/services/orm/query"
	"admin/services/orm/repo"
	"context"
	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

var E casbin.IEnforcer

var Adapter *gormadapter.Adapter

func New(db *gorm.DB) {
	var err error
	Adapter, err = gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(err)
	}
	sysCasbinModel := query.SysCasbinModel
	result, err := repo.SysCasbinModelRepo.Get(context.Background(), orm.DB().Where(sysCasbinModel.IsEnabled.Is(true)), sysCasbinModel.Content)
	if err != nil {
		panic(err)
	}
	m, err := model.NewModelFromString(result.Content)
	if err != nil {
		panic(err)
	}
	E, err = casbin.NewSyncedCachedEnforcer(m, Adapter)
	if err != nil {
		panic(err)
	}
}
