package config

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config contains app configuration variables
type Config struct {
    Env string `yaml:"env" envconfig:"ENV" default:"local"`
	Server struct {
		Port int    `yaml:"port" envconfig:"SERVER_PORT" default:"8080"`
		Addr string `yaml:"host" envconfig:"SERVER_ADDR" default:"0.0.0.0"`
	    ReadTimeout   time.Duration    `yaml:"read_timeout" envconfig:"SERVER_READ_TIMEOUT" default:"10s"` 
    	WriteTimeout  time.Duration    `yaml:"write_timeout" envconfig:"SERVER_WRITE_TIMEOUT" default:"10s"`
        IdleTimeout  time.Duration    `yaml:"idle_timeout" envconfig:"SERVER_IDLE_TIMEOUT" default:"60s"`
	}
	Database struct {
        User     string `yaml:"user" envconfig:"DB_USER" default:"postgres"`
        Password string `yaml:"password" envconfig:"DB_PASSWORD" default:"postgres"`
		Host     string `yaml:"host" envconfig:"DB_HOST" default:"localhost"`
		Port     string `yaml:"port" envconfig:"DB_PORT" default:"5432"`
		Name     string `yaml:"name" envconfig:"DB_NAME" default:"authdb"`
	}
    // PrivateKey and PublicKey is base64 encoded ecdsa256 keys
    PrivateKey string `yaml:"private_key" envconfig:"PRIVATE_KEY"`
    PublicKey string `yaml:"public_key" envconfig:"PUBLIC_KEY"`
    GoogleClientID string `yaml:"google_client_id" envconfig:"GOOGLE_CLIENT_ID"`
    GoogleClientSecret string `yaml:"google_client_secret" envconfig:"GOOGLE_CLIENT_SECRET"`
    TokenTTL time.Duration `yaml:"token_ttl" envconfig:"TOKEN_TTL" default:"15m"`
}

// MustLoadConfig load configs from env variables
func MustLoadConfig() *Config {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
        log.Fatal("Env variables error: ", err)
	}
	return &config
}

