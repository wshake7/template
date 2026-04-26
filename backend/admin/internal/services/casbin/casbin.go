package casbin

import (
	"admin/internal/services/orm/query"
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
	result, err := sysCasbinModel.Where(sysCasbinModel.IsEnabled.Is(true)).Select(sysCasbinModel.Content).First()
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
