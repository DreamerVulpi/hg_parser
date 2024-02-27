package bot

import (
	"database/sql"
	"hg_parser/db"
	"hg_parser/web_scraper"
	"log/slog"

	telegramBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gocolly/colly"
)

func handleCommand(bot *telegramBot.BotAPI, conn *sql.DB, update *telegramBot.Update) error {
	collector := web_scraper.Init()

	msg := telegramBot.NewMessage(update.Message.Chat.ID, update.Message.Text)

	// TODO: Add commands /help, /filter, /settings
	switch update.Message.Command() {
	case "search":
		return commandSearch(collector, update.Message.Text)
	case "start":
		return commandStart(&msg, bot)
	case "register":
		return commandRegistrationAccount(update, conn, &msg, bot)
	case "login":
		return commandLogInAccount(update, conn, &msg, bot)
	default:
		return unknownCommand(&msg, bot)
	}
}

func commandSearch(c *colly.Collector, keyword string) error {
	// TODO: Get list-filter from user account <- postgres
	// Example list-filter
	filter := map[string]string{
		"price":       "3000",
		"ageMin":      "12",
		"countPlayer": "2",
		"timeSession": "60",
	}
	// Debug parser
	web_scraper.WriteJSON(web_scraper.ParseProducts(c, filter, keyword))

	// Code for send messages from bot to user
	// result := web_scraper.ParseProducts(collector, update.Message.Text)
	// for _, element := range result {
	// 	text := element.Name + "\n" + element.Price + "\n" + element.Link + "\n" + element.AgePlayers + "\n" + element.CountPlayers + "\n" + element.TimeSession
	// 	slog.Info(text)
	// 	telegramBot.NewInputMediaPhoto(telegramBot.FileURL(element.Img))
	// 	msg.Text = text
	// 	bot.Send(msg)
	// }
	return nil
}

func commandStart(msg *telegramBot.MessageConfig, bot *telegramBot.BotAPI) error {
	msg.Text = "For using this bot you are need log in an account service.\n Don't know commands? Check list commands, just write /help"
	_, err := bot.Send(msg)
	return err
}

func commandRegistrationAccount(update *telegramBot.Update, conn *sql.DB, msg *telegramBot.MessageConfig, bot *telegramBot.BotAPI) error {
	if update.Message.Text != "" {
		err := db.RegisterAccount(conn, update.Message.From.UserName, update.Message.Text)
		if err != nil {
			slog.Warn(err.Error())
			msg.Text = "Account is registred already. Log in, please"
			bot.Send(msg)
			return err
		}
	}
	return nil
}

func commandLogInAccount(update *telegramBot.Update, conn *sql.DB, msg *telegramBot.MessageConfig, bot *telegramBot.BotAPI) error {
	_, err := db.GetUser(conn, update.Message.From.UserName, update.Message.Text)
	if err != nil {
		slog.Warn(err.Error())
		msg.Text = "Log in is failded. Try again"
		bot.Send(msg)
		return err
	}
	msg.Text = "Success!"
	bot.Send(msg)
	return nil
}

func unknownCommand(msg *telegramBot.MessageConfig, bot *telegramBot.BotAPI) error {
	msg.Text = "Unknown command"
	bot.Send(msg)
	return nil
}
