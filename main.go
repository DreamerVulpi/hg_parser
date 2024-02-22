package main

import (
	telegramBot "hg_parser/bot"
	"log/slog"
)

func main() {
	slog.Info("Start program..")
	bot, err := telegramBot.Init()
	if err != nil {
		slog.Warn(err.Error())
	}
	telegramBot.Work(bot)
}
