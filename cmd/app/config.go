package main

import (
	"github.com/spf13/viper"
)

const (
	defaultMode  = "prod"
	defaultPort  = "8080"
	defaultProto = "http"
)

type Config struct {
	Mode                         string
	Host                         string
	Port                         string
	Proto                        string
	GcpProjectId                 string
	GoogleApplicationCredentials string
	SubscriptionServiceGRPCAddr  string
}

func NewConfig() *Config {
	viper.AutomaticEnv()
	viper.SetDefault("MODE", defaultMode)
	viper.SetDefault("PORT", defaultPort)
	viper.SetDefault("PROTO", defaultProto)

	return &Config{
		Mode:                         viper.GetString("MODE"),
		Host:                         viper.GetString("HOST"),
		Port:                         viper.GetString("PORT"),
		Proto:                        viper.GetString("PROTO"),
		GcpProjectId:                 viper.GetString("GCP_PROJECT_ID"),
		GoogleApplicationCredentials: viper.GetString("GOOGLE_APPLICATION_CREDENTIALS"),
		SubscriptionServiceGRPCAddr:  viper.GetString("SBS_SERVICE_GRPC_ADDR"),
	}
}
