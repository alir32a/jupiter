package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

var cfg *Config

type Config struct {
	DB               *DBConfig
	Manager          *ManagerConfig
	HTTPServerConfig *HTTPServerConfig
	MainBot          *MainBotConfig
	OCCTL            *OCCTLConfig
	TrialPackage     *TrialPackageConfig
}

type DBConfig struct {
	User          string `envconfig:"DB_USER"`
	Password      string `envconfig:"DB_PASSWORD"`
	DBName        string `envconfig:"DB_NAME"`
	Host          string `envconfig:"DB_HOST"`
	Port          int    `envconfig:"DB_PORT" default:"5432"`
	SSLMode       string `envconfig:"DB_SSL_MODE" default:"disable"`
	MigrationPath string `envconfig:"DB_MIGRATION_PATH" default:"database/migrations"`
}

type ManagerConfig struct {
	UpdateInterval       time.Duration `envconfig:"MANAGER_UPDATE_INTERVAL" default:"5s"`
	UpdateTimeout        time.Duration `envconfig:"MANAGER_UPDATE_INTERVAL" default:"60s"`
	MaxFailures          int           `envconfig:"MANAGER_MAX_FAILURES" default:"10"`
	ShutdownOnMaxFailure bool          `envconfig:"MANAGER_SHUTDOWN_MAX_FAILURES" default:"true"`
}

type HTTPServerConfig struct {
	Host                  string        `envconfig:"SERVER_HOST" default:"127.0.0.1"`
	Port                  int           `envconfig:"SERVER_PORT" default:"8080"`
	AccessTokenSecret     string        `envconfig:"SERVER_ACCESS_TOKEN_SECRET" required:"true"`
	AccessTokenExpireTime time.Duration `envconfig:"SERVER_ACCESS_TOKEN_EXPIRE_TIME" required:"true"`
	ENV                   string        `envconfig:"SERVER_ENV" required:"true"`
}

type MainBotConfig struct {
	Token    string `envconfig:"MAIN_BOT_TOKEN"`
	OCCTLCfg OCCTLConfig
}

type OCCTLConfig struct {
	PasswordFile string `envconfig:"OCCTL_PASSWORD_FILE"`
}

type TrialPackageConfig struct {
	Activated        bool    `envconfig:"TRIAL_PACKAGE_ACTIVATED" default:"false"`
	TrafficLimit     float64 `envconfig:"TRIAL_PACKAGE_TRAFFIC_LIMIT" default:"5"`
	MaxConnections   int     `envconfig:"TRIAL_PACKAGE_MAX_CONNECTIONS" default:"2"`
	ExpirationInDays int     `envconfig:"TRIAL_PACKAGE_EXPIRATION" default:"7"`
}

func GetConfig() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{}
	if err := envconfig.Process("jupiter", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
