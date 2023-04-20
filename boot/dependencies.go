package boot

import "gorm.io/gorm"

type Dependencies struct {
	cfg Config
	db  *gorm.DB
}

func InitDependencies(cfg Config) (*Dependencies, error) {
	db, err := InitDB(cfg.DatabaseConfig)
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
