package config

import (
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	SnowflakeNode int64  `env:"SNOWFLAKE_NODE" envDefault:"1"`
	LogLevel      string `env:"LOG_LEVEL" envDefault:"info"`

	JWTAuthentication struct {
		SigningKey       string        `env:"SIGNING_KEY,required"`
		Issuer           string        `env:"ISSUER" envDefault:"https://localhost:8080/auth"`
		ExpectedAudience string        `env:"EXPECTED_AUDIENCE" envDefault:"https://localhost:8080"`
		TargetAudience   []string      `env:"TARGET_AUDIENCE" envDefault:"https://localhost:8080,https://cource-watch.com"`
		TokenTTL         time.Duration `env:"TOKEN_TTL" envDefault:"1h"`
	}

	HTTP struct {
		Host               string        `env:"HOST" envDefault:"localhost"`
		Port               string        `env:"PORT" envDefault:"8080"`
		ReadTimeout        time.Duration `env:"READ_TIMEOUT" envDefault:"5s"`
		WriteTimeout       time.Duration `env:"WRITE_TIMEOUT" envDefault:"5s"`
		MaxHeaderMegabytes int           `env:"MAX_HEADER_MEGABYTES" envDefault:"1"`
	}

	Postgres struct {
		User     string `env:"POSTGRES_USER" envDefault:"postgres"`
		Password string `env:"POSTGRES_PASSWORD,required"`
		Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
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
