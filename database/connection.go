package database

import (
	"errors"
	"log"
	"time"

	"github.com/styerr-development/libs/configuration"
	"gorm.io/gorm"
)

type Connection interface {
	Connect(cfg configuration.GeneralConfig) (Conn, error)
}

type Conn struct {
	Gorm *gorm.DB
}

func (c Conn) GetDB() *gorm.DB {
	return c.Gorm
}

func NewConnection(driver Connection, cfg configuration.GeneralConfig, modelsToMigrate ...interface{}) (*gorm.DB, error) {
	var err error
	var conn Conn

	for count := 0; count < 5; count++ {
		conn, err = driver.Connect(cfg)
		if err == nil {
			break
		}

		log.Printf("\nstep : %d, error: %s", count+1, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		return nil, errors.New("connection failed after multiple attempts")
	}

	if len(modelsToMigrate) > 0 {
		err = conn.Gorm.AutoMigrate(modelsToMigrate...)
		if err != nil {
			return nil, err
		}
	}

	return conn.Gorm, nil
}
