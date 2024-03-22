package bot

import (
	telegramBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bot *Bot) handleCommand() error {
	bot.msg = telegramBot.NewMessage(bot.update.Message.Chat.ID, bot.update.Message.Text)

	switch bot.update.Message.Command() {
	case "search":
		return bot.commandSearch()
	case "filter":
		return bot.commandFilter()
	case "start":
		return bot.commandStart()
	case "help":
		return bot.commandHelp()
	default:
		return bot.unknown()
	}
}

func (bot *Bot) handleKeyboardButton() error {
	bot.msg = telegramBot.NewMessage(bot.update.Message.Chat.ID, bot.update.Message.Text)

	switch bot.update.Message.Text {
	case SearchMenu.Keyboard[0][0].Text:
		return bot.SearchProducts()
	case SearchMenu.Keyboard[0][1].Text:
		return bot.Filter()

	case FilterMenu.Keyboard[0][0].Text:
		return bot.UpdatePrice()
	case FilterMenu.Keyboard[0][1].Text:
		return bot.UpdateAge()
	case FilterMenu.Keyboard[0][2].Text:
		return bot.UpdateCountPlayers()
	case FilterMenu.Keyboard[0][3].Text:
		return bot.UpdateTimesession()

	case FilterMenu.Keyboard[1][0].Text:
		return bot.UpdateStatusFilter()

	default:
		if err := bot.eventsButtons(); err != nil {
			return err
		}
		return nil
	}
}

func (bot *Bot) unknown() error {
	bot.msg.Text = "Unknown data"
	bot.telegramBot.Send(bot.msg)
	return nil
}
