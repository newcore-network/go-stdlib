package stdlib

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Connection defines an interface for database drivers to implement.
// It provides the method to establish a database connection using a given configuration.
type Connection interface {
	Connect(cfg StdLibConfiguration) (Conn, error)
}

// Conn represents a database connection encapsulated with Gorm.
type Conn struct {
	Gorm *gorm.DB
}

// DBWrapper is a wrapper around the Gorm database connection.
// It provides additional methods for managing the database, such as enum migrations and connection pool configuration.
type DBWrapper struct {
	Gorm *gorm.DB
}

// GetDB returns the underlying Gorm database connection.
func (c Conn) GetDB() *gorm.DB {
	return c.Gorm
}

// NewConnection establishes a database connection with retry logic and wraps it in a DBWrapper.
// If the connection fails after multiple attempts (3), it returns an error.
func NewConnection(driver Connection, cfg StdLibConfiguration) (*DBWrapper, error) {
	var err error
	var conn Conn

	// Retry logic for establishing the connection
	for count := 0; count < 5; count++ {
		conn, err = driver.Connect(cfg)
		if err == nil {
			break
		}

		color.New(color.FgRed).Printf("step: %d, error: %s\n", count+1, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		return nil, errors.New("connection failed after multiple attempts")
	}

	color.New(color.FgGreen).Println("Database connection established")

	return &DBWrapper{Gorm: conn.Gorm}, nil
}

// NewRedisConnection creates a new Redis client and pings the server to ensure connectivity.
// If the connection fails, it panics.
func NewRedisConnection(ctx context.Context, cfg StdLibConfiguration) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + strconv.Itoa(cfg.RedisPort),
		Password: cfg.Password,
		DB:       cfg.RedisDB,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic("cannot connect to redis")
	}

	return client
}

// MigrateEnums adds or updates an ENUM type in the PostgreSQL database.
// If the ENUM type does not exist, it creates it with the provided values.
func (db *DBWrapper) MigrateEnums(enumTypeName string, values []string) *DBWrapper {
	if len(values) == 0 {
		color.New(color.FgRed).Printf("values for enum '%s' are empty, skipping...\n", enumTypeName)
		return db
	}

	valueList := "'" + strings.Join(values, "', '") + "'"
	sql := fmt.Sprintf(`
		DO $$ BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = '%s') THEN
				CREATE TYPE %s AS ENUM (%s);
			END IF;
		END $$;
	`, enumTypeName, enumTypeName, valueList)

	if err := db.Gorm.Exec(sql).Error; err != nil {
		color.New(color.FgRed).Printf("Error adding enum '%s': %v\n", enumTypeName, err)
	}

	color.New(color.FgHiGreen).Printf("Enum '%s' added\n", enumTypeName)
	return db
}

// Migrate applies database migrations for the specified models.
// It automatically migrates the database schema based on the provided models.
func (db *DBWrapper) Migrate(models ...interface{}) *DBWrapper {
	if len(models) == 0 {
		color.New(color.FgYellow).Println("No models provided for migration, skipping...")
		return db
	}

	if err := db.Gorm.AutoMigrate(models...); err != nil {
		color.New(color.FgRed).Printf("Error migrating models: %v\n", err)
	}

	var modelNames []string
	for _, model := range models {
		modelNames = append(modelNames, fmt.Sprintf("%T", model))
	}

	color.New(color.FgHiGreen).Printf("Models migrated: %s\n", modelNames)
	return db

}

// SetConnectionPool configures the connection pool settings for the database.
// It allows setting the maximum number of open connections, idle connections, and connection lifetime.
func (db *DBWrapper) SetConnectionPool(maxOpen, maxIdle int, maxLifetime time.Duration) *DBWrapper {
	sqlDB, err := db.Gorm.DB()
	if err != nil {
		color.New(color.FgRed).Printf("Error getting SQL connection: %v\n", err)
		os.Exit(1)
	}
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(maxLifetime)

	color.New(color.FgHiGreen).Println("Connection pool configured")
	return db
}

// EnableExtension ensures a PostgreSQL extension is enabled in the database.
// If the extension does not exist, it creates it.
func (db *DBWrapper) EnableExtension(extensionName string) *DBWrapper {
	query := `CREATE EXTENSION IF NOT EXISTS "` + extensionName + `";`
	if err := db.Gorm.Exec(query).Error; err != nil {
		color.New(color.FgRed).Printf("Error enabling extension '%s': %v\n", extensionName, err)
	}
	color.New(color.FgHiGreen).Printf("Extension '%s' enabled\n", extensionName)
	return db
}

// EnableUUIDExtension is a helper method to enable the 'pgcrypto' extension for UUID generation.
func (db *DBWrapper) EnableUUIDExtension() *DBWrapper {
	db.EnableExtension("pgcrypto")
	return db
}
