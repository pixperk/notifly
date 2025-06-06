package graphql

import "github.com/spf13/viper"

type Config struct {
	UserURL    string `mapstructure:"USER_URL"`
	TriggerURL string `mapstructure:"TRIGGER_URL"`
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
