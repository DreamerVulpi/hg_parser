package bot

import (
	"hg_parser/web_scraper"
	"log/slog"

	telegramBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Init() (*telegramBot.BotAPI, error) {
	bot, err := telegramBot.NewBotAPI("TelegramApiKey")
	bot.Debug = true
	return bot, err
}

func Work(bot *telegramBot.BotAPI) {
	updateConfig := telegramBot.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)
	collector := web_scraper.Init()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := telegramBot.NewMessage(update.Message.Chat.ID, update.Message.Text)

		msg.ReplyToMessageID = update.Message.MessageID

		switch update.Message.Command() {
		case "search":
			result := web_scraper.ParseProducts(collector, update.Message.Text)
			for _, element := range result {
				text := element.Name + "\n" + element.Price + "\n" + element.Link
				slog.Info(text)
				telegramBot.NewInputMediaPhoto(telegramBot.FileURL(element.Img))
				msg.Text = text
				bot.Send(msg)
			}
		}

		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}
	}
}
