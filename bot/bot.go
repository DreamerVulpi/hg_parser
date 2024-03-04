package bot

import (
	"context"
	"database/sql"
	"log/slog"

	telegramBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/looplab/fsm"
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
	var ctx = context.Background()
	fsm := fsm.NewFSM(
		"start",
		fsm.Events{
			{Name: FilterMenu.Keyboard[0][0].Text, Src: []string{"start"}, Dst: "updatePrice"},
			{Name: "cancel", Src: []string{"updatePrice"}, Dst: "start"},

			{Name: FilterMenu.Keyboard[0][1].Text, Src: []string{"start"}, Dst: "updateAge"},
			{Name: "cancel", Src: []string{"updateAge"}, Dst: "start"},

			{Name: FilterMenu.Keyboard[0][2].Text, Src: []string{"start"}, Dst: "updateCountplayers"},
			{Name: "cancel", Src: []string{"updateCountplayers"}, Dst: "start"},

			{Name: FilterMenu.Keyboard[0][3].Text, Src: []string{"start"}, Dst: "updateTimesession"},
			{Name: "cancel", Src: []string{"updateTimesession"}, Dst: "start"},

			{Name: SearchMenu.Keyboard[0][0].Text, Src: []string{"start"}, Dst: "search"},
			{Name: "cancel", Src: []string{"search"}, Dst: "start"},
		},
		fsm.Callbacks{},
	)

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
		} else {
			err := handleKeyboardButton(bot, conn, &update, ctx, fsm)
			if err != nil {
				slog.Warn(err.Error())
				continue
			}
			continue
		}
	}
}
