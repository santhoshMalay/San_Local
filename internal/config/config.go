package config

import (
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

type Config struct {
	SnowflakeNode int64  `env:"SNOWFLAKE_NODE" envDefault:"1"`
	SigningKey    string `env:"SIGNING_KEY,required"`
	LogLevel      string `env:"LOG_LEVEL" envDefault:"info"`
	
	HTTP struct {
		Port               string        `env:"PORT" envDefault:"8080"`
		ReadTimeout        time.Duration `env:"READ_TIMEOUT" envDefault:"5s"`
		WriteTimeout       time.Duration `env:"WRITE_TIMEOUT" envDefault:"5s"`
		MaxHeaderMegabytes int           `env:"MAX_HEADER_MEGABYTES" envDefault:"1"`
	}
	
	Postgres struct {
		User     string `env:"POSTGRES_USER" envDefault:"postgres"`
		Password string `env:"POSTGRES_PASSWORD,required"`
		Host     string `env:"POSTGRES_HOST" envDefault:"0.0.0.0"`
		Port     string `env:"POSTGRES_PORT" envDefault:"6543"`
		Database string `env:"POSTGRES_DATABASE" envDefault:"postgres"`
	}
}

func GetConfig() (*Config, error) {
	cfg := &Config{}
	
	log.Println("Gathering config...")
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	
	return cfg, nil
}
