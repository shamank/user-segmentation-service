package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type (
	Config struct {
		Env            string `yaml:"env"`
		HTTPServer     `yaml:"http-server"`
		PostgresConfig `yaml:"postgres"`

		UserLogsDir string `yaml:"user-logs-dir"`
	}

	Addr struct {
	}

	HTTPServer struct {
		Host               string        `yaml:"host"`
		Port               string        `yaml:"port"`
		WriteTimeout       time.Duration `yaml:"writeTimeout" env-default:"10s"`
		ReadTimeout        time.Duration `yaml:"readTimeout" env-default:"10s"`
		MaxHeaderMegabytes int           `yaml:"maxHeaderMegabytes" env-default:"1"`
	}

	PostgresConfig struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `env-required:"true" env:"DB_PASSWORD"`
		DBName   string `yaml:"db"`
		SSL      string `yaml:"ssl"`
	}
)

func InitConfig(configPath string) *Config {

	if configPath == "" {
		log.Fatal("configPath is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("cannot find config file.. \"%s\"", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config file: \"%s\"", err.Error())
	}

	return &cfg
}
