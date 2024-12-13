package stdlib

import (
	"log"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

type StdLibConfiguration struct {
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

// LoadCfg loads the configuration from the specified file.
// If the file is not provided, it defaults to ".env".
// It return a GeneralConfig instance.
func LoadCfg(file string) StdLibConfiguration {
	var err error
	if file == "" {
		file = ".env"
	}

	err = godotenv.Load(file)
	if err != nil {
		log.Fatalf("%s -> %v", ERRLoading, err)
	}

	port, err := strconv.Atoi(os.Getenv(POSTGRES_PORT))
	if err != nil {
		color.New(color.BgRed, color.FgHiWhite).Println(ERRPort)
		return StdLibConfiguration{}
	}

	redisPort, err := strconv.Atoi(os.Getenv(REDISPORT))
	if err != nil {
		log.Println(ERRPort)
		return StdLibConfiguration{}
	}

	redisDB, err := strconv.Atoi(os.Getenv(REDISDB))
	if err != nil {
		log.Println(ERRPort)
		return StdLibConfiguration{}
	}

	var SSLMode string
	if os.Getenv(POSTGRES_SSLMODE) == ENABLE {
		SSLMode = ENABLE
	} else {
		SSLMode = DISABLE
	}

	return StdLibConfiguration{
		Host:     os.Getenv(POSTGRES_HOST),
		User:     os.Getenv(POSTGRES_USER),
		Password: os.Getenv(POSTGRES_PASSWORD),
		Database: os.Getenv(POSTGRES_DATABASE),
		Port:     port,
		SSLMode:  SSLMode,

		RedisHost: os.Getenv(REDISHOST),
		RedisPort: redisPort,
		RedisDB:   redisDB,
	}
}
