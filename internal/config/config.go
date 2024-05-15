package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Mode                         string
	WebPort                      int
	GcpProjectId                 string
	GoogleApplicationCredentials string
	SubsServiceRPCAddress        string
}

func NewConfig() *Config {
	viper.AutomaticEnv()
	return &Config{
		Mode:                         viper.GetString("PRJCT_MODE"),
		WebPort:                      viper.GetInt("PORT"),
		GcpProjectId:                 viper.GetString("GCP_PROJECT_ID"),
		GoogleApplicationCredentials: viper.GetString("GOOGLE_APPLICATION_CREDENTIALS"),
		SubsServiceRPCAddress:        viper.GetString("PRJCT_SUBS_SERVICE_RPC_ADDR"),
	}
}
