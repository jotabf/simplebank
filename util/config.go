package util

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	HTTPServerAddr       string        `mapstructure:"HTTP_SERVER_ADDR"`
	GRPCServerAddr       string        `mapstructure:"GRPC_SERVER_ADDR"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

// LoadConfig loads the configuration from the given path. It uses Viper to read
// the configuration file and Unmarshal it into the given config struct.
// It returns an error if the configuration file could not be read or
// unmarshalled.
//
// Parameters:
//
//	path: The path of the configuration file.
//
// Returns:
//
//	config: The configuration struct.
//	err:    An error if the configuration file could not be read or
//	        unmarshalled.
func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
