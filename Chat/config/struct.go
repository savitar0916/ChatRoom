package config

import "github.com/sirupsen/logrus"

type Config struct {
	Server struct {
		Port int `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
	} `yaml:"server"`

	MySQL struct {
		User     string `yaml:"user" env:"MYSQL_USER" env-default:"root"`
		Password string `yaml:"password" env:"MYSQL_PASSWORD" env-default:"password"`
		Host     string `yaml:"host" env:"MYSQL_HOST" env-default:"127.0.0.1"`
		Port     int    `yaml:"port" env:"MYSQL_PORT" env-default:"3306"`
		Database string `yaml:"database" env:"MYSQL_DATABASE" env-default:"chat"`
	} `yaml:"mysql"`

	Logger struct {
		Level string `yaml:"level" env:"LOGGER_LEVEL" env-default:"info"`
	} `yaml:"logger"`
}

var AppConfig *Config
var Logger *logrus.Logger
