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
	Server    ServerConfig
	JWT       JWTConfig
	Postgre   PostgreConfig
	Redis     RedisConfig
	Password  PasswordConfig
	Minio     MinioConfig
	Cors      CorsConfig
	Logger    LoggerConfig
	Otp       OtpConfig
	Oidc      OIDC
	Imagor    ImagorConfig
	RabbitMq  RabbitMQConfig
	TypeSense TypeSenseConfig
}

type ServerConfig struct {
	InternalPort string `env:"INTERNAL_PORT"`
	ExternalPort string `env:"EXTERNAL_PORT"`
	RunMode      string `env:"APP_MODE"`
}

type JWTConfig struct {
	AccessTokenExpireHour  int    `env:"JWT_ACCESS_TOKEN_EXPIRY_HOUR" env-default:"5"`
	RefreshTokenExpireHour int    `env:"JWT_REFRESH_TOKEN_EXPIRY_HOUR" env-default:"168"`
	AccessTokenSecret      string `env:"JWT_ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `env:"JWT_REFRESH_TOKEN_SECRET"`
}

type LoggerConfig struct {
	FilePath string `env:"JWT_REFRESH_SECRET"`
	Encoding string `env:"JWT_REFRESH_SECRET"`
	Level    string `env:"LOG_LEVEL"`
	Logger   string `env:"LOGGER"`
}

type RedisConfig struct {
	Host               string        `env:"REDIS_HOST"`
	Port               string        `env:"REDIS_PORT"`
	Password           string        `env:"REDIS_PASSWORD"`
	Db                 int           `env:"REDIS_DB"`
	DialTimeout        time.Duration `env:"REDIS_DIAL_TIMEOUT"`
	ReadTimeout        time.Duration `env:"REDIS_READ_TIMEOUT"`
	WriteTimeout       time.Duration `env:"REDIS_WRITE_TIMEOUT"`
	IdleCheckFrequency time.Duration `env:"REDIS_IDLE_CHECK_FREQ"`
	PoolSize           int           `env:"REDIS_POOLSIZE"`
	PoolTimeout        time.Duration `env:"REDIS_POOL_TIMEOUT"`
}

type ImagorConfig struct {
	Host string `env:"IMAGOR_HOST"`
	Port string `env:"IMAGOR_PORT"`
}

type PasswordConfig struct {
	IncludeChars     bool
	IncludeDigits    bool
	MinLength        int
	MaxLength        int
	IncludeUppercase bool
	IncludeLowercase bool
}

type PostgreConfig struct {
	PG_Username      string `env:"PG_USERNAME"`
	PG_Password      string `env:"PG_PASSWORD"`
	PG_PoolMax       int    `env:"PG_POOLMAX"`
	PG_Port          string `env:"PG_PORT"`
	PG_Host          string `env:"PG_HOST"`
	PG_Name          string `env:"PG_NAME"`
	PG_MIGRATION_URL string `env:"PG_MIGRATION_URL"`
}

type MinioConfig struct {
	EndPoint  string `env:"MINIO_ENDPOINT"`
	AccessKey string `env:"MINIO_ACCESS_KEY"`
	SecretKey string `env:"MINIO_SECRET_KEY"`
	Region    string `env:"MINIO_REGION"`
	Bucket1   string `env:"MINIO_BUCKET1"`
	Bucket2   string `env:"MINIO_BUCKET2"`
	SSL       bool   `env:"MINIO_SSL"`
}

type CorsConfig struct {
	AllowOrigins string
}

type OtpConfig struct {
	ExpireTime time.Duration
	Digits     int
	Limiter    time.Duration
}

type OIDC struct {
	ClientId     string `env:"OIDC_CLIENT_ID"`
	ClientSecret string `env:"OIDC_CLIENT_SECRET"`
	RedirectUrl  string `env:"OIDC_REDIRECT_URL"`
	IssuerUrl    string `env:"OIDC_ISSUER_URL"`
}

type RabbitMQConfig struct {
	Host     string `env:"RABBITMQ_HOST"`
	Port     string `env:"RABBITMQ_PORT"`
	User     string `env:"RABBITMQ_USER"`
	Password string `env:"RABBITMQ_PASSWORD"`
}

type TypeSenseConfig struct {
	Host           string `env:"TYPESENSE_HOST"`
	Port           string `env:"TYPESENSE_PORT"`
	ApiKey         string `env:"TYPESENSE_API_KEY"`
	MeiliImageHost string `env:"TYPSENSE_IMAGE_HOST"`
}

func NewConfig() *Config {
	cfg := &Config{}
	cwd := projectRoot()
	envFilePath := cwd + ".env"

	if checkFileExists(envFilePath) {
		if err := cleanenv.ReadConfig(envFilePath, cfg); err != nil {
			panic(fmt.Errorf("config file error: %w", err))
		}
	}

	// 2. Baca/Timpa dari Environment Variable (ENV) sistem
	//    Ini akan menimpa nilai dari file .env (jika ada) dan mengambil nilai dari Podman
	if err := cleanenv.ReadEnv(cfg); err != nil {
		// Logika ini harus disesuaikan, tapi intinya: selalu baca ENV sistem
		panic(fmt.Errorf("environment variable error: %w", err))
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

	return projectRoot + "/../"
}
