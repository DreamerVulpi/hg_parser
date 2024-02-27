package config

import (
	"log/slog"

	"github.com/spf13/viper"
)

type Postgres struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBname   string `mapstructure:"dbname"`
	Sslmode  string `mapstructure:"sslmode"`
}

type AppConfig struct {
	Path         string
	NameFile     string
	TypeFile     string
	BotSecretKey string `mapstructure:"secretkey"`
	Postgres     Postgres
}

func Load(config *AppConfig) error {
	v := viper.New()
	v.SetConfigName(config.NameFile)
	v.SetConfigType(config.TypeFile)
	v.AddConfigPath(config.Path)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		slog.Warn(err.Error())
		return err
	} else {
		config.BotSecretKey = v.GetString("telegrambot.secretkey")
		config.Postgres.User = v.GetString("postgres.user")
		config.Postgres.Password = v.GetString("postgres.password")
		config.Postgres.DBname = v.GetString("postgres.dbname")
		config.Postgres.Sslmode = v.GetString("postgres.sslmode")
		return err
	}
}
