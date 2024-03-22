package bot

import (
	"context"
	"database/sql"
	"log/slog"

	telegramBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/looplab/fsm"
)

type Bot struct {
	telegramBot *telegramBot.BotAPI
	conn        *sql.DB
	update      *telegramBot.Update
	fsm         *fsm.FSM
	msg         telegramBot.MessageConfig
	ctx         context.Context
}

func (bot *Bot) Init(botsecretkey string, conn *sql.DB) (err error) {
	bot.telegramBot, err = telegramBot.NewBotAPI(botsecretkey)
	bot.conn = conn
	bot.telegramBot.Debug = true
	return err
}

func (bot *Bot) Start() {
	updateConfig := telegramBot.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.telegramBot.GetUpdatesChan(updateConfig)
	bot.handleUpdates(updates)
}

func (bot *Bot) handleUpdates(updates telegramBot.UpdatesChannel) {
	bot.ctx = context.Background()
	bot.fsm = fsm.NewFSM(
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

			{Name: FilterMenu.Keyboard[1][0].Text, Src: []string{"start"}, Dst: "updateStatusFilter"},
			{Name: "cancel", Src: []string{"updateStatusFilter"}, Dst: "start"},

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
			bot.update = &update
			err := bot.handleCommand()
			if err != nil {
				slog.Warn(err.Error())
				continue
			}
			continue
		} else {
			bot.update = &update
			err := bot.handleKeyboardButton()
			if err != nil {
				slog.Warn(err.Error())
				continue
			}
			continue
		}
	}
}
