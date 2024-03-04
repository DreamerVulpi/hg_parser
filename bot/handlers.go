package bot

import (
	"context"
	"database/sql"

	telegramBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/looplab/fsm"
)

func handleCommand(bot *telegramBot.BotAPI, conn *sql.DB, update *telegramBot.Update) error {
	msg := telegramBot.NewMessage(update.Message.Chat.ID, update.Message.Text)

	switch update.Message.Command() {
	case "search":
		return commandSearch(bot, &msg)
	case "filter":
		return commandFilter(bot, conn, &msg, update.Message.Chat.ID)
	case "start":
		return commandStart(&msg, conn, update.Message.Chat.ID, bot)
	case "help":
		return commandHelp(&msg, bot)
	default:
		return unknown(&msg, bot)
	}
}

func handleKeyboardButton(bot *telegramBot.BotAPI, conn *sql.DB, update *telegramBot.Update, ctx context.Context, fsm *fsm.FSM) error {
	msg := telegramBot.NewMessage(update.Message.Chat.ID, update.Message.Text)

	switch update.Message.Text {
	case SearchMenu.Keyboard[0][0].Text:
		return buttonSearchProducts(fsm, ctx, bot, update)
	case SearchMenu.Keyboard[0][1].Text:
		return buttonFilter(bot, conn, update.Message.Chat.ID, &msg)

	case FilterMenu.Keyboard[0][0].Text:
		return buttonUpdatePrice(fsm, ctx, bot, update)
	case FilterMenu.Keyboard[0][1].Text:
		return buttonUpdateAge(fsm, ctx, bot, update)
	case FilterMenu.Keyboard[0][2].Text:
		return buttonUpdateCountPlayers(fsm, ctx, bot, update)
	case FilterMenu.Keyboard[0][3].Text:
		return buttonUpdateTimesession(fsm, ctx, bot, update)

	default:
		if err := eventsButtons(fsm, conn, bot, update, ctx, &msg); err != nil {
			return err
		}
		return nil
	}
}

func unknown(msg *telegramBot.MessageConfig, bot *telegramBot.BotAPI) error {
	msg.Text = "Unknown data"
	bot.Send(msg)
	return nil
}
