package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Judge   JudgeConfig   `mapstructure:"judge"`
	Queue   QueueConfig   `mapstructure:"queue"`
	JWT     JWTConfig     `mapstructure:"jwt"`
	Log     LogConfig     `mapstructure:"log"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type JudgeConfig struct {
	WorkerCount  int    `mapstructure:"worker_count"`
	TimeLimit    int    `mapstructure:"time_limit"`
	MemoryLimit  int    `mapstructure:"memory_limit"`
	OutputLimit  int    `mapstructure:"output_limit"`
	GoPath       string `mapstructure:"go_path"`
	GccPath      string `mapstructure:"gcc_path"`
	GppPath      string `mapstructure:"gpp_path"`
	JavacPath    string `mapstructure:"javac_path"`
	JavaPath     string `mapstructure:"java_path"`
}

type QueueConfig struct {
	MaxSize       int `mapstructure:"max_size"`
	WorkerTimeout int `mapstructure:"worker_timeout"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

var AppConfig *Config

func Load(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode)
}
