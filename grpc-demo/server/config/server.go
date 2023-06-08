package config

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Name                string `mapstructure:"name"`
	GrpcServerAddr      string `mapstructure:"grpc_server_addr"`
	LocalGrpcServerAddr string `mapstructure:"local_grpc_server_addr"`
	HttpServerAddr      string `mapstructure:"http_server_addr"`
	Logger              string `mapstructure:"logger"`
}

var (
	serverConfig   = ServerConfig{}
	configFilename string
)

func init() {
	flag.StringVar(&configFilename, "config", "./deploy/dev/config.yaml", "config file")
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
