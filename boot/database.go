package boot

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg DatabaseConfig) (*gorm.DB, error) {

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&timeout=10s&parseTime=true", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name)
	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(0)

	err = sqlDB.Ping()
	if err != nil {
		err = fmt.Errorf("error pinging db: %w", err)
	}

	return db, err
}
