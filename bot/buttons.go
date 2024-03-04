package bot

import (
	"context"
	"database/sql"
	"fmt"
	"hg_parser/db"
	"hg_parser/web_scraper"

	telegramBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/looplab/fsm"
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
)

func buttonUpdateTimesession(fsm *fsm.FSM, ctx context.Context, bot *telegramBot.BotAPI, update *telegramBot.Update) error {
	switch fsm.Current() {
	case "start":
		msg := telegramBot.NewMessage(update.Message.Chat.ID, "Укажите количество своего свободного времени (в минутах)")
		_, err := bot.Send(msg)
		if err != nil {
			return err
		}
		err = fsm.Event(ctx, FilterMenu.Keyboard[0][3].Text)
		if err != nil {
			return err
		}
	case "updateTimesession":
		msg := telegramBot.NewMessage(update.Message.Chat.ID, "?")
		_, err := bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func buttonUpdateCountPlayers(fsm *fsm.FSM, ctx context.Context, bot *telegramBot.BotAPI, update *telegramBot.Update) error {
	switch fsm.Current() {
	case "start":
		msg := telegramBot.NewMessage(update.Message.Chat.ID, "Укажите количество игроков (от 1 до 8 человек)")
		_, err := bot.Send(msg)
		if err != nil {
			return err
		}
		err = fsm.Event(ctx, FilterMenu.Keyboard[0][2].Text)
		if err != nil {
			return err
		}
	case "updateCountplayers":
		msg := telegramBot.NewMessage(update.Message.Chat.ID, "?")
		_, err := bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func buttonUpdateAge(fsm *fsm.FSM, ctx context.Context, bot *telegramBot.BotAPI, update *telegramBot.Update) error {
	switch fsm.Current() {
	case "start":
		msg := telegramBot.NewMessage(update.Message.Chat.ID, "Укажите минимальный возраст игрока. (0 - любой, до 18)")
		_, err := bot.Send(msg)
		if err != nil {
			return err
		}
		err = fsm.Event(ctx, FilterMenu.Keyboard[0][1].Text)
		if err != nil {
			return err
		}
	case "updateAge":
		msg := telegramBot.NewMessage(update.Message.Chat.ID, "?")
		_, err := bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func buttonUpdatePrice(fsm *fsm.FSM, ctx context.Context, bot *telegramBot.BotAPI, update *telegramBot.Update) error {
	switch fsm.Current() {
	case "start":
		msg := telegramBot.NewMessage(update.Message.Chat.ID, "Укажите минимальную цену покупки (в рублях, например: 1000)")
		_, err := bot.Send(msg)
		if err != nil {
			return err
		}
		err = fsm.Event(ctx, FilterMenu.Keyboard[0][0].Text)
		if err != nil {
			return err
		}
	case "updatePrice":
		msg := telegramBot.NewMessage(update.Message.Chat.ID, "?")
		_, err := bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func buttonSearchProducts(fsm *fsm.FSM, ctx context.Context, bot *telegramBot.BotAPI, update *telegramBot.Update) error {
	switch fsm.Current() {
	case "start":
		msg := telegramBot.NewMessage(update.Message.Chat.ID, "Введите название искомого товара")
		_, err := bot.Send(msg)
		if err != nil {
			return err
		}
		err = fsm.Event(ctx, SearchMenu.Keyboard[0][0].Text)
		if err != nil {
			return err
		}
	case "search":
		msg := telegramBot.NewMessage(update.Message.Chat.ID, "?")
		_, err := bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func buttonFilter(bot *telegramBot.BotAPI, conn *sql.DB, chat_id int64, msg *telegramBot.MessageConfig) error {
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

func eventsButtons(fsm *fsm.FSM, conn *sql.DB, bot *telegramBot.BotAPI, update *telegramBot.Update, ctx context.Context, msg *telegramBot.MessageConfig) error {
	switch fsm.Current() {
	case "updatePrice":
		err := db.UpdatePrice(conn, update.Message.Chat.ID, update.Message.Text)
		if err != nil {
			return err
		}
		msg.Text = "Цена в фильтре успешно изменена."
		_, err = bot.Send(msg)
		if err != nil {
			return nil
		}
		err = fsm.Event(ctx, "cancel")
		if err != nil {
			return err
		}
	case "updateAge":
		err := db.UpdateAge(conn, update.Message.Chat.ID, update.Message.Text)
		if err != nil {
			return err
		}
		msg.Text = "Минимальный возраст игрока в фильтре успешно изменен."
		_, err = bot.Send(msg)
		if err != nil {
			return nil
		}
		err = fsm.Event(ctx, "cancel")
		if err != nil {
			return err
		}
	case "updateCountplayers":
		err := db.UpdateCountPlayers(conn, update.Message.Chat.ID, update.Message.Text)
		if err != nil {
			return err
		}
		msg.Text = "Количество игроков в фильтре успешно изменено."
		_, err = bot.Send(msg)
		if err != nil {
			return nil
		}
		err = fsm.Event(ctx, "cancel")
		if err != nil {
			return err
		}
	case "updateTimesession":
		err := db.UpdateTimeSession(conn, update.Message.Chat.ID, update.Message.Text)
		if err != nil {
			return err
		}
		msg.Text = "Свободное время изменено в фильтре успешно."
		_, err = bot.Send(msg)
		if err != nil {
			return nil
		}
		err = fsm.Event(ctx, "cancel")
		if err != nil {
			return err
		}
	case "search":
		filter, err := db.GetConfig(conn, update.Message.Chat.ID)
		if err != nil {
			return err
		}
		collector := web_scraper.Init()
		web_scraper.WriteJSON(web_scraper.ParseProducts(collector, filter, update.Message.Text))
		msg.Text = "Найдены товары."
		result := web_scraper.ParseProducts(collector, filter, update.Message.Text)
		for _, element := range result {
			text := element.Name + "\n" + element.Price + "\n" + element.Link + "\n" + element.AgePlayers + "\n" + element.CountPlayers + "\n" + element.TimeSession
			telegramBot.NewInputMediaPhoto(telegramBot.FileURL(element.Img))
			msg.Text = text
			bot.Send(msg)
		}
		_, err = bot.Send(msg)
		if err != nil {
			return err
		}
		err = fsm.Event(ctx, "cancel")
		if err != nil {
			return err
		}
	}
	return nil
}
