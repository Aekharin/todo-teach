package config

import (
	"log"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

var config *Config
var Mode string

type Config struct {
	Server ServerConfig `mapstructure:"server" validate:"required"`
	DB     DBConfig     `mapstructure:"db" validate:"required"`
}

type ServerConfig struct {
	Name string `mapstructure:"name" validate:"required"`
	Port string `mapstructure:"port" validate:"required"`
}

type DBConfig struct {
	Host     string `mapstructure:"host" validate:"required"`
	Port     string `mapstructure:"port" validate:"required"`
	Username string `mapstructure:"username" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	Database string `mapstructure:"database" validate:"required"`
}

func InitialConfig() *Config {
	configPath, ok := os.LookupEnv("API_CONFIG_PATH")
	if !ok {
		configPath = "./config"
	}

	configName, ok := os.LookupEnv("API_CONFIG_NAME")
	if !ok {
		configName = "config"
	}

	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("cannot read in viper config: %s", err)
	}

	viper.AutomaticEnv()

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("failed to unmarshal config: %s", err)
	}

	if err := config.validate(); err != nil {
		log.Fatalf("config validation failed: %s", err)
	}

	return config
}

func (c *Config) validate() error {
	validate := validator.New()
	err := validate.Struct(c)
	if err != nil {
		return err
	}
	return nil
}
