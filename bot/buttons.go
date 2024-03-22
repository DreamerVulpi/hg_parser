package bot

import (
	"fmt"
	"hg_parser/db"
	"hg_parser/web_scraper"

	telegramBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var SearchMenu = telegramBot.NewReplyKeyboard(
	telegramBot.NewKeyboardButtonRow(
		telegramBot.NewKeyboardButton("🔎 Поиск товаров"),
		telegramBot.NewKeyboardButton("🎛 Фильтр"),
	),
)

var FilterMenu = telegramBot.NewReplyKeyboard(
	telegramBot.NewKeyboardButtonRow(
		telegramBot.NewKeyboardButton("📉 Максимальная цена"),
		telegramBot.NewKeyboardButton("🌚 Минимальный возраст игрока"),
		telegramBot.NewKeyboardButton("🎲 Количество игроков"),
		telegramBot.NewKeyboardButton("⌛ Время игровой сессии"),
	),
	telegramBot.NewKeyboardButtonRow(
		telegramBot.NewKeyboardButton("🎛 Статус фильтра"),
	),
)

func (bot *Bot) UpdateTimesession() error {
	switch bot.fsm.Current() {
	case "start":
		bot.msg = telegramBot.NewMessage(bot.update.Message.Chat.ID, "Укажите количество своего свободного времени (в минутах)")
		_, err := bot.telegramBot.Send(bot.msg)
		if err != nil {
			return err
		}
		err = bot.fsm.Event(bot.ctx, FilterMenu.Keyboard[0][3].Text)
		if err != nil {
			return err
		}
	case "updateTimesession":
		bot.msg = telegramBot.NewMessage(bot.update.Message.Chat.ID, "?")
		_, err := bot.telegramBot.Send(bot.msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bot *Bot) UpdateCountPlayers() error {
	switch bot.fsm.Current() {
	case "start":
		bot.msg = telegramBot.NewMessage(bot.update.Message.Chat.ID, "Укажите количество игроков (от 1 до 8 человек)")
		_, err := bot.telegramBot.Send(bot.msg)
		if err != nil {
			return err
		}
		err = bot.fsm.Event(bot.ctx, FilterMenu.Keyboard[0][2].Text)
		if err != nil {
			return err
		}
	case "updateCountplayers":
		bot.msg = telegramBot.NewMessage(bot.update.Message.Chat.ID, "?")
		_, err := bot.telegramBot.Send(bot.msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bot *Bot) UpdateAge() error {
	switch bot.fsm.Current() {
	case "start":
		bot.msg = telegramBot.NewMessage(bot.update.Message.Chat.ID, "Укажите минимальный возраст игрока. (0 - любой, до 18)")
		_, err := bot.telegramBot.Send(bot.msg)
		if err != nil {
			return err
		}
		err = bot.fsm.Event(bot.ctx, FilterMenu.Keyboard[0][1].Text)
		if err != nil {
			return err
		}
	case "updateAge":
		msg := telegramBot.NewMessage(bot.update.Message.Chat.ID, "?")
		_, err := bot.telegramBot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bot *Bot) UpdatePrice() error {
	switch bot.fsm.Current() {
	case "start":
		bot.msg = telegramBot.NewMessage(bot.update.Message.Chat.ID, "Укажите минимальную цену покупки (в рублях, например: 1000)")
		_, err := bot.telegramBot.Send(bot.msg)
		if err != nil {
			return err
		}
		err = bot.fsm.Event(bot.ctx, FilterMenu.Keyboard[0][0].Text)
		if err != nil {
			return err
		}
	case "updatePrice":
		bot.msg = telegramBot.NewMessage(bot.update.Message.Chat.ID, "?")
		_, err := bot.telegramBot.Send(bot.msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bot *Bot) UpdateStatusFilter() error {
	switch bot.fsm.Current() {
	case "start":
		bot.msg = telegramBot.NewMessage(bot.update.Message.Chat.ID, "Укажите статус фильтра (Активировать - 1 или выключить - 0)")
		_, err := bot.telegramBot.Send(bot.msg)
		if err != nil {
			return err
		}
		err = bot.fsm.Event(bot.ctx, FilterMenu.Keyboard[1][0].Text)
		if err != nil {
			return err
		}
	case "updateStatusFilter":
		bot.msg = telegramBot.NewMessage(bot.update.Message.Chat.ID, "?")
		_, err := bot.telegramBot.Send(bot.msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bot *Bot) SearchProducts() error {
	switch bot.fsm.Current() {
	case "start":
		bot.msg = telegramBot.NewMessage(bot.update.Message.Chat.ID, "Введите название искомого товара")
		_, err := bot.telegramBot.Send(bot.msg)
		if err != nil {
			return err
		}
		err = bot.fsm.Event(bot.ctx, SearchMenu.Keyboard[0][0].Text)
		if err != nil {
			return err
		}
	case "search":
		bot.msg = telegramBot.NewMessage(bot.update.Message.Chat.ID, "?")
		_, err := bot.telegramBot.Send(bot.msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bot *Bot) Filter() error {
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

func (bot *Bot) eventsButtons() error {
	switch bot.fsm.Current() {
	case "updatePrice":
		err := db.UpdatePrice(bot.conn, bot.update.Message.Chat.ID, bot.update.Message.Text)
		if err != nil {
			return err
		}
		bot.msg.Text = "Цена в фильтре успешно изменена."
		_, err = bot.telegramBot.Send(bot.msg)
		if err != nil {
			return nil
		}
		err = bot.fsm.Event(bot.ctx, "cancel")
		if err != nil {
			return err
		}
	case "updateAge":
		err := db.UpdateAge(bot.conn, bot.update.Message.Chat.ID, bot.update.Message.Text)
		if err != nil {
			return err
		}
		bot.msg.Text = "Минимальный возраст игрока в фильтре успешно изменен."
		_, err = bot.telegramBot.Send(bot.msg)
		if err != nil {
			return nil
		}
		err = bot.fsm.Event(bot.ctx, "cancel")
		if err != nil {
			return err
		}
	case "updateCountplayers":
		err := db.UpdateCountPlayers(bot.conn, bot.update.Message.Chat.ID, bot.update.Message.Text)
		if err != nil {
			return err
		}
		bot.msg.Text = "Количество игроков в фильтре успешно изменено."
		_, err = bot.telegramBot.Send(bot.msg)
		if err != nil {
			return nil
		}
		err = bot.fsm.Event(bot.ctx, "cancel")
		if err != nil {
			return err
		}
	case "updateTimesession":
		err := db.UpdateTimeSession(bot.conn, bot.update.Message.Chat.ID, bot.update.Message.Text)
		if err != nil {
			return err
		}
		bot.msg.Text = "Свободное время изменено в фильтре успешно."
		_, err = bot.telegramBot.Send(bot.msg)
		if err != nil {
			return nil
		}
		err = bot.fsm.Event(bot.ctx, "cancel")
		if err != nil {
			return err
		}
	case "updateStatusFilter":
		err := db.UpdateStateFilter(bot.conn, bot.update.Message.Chat.ID, bot.update.Message.Text)
		if err != nil {
			return err
		}
		bot.msg.Text = "Статус фильтра изменен успешно."
		_, err = bot.telegramBot.Send(bot.msg)
		if err != nil {
			return err
		}
		err = bot.fsm.Event(bot.ctx, "cancel")
		if err != nil {
			return err
		}
	case "search":
		filter, err := db.GetConfig(bot.conn, bot.update.Message.Chat.ID)
		if err != nil {
			return err
		}
		collector := web_scraper.Init()
		bot.msg.Text = "поиск товаров..."
		bot.telegramBot.Send(bot.msg)
		result := web_scraper.ParseProducts(collector, filter, bot.update.Message.Text)
		for _, element := range result {
			text := element.Name + "\n" + "📉 " + element.Price + "\n" + element.Link + "\n"
			if element.AgePlayers != "" {
				text += "🌚 " + element.AgePlayers + "\n"
			}
			if element.TimeSession != "" {
				text += "⌛ " + element.TimeSession + "\n"
			}
			if element.CountPlayers != "" {
				text += "🎲 " + element.CountPlayers + "\n"
			}
			telegramBot.NewInputMediaPhoto(telegramBot.FileURL(element.Img))
			bot.msg.Text = text
			bot.telegramBot.Send(bot.msg)
		}
		err = bot.fsm.Event(bot.ctx, "cancel")
		if err != nil {
			return err
		}
	}
	return nil
}
