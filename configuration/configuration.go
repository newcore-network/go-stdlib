package configuration

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/styerr-development/libs/configuration/constants"
)

type GeneralConfig struct {
	Host     string
	User     string
	Password string
	Database string
	Port     int
	SSLMode  string

	RedisHost string
	RedisPort int
	RedisDB   int
}

func GetFromEnvFile(file string) GeneralConfig {
	var err error

	err = godotenv.Load(file)
	if err != nil {
		log.Fatalf("%s -> %v", constants.ERRLoading, err)
	}

	port, err := strconv.Atoi(os.Getenv(constants.POSTGRES_PORT))
	if err != nil {
		log.Println(constants.ERRPort)
		return GeneralConfig{}
	}

	redisPort, err := strconv.Atoi(os.Getenv(constants.REDISPORT))
	if err != nil {
		log.Println(constants.ERRPort)
		return GeneralConfig{}
	}

	redisDB, err := strconv.Atoi(os.Getenv(constants.REDISDB))
	if err != nil {
		log.Println(constants.ERRPort)
		return GeneralConfig{}
	}

	var SSLMode string
	if os.Getenv(constants.POSTGRES_SSLMODE) == constants.ENABLE {
		SSLMode = constants.ENABLE
	} else {
		SSLMode = constants.DISABLE
	}

	return GeneralConfig{
		Host:     os.Getenv(constants.POSTGRES_HOST),
		User:     os.Getenv(constants.POSTGRES_USER),
		Password: os.Getenv(constants.POSTGRES_PASSWORD),
		Database: os.Getenv(constants.POSTGRES_DATABASE),
		Port:     port,
		SSLMode:  SSLMode,

		RedisHost: os.Getenv(constants.REDISHOST),
		RedisPort: redisPort,
		RedisDB:   redisDB,
	}
}
