package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Mode                          string
	WebPort                       int
	GcpProjectId                  string
	GoogleApplicationCredentials  string
	SubsServiceRPCAddress         string
	CheckmailTopicId              string
	CheckmailSubName              string
	CheckmailCreateAccessEndpoint string
}

func NewConfig() *Config {
	viper.AutomaticEnv()
	return &Config{
		Mode:                          viper.GetString("PRJCT_MODE"),
		WebPort:                       viper.GetInt("PORT"),
		GcpProjectId:                  viper.GetString("GCP_PROJECT_ID"),
		GoogleApplicationCredentials:  viper.GetString("GOOGLE_APPLICATION_CREDENTIALS"),
		SubsServiceRPCAddress:         viper.GetString("PRJCT_SUBS_SERVICE_RPC_ADDR"),
		CheckmailTopicId:              viper.GetString("PRJCT_CHECKMAIL_TOPIC_ID"),
		CheckmailSubName:              viper.GetString("PRJCT_CHECKMAIL_SUB_NAME"),
		CheckmailCreateAccessEndpoint: viper.GetString("PRJCT_CHECKMAIL_CREATE_ACCESS_ENDPOINT"),
	}
}
