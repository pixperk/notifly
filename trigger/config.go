package trigger

import (
	"github.com/spf13/viper"
)

type Config struct {
	NatsURL       string `mapstructure:"NATS_URL"`
	NatsClusterID string `mapstructure:"NATS_CLUSTER_ID"`
	NatsClientID  string `mapstructure:"NATS_CLIENT_ID"`
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
