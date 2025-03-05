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

    GRPC struct {
        Server struct {
            Host     string        `yaml:"host" envconfig:"GRPC_SERVER_HOST" default:"0.0.0.0"`
            Port     int           `yaml:"port" envconfig:"GRPC_SERVER_PORT" default:"9090"`
            UseTLS   bool          `yaml:"use_tls" envconfig:"GRPC_SERVER_USE_TLS" default:"true"`
            CertFile string        `yaml:"cert_file" envconfig:"GRPC_SERVER_CERT_FILE" default:"certs/server.crt"`
            KeyFile  string        `yaml:"key_file" envconfig:"GRPC_SERVER_KEY_FILE" default:"certs/server.key"`
            CAFile   string        `yaml:"ca_file" envconfig:"GRPC_SERVER_CA_FILE" default:"certs/ca.crt"`
            Timeout  time.Duration `yaml:"timeout" envconfig:"GRPC_SERVER_TIMEOUT" default:"5s"`
        }
        Client struct {
            Host         string        `yaml:"host" envconfig:"GRPC_CLIENT_HOST" default:"localhost"`
            Port         int           `yaml:"port" envconfig:"GRPC_CLIENT_PORT" default:"9090"`
            UseTLS       bool          `yaml:"use_tls" envconfig:"GRPC_CLIENT_USE_TLS" default:"true"`
            CertFile     string        `yaml:"cert_file" envconfig:"GRPC_CLIENT_CERT_FILE" default:"certs/client.crt"`
            KeyFile      string        `yaml:"key_file" envconfig:"GRPC_CLIENT_KEY_FILE" default:"certs/client.key"`
            CAFile       string        `yaml:"ca_file" envconfig:"GRPC_CLIENT_CA_FILE" default:"certs/ca.crt"`
            Timeout      time.Duration `yaml:"timeout" envconfig:"GRPC_CLIENT_TIMEOUT" default:"5s"`
            RetryCount   int           `yaml:"retry_count" envconfig:"GRPC_CLIENT_RETRY_COUNT" default:"3"`
            RetryTimeout time.Duration `yaml:"retry_timeout" envconfig:"GRPC_CLIENT_RETRY_TIMEOUT" default:"2s"`
        }
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

    // Microservices
    AuthServiceHost         string `yaml:"auth_service_host" envconfig:"AUTH_SERVICE_HOST" default:"localhost"`
	AuthServicePort         int    `yaml:"auth_service_port" envconfig:"AUTH_SERVICE_PORT" default:"8081"`
	EventServiceHost        string `yaml:"event_service_host" envconfig:"EVENT_SERVICE_HOST" default:"localhost"`
	EventServicePort        int    `yaml:"event_service_port" envconfig:"EVENT_SERVICE_PORT" default:"8082"`
	RegistrationServiceHost string `yaml:"registration_service_host" envconfig:"REGISTRATION_SERVICE_HOST" default:"localhost"`
	RegistrationServicePort int    `yaml:"registration_service_port" envconfig:"REGISTRATION_SERVICE_PORT" default:"8083"`
	UserServiceHost         string `yaml:"user_service_host" envconfig:"USER_SERVICE_HOST" default:"localhost"`
	UserServicePort         int    `yaml:"user_service_port" envconfig:"USER_SERVICE_PORT" default:"8084"`
	ReviewsServiceHost      string `yaml:"reviews_service_host" envconfig:"REVIEWS_SERVICE_HOST" default:"localhost"`
	ReviewsServicePort      int    `yaml:"reviews_service_port" envconfig:"REVIEWS_SERVICE_PORT" default:"8085"`
	MediaServiceHost        string `yaml:"media_service_host" envconfig:"MEDIA_SERVICE_HOST" default:"localhost"`
	MediaServicePort        int    `yaml:"media_service_port" envconfig:"MEDIA_SERVICE_PORT" default:"8086"`
	NotificationServiceHost string `yaml:"notification_service_host" envconfig:"NOTIFICATION_SERVICE_HOST" default:"localhost"`
	NotificationServicePort int    `yaml:"notification_service_port" envconfig:"NOTIFICATION_SERVICE_PORT" default:"8087"`

	// Конфигурация для Kafka
	Kafka struct {
		Brokers []string `yaml:"brokers" envconfig:"KAFKA_BROKERS" default:"localhost:9092"`
		GroupID string   `yaml:"group_id" envconfig:"KAFKA_GROUP_ID" default:"default_group"`
		Topic   string   `yaml:"topic" envconfig:"KAFKA_TOPIC" default:"default_topic"`
	} `yaml:"kafka"`
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

