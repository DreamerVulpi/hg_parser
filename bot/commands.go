package bot

import (
	"database/sql"
	"fmt"
	"hg_parser/db"

	telegramBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func commandFilter(bot *telegramBot.BotAPI, conn *sql.DB, msg *telegramBot.MessageConfig, chat_id int64) error {
	msg.ReplyMarkup = FilterMenu
	config, err := db.GetConfig(conn, chat_id)
	if err != nil {
		msg.Text = "Не удалось получить конфиг из базы данных."
		return nil
	}
	msg.Text = fmt.Sprintf("Текущий фильтр:\n📉 Цена %s ₽\n🌚 Минимальный возраст: %s\n🎲 Количество игроков: %s\n⌛ Время игровой сессии (в минутах): %s", config["price"], config["age"], config["countplayers"], config["timesession"])
	bot.Send(msg)
	return nil
}

func commandSearch(bot *telegramBot.BotAPI, msg *telegramBot.MessageConfig) error {
	msg.ReplyMarkup = SearchMenu
	bot.Send(msg)
	return nil
}

func commandStart(msg *telegramBot.MessageConfig, conn *sql.DB, chat_id int64, bot *telegramBot.BotAPI) error {
	msg.Text = "Добро пожаловать!\nВы зарегистрированы в сервисе. Вам доступен фукнционал бота. Для ознакомления поможет команда /help\n"
	err := db.RegisterAccount(conn, chat_id)
	if err != nil {
		return err
	}
	_, err = bot.Send(msg)
	return err
}

func commandHelp(msg *telegramBot.MessageConfig, bot *telegramBot.BotAPI) error {
	msg.Text = `
		Список команд:
	/start 	- стартовая команда для запуска бота, обязательна, т.к при выполнении команды регистрирует аккаунт в сервисе
	/search - вызов клавиатуры поиска, где присутствует выбор между поиском и фильтра
	/filter - вызов клавиатуры настройки фильтра и ознакомление с текущим состоянием фильтра
	/help 	- список команд с пояснениями
	`
	_, err := bot.Send(msg)
	return err
}
