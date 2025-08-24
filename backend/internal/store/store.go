package store

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Open(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name,
	)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
