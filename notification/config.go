package notification

import (
	"github.com/spf13/viper"
)

type Config struct {
	NatsURL           string `mapstructure:"NATS_URL"`
	NatsClusterID     string `mapstructure:"NATS_CLUSTER_ID"`
	NatsClientID      string `mapstructure:"NATS_CLIENT_ID"`
	TwilioAuthToken   string `mapstructure:"TWILIO_AUTH_TOKEN"`
	TwilioAccountSID  string `mapstructure:"TWILIO_ACCOUNT_SID"`
	TwilioPhoneNumber string `mapstructure:"TWILIO_PHONE_NUMBER"`
	BrevoAPIKey       string `mapstructure:"BREVO_API_KEY"`
	FromAddress       string `mapstructure:"FROM_ADDRESS"`
	Port              int    `mapstructure:"PORT"`
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
