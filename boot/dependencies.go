package boot

import (
	"gorm.io/gorm"

	"github.com/intermediate-service-ta/internal/storage/dao"
)

type Dependencies struct {
	cfg Config
	db  *gorm.DB
}

func InitDependencies(cfg Config) (*Dependencies, error) {
	db, err := InitDB(cfg.DatabaseConfig)
	db.AutoMigrate(&dao.File{}, &dao.User{}, &dao.Group{}, &dao.ChunkFile{})
	if err != nil {
		return nil, err
	}

	return &Dependencies{
		cfg: cfg,
		db:  db,
	}, nil
}

func (d Dependencies) Config() Config {
	return d.cfg
}

func (d Dependencies) DB() *gorm.DB {
	return d.db
}
