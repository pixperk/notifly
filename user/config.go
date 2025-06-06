package user

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port                int           `mapstructure:"PORT"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	NatsURL             string        `mapstructure:"NATS_URL"`
	NatsClusterID       string        `mapstructure:"NATS_CLUSTER_ID"`
	NatsClientID        string        `mapstructure:"NATS_CLIENT_ID"`
}

func LoadConfig(path string) (config Config, err error) {

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	if err = viper.Unmarshal(&config); err != nil {
		return
	}

	return
}
