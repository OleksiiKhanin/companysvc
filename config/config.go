package config

import "time"

type Config struct {
	Server   ServerConfig   `json:"server"`
	Loc      LocatorConfig  `yaml:"loc"`
	Db       DatabaseConfig `yaml:"db"`
	Event    QueueConfig    `yaml:"event"`
	LogLevel string         `yaml:"logLevel"`
}

type LocatorConfig struct {
	URL                  string   `yaml:"url"`
	RetryAttempt         int      `yaml:"retryAttempt"`
	AllowedCountiesCodes []string `yaml:"allowedCountiesCodes"`
}

type DatabaseConfig struct {
	URL        string `yaml:"url"`
	Port       int    `yaml:"port"`
	Login      string `yaml:"login"`
	Password   string `yaml:"password"`
	NameDB     string `yaml:"nameDB"`
	MaxConns   int    `yaml:"maxConns"`
	Migrations string `yaml:"migrations"`
}

type QueueConfig struct {
	URL           string        `yaml:"url"`
	Port          int           `yaml:"port"`
	EventChannel  string        `yaml:"eventChannel"`
	ReconnectWait time.Duration `yaml:"reconnectWait"`
	PingInterval  time.Duration `yaml:"pingInterval"`
}

type ServerConfig struct {
	URL       string `yaml:"url"`
	PrefixAPI string `yaml:"prefixAPI"`
}
