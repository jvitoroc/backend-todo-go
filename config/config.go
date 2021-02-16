package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	CFG_PROD = "./config/prod.yaml"
	CFG_TEST = "./config/test.yaml"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	}
	DB struct {
		SQLiteDSN string `yaml:"sqliteDSN"`
	}
	Email struct {
		SmtpUser string `yaml:"smtpUser"`
		SmtpPass string `yaml:"smtpPass"`
		SmtpHost string `yaml:"smtpHost"`
		SmtpAddr string `yaml:"smtpAddr"`
	}
	Auth struct {
		JwtSecret      string `yaml:"jwtSecret"`
		GoogleClientID string `yaml:"googleClientID"`
	}
}

func NewConfig(path string) *Config {
	cfg := &Config{}
	if err := cfg.load(path); err != nil {
		log.Fatalf("Could not load configuration: %s", err.Error())
	}

	return cfg
}

func (cfg *Config) load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(cfg); err != nil {
		return err
	}

	return nil
}
