package main

import (
	"fmt"
	telegramBot "hg_parser/bot"
	"hg_parser/config"
	"hg_parser/db"
	"log/slog"
)

func main() {
	slog.Info("Start program..")

	var conf config.AppConfig
	conf.Path = "."
	conf.NameFile = "config"
	conf.TypeFile = "yaml"
	config.Load(&conf)

	params := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", conf.Postgres.User, conf.Postgres.Password, conf.Postgres.DBname, conf.Postgres.Sslmode)
	conn, err := db.ConnectToPostgres("postgres", params)
	if err != nil {
		slog.Warn(err.Error())
	}

	bot, err := telegramBot.Init(conf.BotSecretKey)
	if err != nil {
		slog.Warn(err.Error())
	}
	telegramBot.Start(conn, bot)
}
