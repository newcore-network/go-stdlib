package drivers

import (
	"fmt"

	"github.com/newcore-network/libs/configuration"
	"github.com/newcore-network/libs/database"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MariaDBConnection struct{}

func (m *MariaDBConnection) Connect(cfg configuration.GeneralConfig) (database.Conn, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return database.Conn{}, err
	}

	return database.Conn{Gorm: db}, nil
}
