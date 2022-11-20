package config

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	GatewayPort uint32
	Host string
	Port uint32
	UserServiceUrl string
}

const (
	appGatewayPort string = "app.grpc-gateway-port"
	appHost string = "app.host"
	appPort string = "app.port"
	userServiceUrl string = "services.user.url"
)


func GetAppConfig() AppConfig {
return AppConfig{
		GatewayPort: viper.GetUint32(appGatewayPort),
		Host: viper.GetString(appHost),
		Port: viper.GetUint32(appPort),
		UserServiceUrl: viper.GetString(userServiceUrl),
	}
};