package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server  ServerConfig
	JWT     JWTConfig
	Postgre PostgreConfig
	OpenID  OpenIDConfig
	S3       S3Config
	Redis    RedisConfig
}

type ServerConfig struct {
	InternalPort string `env:"INTERNAL_PORT"`
	ExternalPort string `env:"EXTERNAL_PORT"`
	RunMode      string `env:"GIN_MODE"`
}

type JWTConfig struct {
	AccessTokenExpireMinutes int    `env:"JWT_ACCESS_TOKEN_EXPIRY_MINUTES"`
	RefreshTokenExpireHour   int    `env:"JWT_REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret        string `env:"JWT_ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret       string `env:"JWT_REFRESH_TOKEN_SECRET"`
}

type RedisConfig struct {
	Host               string        `env:"REDIS_HOST"`
	Port               string        `env:"REDIS_PORT"`
	Password           string        `env:"REDIS_PASSWORD"`
	Db                 string        `env:"REDIS_DB"`
	DialTimeout        time.Duration `env:"REDIS_DIAL_TIMEOUT"`
	ReadTimeout        time.Duration `env:"REDIS_READ_TIMEOUT"`
	WriteTimeout       time.Duration `env:"REDIS_WRITE_TIMEOUT"`
	IdleCheckFrequency time.Duration `env:"REDIS_IDLE_CHECK_FREQ"`
	PoolSize           int           `env:"REDIS_POOLSIZE"`
	PoolTimeout        time.Duration `env:"REDIS_POOL_TIMEOUT"`
}

type PostgreConfig struct {
	PG_Username      string `env:"PG_USERNAME"`
	PG_Password      string `env:"PG_PASSWORD"`
	PG_PoolMax       int    `env:"PG_POOLMAX"`
	PG_Port          string `env:"PG_PORT"`
	PG_Host          string `env:"PG_HOST"`
	PG_DB            string `env:"PG_DBNAME"`
	PG_MIGRATION_URL string `env:"PG_MIGRATION_URL"`
}

type OpenIDConfig struct {
	ClientId     string `env:"OPENID_CLIENT_ID"`
	ClientSecret string `env:"OPENID_CLIENT_SECRET"`
	RedirectUrl  string `env:"OPENID_REDIRECT_URL"`
	IssuerUrl    string `env:"OPENID_ISSUER_URL"`
}

type S3Config struct{
	Baseurl string  `env:"S3_BASEURL"`
	Path string  `env:"S3_PATH"`
}

func NewConfig() *Config {
	cfg := &Config{}
	cwd := projectRoot()
	envFilePath := filepath.Join(cwd, ".env")

	err := readEnv(envFilePath, cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

func readEnv(envFilePath string, cfg *Config) error {
	envFileExists := checkFileExists(envFilePath)

	if envFileExists {
		err := cleanenv.ReadConfig(envFilePath, cfg)
		if err != nil {
			return fmt.Errorf("config error: %w", err)
		}
	} else {
		err := cleanenv.ReadEnv(cfg)
		if err != nil {

			if _, statErr := os.Stat(envFilePath + ".example"); statErr == nil {
				return fmt.Errorf("missing environmentvariables: %w\n\nprovide all required environment variables or rename and update .env.example to .env for convinience", err)
			}

			return err
		}
	}
	return nil
}

func checkFileExists(fileName string) bool {
	envFileExists := false
	if _, err := os.Stat(fileName); err == nil {
		envFileExists = true
	}
	return envFileExists
}

func projectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(b)

	return filepath.Join(projectRoot, "..")
}
