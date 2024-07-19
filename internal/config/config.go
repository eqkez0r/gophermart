package config

import (
	"errors"
	"flag"
	e "github.com/eqkez0r/gophermart/pkg/error"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

const (
	defaultRunAddr           = "127.0.0.1:8880"
	defaultAccrualSystemAddr = "http://127.0.0.1:8080"
)

var (
	errEmptyDatabaseURI = errors.New("empty database uri")
)

func NewConfig() (*Config, error) {
	const op = "Initial config error: "

	cfg := &Config{}
	flag.StringVar(&cfg.RunAddress, "a", defaultRunAddr, "run address")
	flag.StringVar(&cfg.DatabaseURI, "d", "", "database uri")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", defaultAccrualSystemAddr, "")
	flag.Parse()

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, e.Wrap(op, err)
	}
	if cfg.DatabaseURI == "" {
		return nil, e.Wrap(op, errEmptyDatabaseURI)
	}

	return cfg, nil
}
