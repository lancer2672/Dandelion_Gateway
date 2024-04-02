package utils

import "github.com/spf13/viper"

type Config struct {
	GatewayAddress             string `mapstructure:"GATEWAY_ADDRESS"`
	NotificationServiceAddress string `mapstructure:"NOTI_SERVICE_ADDR"`
	MainServiceAddress         string `mapstructure:"MAIN_SERVICE_ADDR"`
	MovieRESTAddress           string `mapstructure:"MOVIE_SERVICE_ADDR"`
	MovieGRPCAddress           string `mapstructure:"MOVIE_SERVICE_GRPC_ADDR"`
	GatewayApiKey              string `mapstructure:"GATEWAY_API_KEY"`
	RequestLimitTimeUnit       string `mapstructure:"REQUEST_LIMIT_TIMEUNIT"`
	RequestLimitPerTimeUnit    string `mapstructure:"REQUEST_LIMIT_PER_TIMEUNIT"`
	RedisURL                   string `mapstructure:"REDIS_URL"`
	RedisUsername              string `mapstructure:"REDIS_USERNAME"`
	RedisPassword              string `mapstructure:"REDIS_PASSWORD"`
}

// overrided by env if exists
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("dev")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
