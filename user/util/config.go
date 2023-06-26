package util

import "github.com/spf13/viper"

type Config struct {
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBName            string `mapstructure:"DB_NAME"`
	GrpcServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`
	FirebaseProjectID string `mapstructure:"FIREBASE_PROJECT_ID"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
