package services

import (
	"admin/config"
	"admin/services/orm"
	"context"
	"fmt"
	"go.uber.org/zap"
	gormCrud "orm-crud/gorm"
)

type Repo struct {
	ormConf config.OrmConfig
	client  *gormCrud.Client
}

func NewOrm(conf config.OrmConfig) *Repo {
	return &Repo{
		ormConf: conf,
		client:  orm.New(conf),
	}
}

func (o *Repo) Start(ctx context.Context) error {
	return nil
}

func (o *Repo) String() string {
	return "repo"
}

func (o *Repo) State(ctx context.Context) (string, error) {
	sqlDB, err := o.client.DB.DB()
	if err != nil {
		return "UNHEALTHY", fmt.Errorf("failed to get sql.DB: %w", err)
	}
	if err = sqlDB.PingContext(ctx); err != nil {
		return "UNHEALTHY", fmt.Errorf("database ping failed: %w", err)
	}
	stats := sqlDB.Stats()
	state := fmt.Sprintf(
		"HEALTHY | open=%d idle=%d inUse=%d waitCount=%d",
		stats.OpenConnections,
		stats.Idle,
		stats.InUse,
		stats.WaitCount,
	)
	return state, nil
}

func (o *Repo) Terminate(ctx context.Context) error {
	zap.L().Info("database close start")
	sqlDB, err := o.client.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	if err = sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}
	zap.L().Info("database closed")
	return nil
}
