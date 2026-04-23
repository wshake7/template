package services

import (
	"admin/config"
	"admin/services/repo"
	"context"
	"fmt"
	"go.uber.org/zap"
	gormCrud "orm-crud/gorm"
)

type Repo struct {
	repoConf config.RepoConfig
	client   *gormCrud.Client
}

func NewRepo(conf config.RepoConfig) *Repo {
	return &Repo{
		repoConf: conf,
		client:   repo.New(conf),
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
		return "unhealthy", fmt.Errorf("failed to get sql.DB: %w", err)
	}
	if err = sqlDB.PingContext(ctx); err != nil {
		return "unhealthy", fmt.Errorf("database ping failed: %w", err)
	}
	stats := sqlDB.Stats()
	state := fmt.Sprintf(
		"healthy | open=%d idle=%d inUse=%d waitCount=%d",
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
