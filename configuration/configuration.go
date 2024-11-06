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
	JWTKey   string
}

func GetFromEnvFile(file string) GeneralConfig {
	var err error

	err = godotenv.Load(file)
	if err != nil {
		log.Fatalf("%s -> %v", constants.ERRLoading, err)
	}

	port, err := strconv.Atoi(os.Getenv(constants.PORT))
	if err != nil {
		log.Println(constants.ERRPort)
		return GeneralConfig{}
	}
	var SSLMode string
	if os.Getenv(constants.SSLMODE) == constants.ENABLE {
		SSLMode = constants.ENABLE
	} else {
		SSLMode = constants.DISABLE
	}

	jwtKey := os.Getenv(constants.JWTKEY)
	if jwtKey == "" {
		log.Fatal(constants.JWTKEY)
	}

	return GeneralConfig{
		Host:     os.Getenv(constants.HOST),
		User:     os.Getenv(constants.USER),
		Password: os.Getenv(constants.PASSWORD),
		Database: os.Getenv(constants.DATABASE),
		Port:     port,
		SSLMode:  SSLMode,
		JWTKey:   jwtKey,
	}
}
