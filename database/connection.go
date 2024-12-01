package database

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/styerr-development/libs/configuration"
	"gorm.io/gorm"
)

type Connection interface {
	Connect(cfg configuration.GeneralConfig) (Conn, error)
}

type Conn struct {
	Gorm *gorm.DB
}

// DBWrapper wraps the Gorm DB connection and provides additional methods for managing enums.
type DBWrapper struct {
	Gorm *gorm.DB
}

func (c Conn) GetDB() *gorm.DB {
	return c.Gorm
}

// NewConnection establishes a connection to the database with retry logic and wraps it in a DBWrapper.
func NewConnection(driver Connection, cfg configuration.GeneralConfig) (*DBWrapper, error) {
	var err error
	var conn Conn

	// Retry logic for establishing the connection
	for count := 0; count < 5; count++ {
		conn, err = driver.Connect(cfg)
		if err == nil {
			break
		}

		color.New(color.FgRed).Printf("step : %d, error: %s\n", count+1, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		return nil, errors.New("connection failed after multiple attempts")
	}

	color.New(color.FgRed).Println("Database connection established", nil)

	return &DBWrapper{Gorm: conn.Gorm}, nil
}

// MigrateEnums migrate enums to the database. If the enum type doesn't exist, it creates it.
func (db *DBWrapper) MigrateEnums(enumTypeName string, values []string) *DBWrapper {
	if len(values) == 0 {
		color.New(color.FgRed).Printf("values for enum '%s' are empty, skipping", enumTypeName)
		return db
	}

	// Prepare the SQL to create the enum type
	valueList := "'" + strings.Join(values, "', '") + "'"
	sql := fmt.Sprintf(`
		DO $$ BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = '%s') THEN
				CREATE TYPE %s AS ENUM (%s);
			END IF;
		END $$;
	`, enumTypeName, enumTypeName, valueList)

	// Execute the SQL
	if err := db.Gorm.Exec(sql).Error; err != nil {
		color.New(color.FgRed).Printf("Error adding enum '%s': %v", enumTypeName, err)
	}

	color.New(color.FgHiGreen).Printf("Enum '%s' added", enumTypeName)
	return db
}

// Migrate adds the specified models to the database.
func (db *DBWrapper) Migrate(models ...interface{}) *DBWrapper {
	if len(models) == 0 {
		color.New(color.FgYellow).Println("No models provided for migration, skipping...")
		return db
	}

	if err := db.Gorm.AutoMigrate(models...); err != nil {
		color.New(color.FgRed).Printf("Error migrating models: %v\n", err)
	}

	color.New(color.FgHiGreen).Println("Models migrated")
	return db
}

func (db *DBWrapper) SetConnectionPool(maxOpen, maxIdle int, maxLifetime time.Duration) *DBWrapper {
	sqlDB, err := db.Gorm.DB()
	if err != nil {
		color.New(color.FgRed).Printf("Error getting SQL connection: %v", err)
		os.Exit(1)
	}
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(maxLifetime)

	color.New(color.FgHiGreen).Println("Connection pool configured")
	return db
}

// EnableExtension ensures a PostgreSQL extension is enabled in the database.
func (db *DBWrapper) EnableExtension(extensionName string) *DBWrapper {
	query := `CREATE EXTENSION IF NOT EXISTS "` + extensionName + `";`
	if err := db.Gorm.Exec(query).Error; err != nil {
		color.New(color.FgRed).Printf("Error enabling extension '%s': %v", extensionName, err)
	}
	color.New(color.FgHiGreen).Printf("Extension '%s' enabled", extensionName)
	return db
}

func (db *DBWrapper) EnableUUIDExtension() *DBWrapper {
	db.EnableExtension("pgcrypto")
	return db
}
