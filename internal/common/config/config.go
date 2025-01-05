package config

import (
	"github.com/spf13/viper"
)

const (
	defaultMode  = "prod"
	defaultPort  = 8080
	defaultProto = "http"
)

type Config struct {
	Mode                          string
	Port                          int
	Proto                         string
	GcpProjectId                  string
	GoogleApplicationCredentials  string
	SubsServiceRPCAddress         string
	CheckmailTopicId              string
	CheckmailSubName              string
	CheckmailCreateAccessEndpoint string
}

func NewConfig() *Config {
	viper.AutomaticEnv()
	mode := viper.GetString("MODE")
	if mode == "" {
		mode = defaultMode
	}
	port := viper.GetInt("PORT")
	if port == 0 {
		port = defaultPort
	}
	proto := viper.GetString("PROTO")
	if proto == "" {
		proto = defaultProto
	}
	return &Config{
		Mode:                          mode,
		Port:                          port,
		Proto:                         proto,
		GcpProjectId:                  viper.GetString("GCP_PROJECT_ID"),
		GoogleApplicationCredentials:  viper.GetString("GOOGLE_APPLICATION_CREDENTIALS"),
		SubsServiceRPCAddress:         viper.GetString("PRJCT_SUBS_SERVICE_RPC_ADDR"),
		CheckmailTopicId:              viper.GetString("PRJCT_CHECKMAIL_TOPIC_ID"),
		CheckmailSubName:              viper.GetString("PRJCT_CHECKMAIL_SUB_NAME"),
		CheckmailCreateAccessEndpoint: viper.GetString("PRJCT_CHECKMAIL_CREATE_ACCESS_ENDPOINT"),
	}
}
