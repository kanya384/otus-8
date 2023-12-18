package config

import (
	"fmt"

	"github.com/vrischmann/envconfig"
)

type Config struct {
	App struct {
		ServiceName string
		Port        int
		Env         string `envconfig:"default=prod"`
	}

	Log struct {
		Level string
	}

	PG struct {
		User    string
		Pass    string
		Port    string
		Host    string
		PoolMax int
		DbName  string
		Timeout int
	}

	Redis struct {
		Host string
		Port int
	}

	Jaeger struct {
		Endpoint string
	}
}

func InitConfig(prefix string) (*Config, error) {
	conf := &Config{}
	if err := envconfig.InitWithPrefix(conf, prefix); err != nil {
		return nil, fmt.Errorf("init config error: %w", err)
	}

	return conf, nil
}

func (c Config) PostgresDsn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.PG.User, c.PG.Pass, c.PG.Host, c.PG.Port, c.PG.DbName)
}

func (c Config) RedisAddress() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}
