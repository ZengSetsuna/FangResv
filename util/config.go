package util

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER" default:"postgres"`
	DBSource             string        `mapstructure:"DB_SOURCE" default:"postgresql://shu:shu@localhost:5432/simple_bank?sslmode=disable"`
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS" default:":8080"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS" default:":8081"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION" default:"15m"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION" default:"72h"`
	SMTPHost             string        `mapstructure:"SMTP_HOST"`
	SMTPPort             int           `mapstructure:"SMTP_PORT"`
	SMTPUsername         string        `mapstructure:"SMTP_USERNAME"`
	SMTPPassword         string        `mapstructure:"SMTP_PASSWORD"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	fmt.Println("config: ", config)
	if err != nil {
		return
	}
	return
}
