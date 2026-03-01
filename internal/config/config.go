package config

import "time"

type Config struct {
	Env    string       `yaml:"env" env:"ENV"`
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"db"`
	Kafka  KafkaConfig  `yaml:"kafka"`
	URL    URLConfig    `yaml:"url"`
	Logger LogConfig    `yaml:"logger"`
}

type ServerConfig struct {
	Address string `yaml:"address" env:"SERVER_ADDRESS"`
}

type DBConfig struct {
	URL             string        `yaml:"url" env:"DB_URL"`
	Name            string        `yaml:"name" env:"DB_NAME"`
	Password        string        `yaml:"password" env:"DB_PASSWORD"`
	MaxOpensConns   int           `yaml:"max_opens_conns" env:"DB_MAX_OPENS_CONNS"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env:"DB_MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env:"DB_CONN_MAX_LIFETIME"`
	ConnMaxIdle     time.Duration `yaml:"conn_max_idle" env:"DB_CONN_MAX_IDLE"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers" env:"KAFKA_BROKERS"`
	Topic   string   `yaml:"topic" env:"KAFKA_TOPIC"`
}

type URLConfig struct {
	Base         string        `yaml:"base" env:"URL_BASE"`
	CodeLength   int           `yaml:"code_length" env:"URL_CODE_LENGTH"`
	CodeAlphabet string        `yaml:"code_alphabet" env:"URL_CODE_ALPHABET"`
	DefaultTTL   time.Duration `yaml:"default_ttl" env:"URL_DEFAULT_TTL"`
}

type LogConfig struct {
	Level        string `yaml:"level" env:"LOG_LEVEL"`
	ShowPathCall bool   `yaml:"show_path_call" env:"LOG_SHOW_PATH_CALL"`
}
