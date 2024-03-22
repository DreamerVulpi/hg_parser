package main

import (
	"fmt"
	telegramBot "hg_parser/bot"
	"hg_parser/config"
	"hg_parser/db"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	slog.Info("Start program..")

	var conf config.AppConfig
	conf.Path = "."
	conf.NameFile = "config"
	conf.TypeFile = "yaml"
	config.Load(&conf)

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", conf.Postgres.User, conf.Postgres.Password, conf.Postgres.Host, conf.Postgres.Port, conf.Postgres.DBname, conf.Postgres.Sslmode)
	m, err := migrate.New("file://db/migrations", url)
	slog.Info(url)
	if err != nil {
		slog.Warn(err.Error())
	}
	if err := m.Up(); err != nil {
		slog.Warn(err.Error())
	}

	params := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s port=%s", conf.Postgres.User, conf.Postgres.Password, conf.Postgres.DBname, conf.Postgres.Sslmode, conf.Postgres.Host, conf.Postgres.Port)
	conn, err := db.ConnectToPostgres("postgres", params)
	if err != nil {
		slog.Warn(err.Error())
	}
	var bot telegramBot.Bot
	err = bot.Init(conf.BotSecretKey, conn)
	if err != nil {
		slog.Warn(err.Error())
	}
	bot.Start()
}
