package bot

import (
	"database/sql"
	"log/slog"

	telegramBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Init(botsecretkey string) (*telegramBot.BotAPI, error) {
	bot, err := telegramBot.NewBotAPI(botsecretkey)
	bot.Debug = true
	return bot, err
}

func Start(conn *sql.DB, bot *telegramBot.BotAPI) {
	updateConfig := telegramBot.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	handleUpdates(bot, conn, updates)
}

func handleUpdates(bot *telegramBot.BotAPI, conn *sql.DB, updates telegramBot.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			err := handleCommand(bot, conn, &update)
			if err != nil {
				slog.Warn(err.Error())
				continue
			}
			continue
		}
	}
}
