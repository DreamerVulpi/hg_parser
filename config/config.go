package config

import (
	"log/slog"

	"github.com/spf13/viper"
)

type TelegramBot struct {
	SecretKey string `mapstructure:"secretKey"`
}

type Postgres struct {
	User     string
	Password string
	DBname   string
	Sslmode  string
}

type Redis struct {
	DBname   string
	Password string
	Port     string
}

type AppConfig struct {
	Path     string
	NameFile string
	TypeFile string
	Bot      TelegramBot
	Postgres Postgres
	Redis    Redis
}

func LoadConfig(config *AppConfig) error {
	v := viper.New()
	v.SetConfigFile(config.NameFile)
	v.SetConfigType(config.TypeFile)
	v.AddConfigPath(config.Path)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		slog.Warn(err.Error())
		return err
	} else {
		config.Bot.SecretKey = v.GetString("telegrambot.secretkey")

		config.Postgres.User = v.GetString("postgres.user")
		config.Postgres.Password = v.GetString("postgres.password")
		config.Postgres.DBname = v.GetString("postgres.dbname")
		config.Postgres.Sslmode = v.GetString("postgres.sslmode")

		config.Redis.DBname = v.GetString("redis.dbname")
		config.Redis.Password = v.GetString("redis.password")
		config.Redis.Port = v.GetString("redis.port")

		return err
	}
}
