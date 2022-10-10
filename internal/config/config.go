package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type AppConfig struct {
	GatewayPort uint32
	Host string
	Port uint32
	UserServiceUrl string
}

type Config struct {
	AppConfig 			AppConfig
}

const (
	appGatewayPort string = "app.grpc-gateway-port"
	appHost string = "app.host"
	appPort string = "app.port"
	userServiceUrl string = "services.user.url"
)


func GetConfig(env string) Config {
	n := fmt.Sprintf("config.%s", env)

	viper.SetConfigName(n)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Println(err)
	}

	return Config{
		AppConfig: AppConfig{
			GatewayPort: viper.GetUint32(appGatewayPort),
			Host: viper.GetString(appHost),
			Port: viper.GetUint32(appPort),
			UserServiceUrl: viper.GetString(userServiceUrl),
		},
	}
}