package config

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Name                string
	GrpcServerAddr      string
	LocalGrpcServerAddr string
	HttpServerAddr      string
	Logger              string
}

var (
	serverConfig   = ServerConfig{}
	configFilename string
)

func init() {
	flag.StringVar(&configFilename, "config", "deploy/config.yml", "config file")
}

func GetServerConfig() *ServerConfig {
	return &serverConfig
}

func SetConfigFilename(filename string) {
	configFilename = filename
}

func Init() {
	flag.Parse()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFilename)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error read config %s, err:%s", configFilename, err))
	}
	err = viper.Unmarshal(&serverConfig)
	if err != nil {
		panic(fmt.Errorf("error viper.Unnarshal %v", err))
	}

	fmt.Printf("serverConfig:%+v\n", serverConfig)
}
