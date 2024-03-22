package bot

import (
	"fmt"
	"hg_parser/db"
)

func (bot *Bot) commandFilter() error {
	bot.msg.ReplyMarkup = FilterMenu
	config, err := db.GetConfig(bot.conn, bot.update.Message.Chat.ID)
	if err != nil {
		bot.msg.Text = "Не удалось получить конфиг из базы данных."
		return nil
	}
	bot.msg.Text = fmt.Sprintf("Текущий фильтр:\n📉 Цена %s ₽\n🌚 Минимальный возраст: %s\n🎲 Количество игроков: %s\n⌛ Время игровой сессии (в минутах): %s\n🎛 Статус фильтра: %s", config["price"], config["age"], config["countplayers"], config["timesession"], config["switch"])
	bot.telegramBot.Send(bot.msg)
	return nil
}

func (bot *Bot) commandSearch() error {
	bot.msg.ReplyMarkup = SearchMenu
	bot.telegramBot.Send(bot.msg)
	return nil
}

func (bot *Bot) commandStart() error {
	bot.msg.Text = "Добро пожаловать!\nВы зарегистрированы в сервисе. Вам доступен фукнционал бота. Для ознакомления поможет команда /help\n"
	err := db.RegisterAccount(bot.conn, bot.update.Message.Chat.ID)
	if err != nil {
		return err
	}
	_, err = bot.telegramBot.Send(bot.msg)
	return err
}

func (bot *Bot) commandHelp() error {
	bot.msg.Text = `
		Список команд:
	/start 	- стартовая команда для запуска бота, обязательна, т.к при выполнении команды регистрирует аккаунт в сервисе
	/search - вызов клавиатуры поиска, где присутствует выбор между поиском и фильтра
	/filter - вызов клавиатуры настройки фильтра и ознакомление с текущим состоянием фильтра
	/help 	- список команд с пояснениями
	`
	_, err := bot.telegramBot.Send(bot.msg)
	return err
}
