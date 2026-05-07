package infra

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

var cfg *EnvConfig

func GetEnvConfig() *EnvConfig {

	if cfg != nil {
		return cfg
	}

	envPath := "../../.env"
	_, envFileErr := os.Stat(envPath)
	if !os.IsNotExist(envFileErr) {
		_ = godotenv.Load(envPath)
	}

	var c EnvConfig
	if cfgErr := env.Parse(&c); cfgErr != nil {
		log.Fatalf("parse env failed: %v", cfgErr)
	}

	hostname, hostNameErr := os.Hostname()
	if hostNameErr != nil || hostname == "" {
		log.Fatalf("get hostname failed: %v", hostNameErr)
	}
	c.ServerName = hostname

	cfg = &c
	return cfg
}

type EnvConfig struct {
	DBPort        int    `env:"DB_PORT" envDefault:"5432"`
	RedisDB       int    `env:"REDIS_DB" envDefault:"0"`
	DBMaxConn     int    `env:"DB_MAX_OPEN" envDefault:"50"`
	DBMaxIdleConn int    `env:"DB_MAX_IDLE" envDefault:"20"`
	DBHost        string `env:"DB_HOST,required"`
	RedisHost     string `env:"REDIS_HOST,required"`
	DBUser        string `env:"DB_USER,required"`
	DBPassword    string `env:"DB_PASSWORD,required"`
	RedisPassword string `env:"REDIS_PASSWORD,required"`
	DBName        string `env:"DB_NAME,required"`
	DBSSLMode     string `env:"DB_SSLMODE" envDefault:"disable"`
	DBTimeZone    string `env:"DB_TIMEZONE" envDefault:"Asia/Taipei"`

	DBConnMaxLifeTime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"30m"`
	DBConnMaxIdleTime time.Duration `env:"DB_CONN_MAX_IDLE" envDefault:"5m"`
	RedisDialTimeout  time.Duration `env:"REDIS_DIAL_TIMEOUT" envDefault:"5s"`

	ServerName string `envDefault:""`
}
