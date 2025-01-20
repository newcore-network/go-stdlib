package stdlib

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type StdLibConfiguration struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBDatabase string
	DBPort     int
	DBSSLMode  string

	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int

	DevMode bool
}

// LoadCfg loads the configuration from the specified file or defaults to ".env".
// It returns a StdLibConfiguration instance.
// If an error occurs during loading, it logs the error and continues with the environment variables already set.
func LoadCfg(file ...string) StdLibConfiguration {
	defaultFile := ".env"
	if len(file) > 0 && file[0] != "" {
		defaultFile = file[0]
	}

	if _, err := os.Stat(defaultFile); err == nil {
		if loadErr := godotenv.Overload(defaultFile); loadErr != nil {
			log.Panicf("Error loading environment file: %v", loadErr)
		}
	} else if len(defaultFile) > 0 {
		log.Panicf("Specified environment file '%s' does not exist", file)
	}

	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		log.Panicf("Error parsing POSTGRES_PORT: %v", err)
	}
	redisPort, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		log.Panicf("Error parsing REDISPORT: %v", err)
	}
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Panicf("Error parsing REDISDB: %v", err)
	}

	SSLMode := "disable"
	if os.Getenv("POSTGRES_SSLMODE") == "enable" {
		SSLMode = "enable"
	}

	return StdLibConfiguration{
		DBHost:     os.Getenv("POSTGRES_HOST"),
		DBUser:     os.Getenv("POSTGRES_USER"),
		DBPassword: os.Getenv("POSTGRES_PASSWORD"),
		DBDatabase: os.Getenv("POSTGRES_DATABASE"),
		DBPort:     port,
		DBSSLMode:  SSLMode,

		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     redisPort,
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       redisDB,

		DevMode: os.Getenv("DEV_MODE") == "true",
	}
}
