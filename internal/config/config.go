package config

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

type Config struct {
	Mode                     string `mapstructure:"MODE"`
	PostgresDSN              string `mapstructure:"POSTGRES_DSN"`
	ProjectServiceRPCAddress string `mapstructure:"PROJECT_SERVICE_RPC_ADDR"`
	AccessSecret             string `mapstructure:"ACCESS_SECRET"`
}

func NewConfig() *Config {
	var cfg Config
	viper.AutomaticEnv()
	executablePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	executableDir := filepath.Dir(executablePath)
	viper.SetConfigFile(filepath.Join(executableDir, ".env"))
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}
	return &cfg
}
