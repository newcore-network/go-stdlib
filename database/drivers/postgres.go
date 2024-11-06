package drivers

import (
	"fmt"

	"github.com/styerr-development/libs/configuration"
	"github.com/styerr-development/libs/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConnection struct{}

func (p *PostgresConnection) Connect(cfg configuration.GeneralConfig) (database.Conn, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Database, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return database.Conn{}, err
	}

	return database.Conn{Gorm: db}, nil
}
