package utils

import (
	"log"

	env "github.com/caarlos0/env/v11"
)

type Env struct {
	Port string `env:"PORT" envDefault:"8080"`
}

func LoadEnv() Env {
	e := Env{}
	err := env.Parse(&e)
	if err != nil {
		log.Fatalf("failed to read env: %v", err)
	}

	return e
}
