package stdlib

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func setupEnv() {
	_ = os.Setenv("POSTGRES_HOST", "localhost")
	_ = os.Setenv("POSTGRES_USER", "test_user")
	_ = os.Setenv("POSTGRES_PASSWORD", "test_password")
	_ = os.Setenv("POSTGRES_DATABASE", "test_db")
	_ = os.Setenv("POSTGRES_PORT", "5432")
	_ = os.Setenv("POSTGRES_SSLMODE", "disable")
	_ = os.Setenv("REDISHOST", "localhost")
	_ = os.Setenv("REDISPORT", "6379")
	_ = os.Setenv("REDISDB", "0")
}

func teardownEnv() {
	_ = os.Unsetenv("POSTGRES_HOST")
	_ = os.Unsetenv("POSTGRES_USER")
	_ = os.Unsetenv("POSTGRES_PASSWORD")
	_ = os.Unsetenv("POSTGRES_DATABASE")
	_ = os.Unsetenv("POSTGRES_PORT")
	_ = os.Unsetenv("POSTGRES_SSLMODE")
	_ = os.Unsetenv("REDISHOST")
	_ = os.Unsetenv("REDISPORT")
	_ = os.Unsetenv("REDISDB")
}

func TestLoadCfg_ValidEnv(t *testing.T) {
	setupEnv()
	defer teardownEnv()

	cfg := LoadCfg()

	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, "test_user", cfg.User)
	assert.Equal(t, "test_password", cfg.Password)
	assert.Equal(t, "test_db", cfg.Database)
	assert.Equal(t, 5432, cfg.Port)
	assert.Equal(t, "disable", cfg.SSLMode)
	assert.Equal(t, "localhost", cfg.RedisHost)
	assert.Equal(t, 6379, cfg.RedisPort)
	assert.Equal(t, 0, cfg.RedisDB)
}

func TestLoadCfg_InvalidPort(t *testing.T) {
	_ = os.Setenv("POSTGRES_PORT", "invalid_port")
	defer os.Unsetenv("POSTGRES_PORT")

	assert.Panics(t, func() {
		_ = LoadCfg()
	}, "LoadCfg should panic when the POSTGRES_PORT environment variable is invalid")
}

func TestLoadCfg_ValidEnvFile(t *testing.T) {
	tempFile := ".env.test"
	defer os.Remove(tempFile)

	envContent := `
	POSTGRES_HOST=env_host
	POSTGRES_USER=env_user
	POSTGRES_PASSWORD=env_password
	POSTGRES_DATABASE=env_db
	POSTGRES_PORT=5433
	POSTGRES_SSLMODE=enable
	REDISHOST=env_redis
	REDISPORT=6380
	REDISDB=1
	`

	_ = os.WriteFile(tempFile, []byte(envContent), 0644)
	_ = godotenv.Load(tempFile)

	cfg := LoadCfg(tempFile)

	assert.Equal(t, "env_host", cfg.Host)
	assert.Equal(t, "env_user", cfg.User)
	assert.Equal(t, "env_password", cfg.Password)
	assert.Equal(t, "env_db", cfg.Database)
	assert.Equal(t, 5433, cfg.Port)
	assert.Equal(t, "enable", cfg.SSLMode)
	assert.Equal(t, "env_redis", cfg.RedisHost)
	assert.Equal(t, 6380, cfg.RedisPort)
	assert.Equal(t, 1, cfg.RedisDB)
}

func TestLoadCfg_MissingEnvFile(t *testing.T) {
	file := ".env.test"
	assert.Panics(t, func() {
		_ = LoadCfg(file)
	}, "LoadCfg should panic when the specified environment file does not exist")
}
