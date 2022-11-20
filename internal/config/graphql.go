package config

import (
	"github.com/spf13/viper"
)

const (
	graphQLSchemeKey	string = "services.graphql.scheme"
	graphQLHostKey 		string = "services.graphql.host"
	graphQLPathKey		string = "services.graphql.path"
)

type GraphQLConfig struct {
	Scheme	string
	Host    string
	Path 		string
}

func GetGraphQLConfig() GraphQLConfig {
	return GraphQLConfig{
		Scheme: 		viper.GetString(graphQLSchemeKey),		
		Host: 		viper.GetString(graphQLHostKey),
		Path: 		viper.GetString(graphQLPathKey),
	}
}