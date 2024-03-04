package db

import (
	"database/sql"
	"fmt"
	"log/slog"
)

func GetUserID(conn *sql.DB, chat_id int64) int {
	chat_id_str := fmt.Sprintf("%d", chat_id)
	var cht_id int
	_ = conn.QueryRow("SELECT id FROM users WHERE chat_id = $1;", chat_id_str).Scan(&cht_id)
	return cht_id
}

func GetUser(conn *sql.DB, chat_id int64) bool {
	chat_id_str := fmt.Sprintf("%d", chat_id)
	var result bool
	conn.QueryRow("SELECT EXISTS(SELECT * FROM users WHERE chat_id = $1);", chat_id_str).Scan(&result)
	return result
}

func RegisterAccount(conn *sql.DB, chat_id int64) error {
	chat_id_str := fmt.Sprintf("%d", chat_id)
	if !GetUser(conn, chat_id) {
		var user_id int
		err := conn.QueryRow("INSERT INTO users (chat_id) VALUES ($1) RETURNING id;", chat_id_str).Scan(&user_id)
		if err != nil {
			return err
		}

		var list_id int
		err = conn.QueryRow("INSERT INTO lists_config (price, age, countplayers, timesession) VALUES ($1, $2, $3, $4) RETURNING id;", 0, 0, 0, 0).Scan(&list_id)
		if err != nil {
			return err
		}

		_, err = conn.Exec("INSERT INTO lists_users (user_id, list_id) VALUES ($1, $2);", user_id, list_id)
		if err != nil {
			slog.Warn(err.Error())
			return err
		}
	}
	return nil
}
