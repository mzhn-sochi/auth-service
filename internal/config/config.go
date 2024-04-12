package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	App struct {
		Host string `env:"APP_HOST" env-default:"0.0.0.0"`
		Port int    `env:"APP_PORT" env-default:"8080"`
	}

	DB struct {
		User string `env:"DB_USER" env-default:"postgres"`
		Pass string `env:"DB_PASS" env-default:"postgres"`
		Host string `env:"DB_HOST" env-default:"localhost"`
		Port int    `env:"DB_PORT" env-default:"5436"`
		Name string `env:"DB_NAME" env-default:"users"`
	}

	JWT struct {
		Access struct {
			TTL    int    `env:"JWT_ACCESS_TTL" env-default:"15"`
			Secret string `env:"JWT_ACCESS_SECRET"`
		}

		Refresh struct {
			TTL    int    `env:"JWT_REFRESH_TTL" env-default:"604800"`
			Secret string `env:"JWT_REFRESH_SECRET"`
		}
	}
}

func New() *Config {
	config := &Config{}

	if err := cleanenv.ReadEnv(config); err != nil {
		header := "AUTH SERVICE ENVs"
		f := cleanenv.FUsage(os.Stdout, config, &header)
		f()
		panic(err)
	}

	log.Println(config)

	return config
}
