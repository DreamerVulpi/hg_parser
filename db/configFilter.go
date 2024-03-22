package db

import (
	"database/sql"
	"log/slog"
)

func UpdatePrice(conn *sql.DB, chat_id int64, price string) error {
	_, err := conn.Exec("UPDATE lists_config lc SET price=$1 FROM lists_users lu WHERE lc.id = lu.list_id AND lu.user_id = $2", price, GetUserID(conn, chat_id))
	return err
}
func UpdateAge(conn *sql.DB, chat_id int64, age string) error {
	_, err := conn.Exec("UPDATE lists_config lc SET age=$1 FROM lists_users lu WHERE lc.id = lu.list_id AND lu.user_id = $2", age, GetUserID(conn, chat_id))
	return err
}
func UpdateCountPlayers(conn *sql.DB, chat_id int64, countplayers string) error {
	_, err := conn.Exec("UPDATE lists_config lc SET countplayers=$1 FROM lists_users lu WHERE lc.id = lu.list_id AND lu.user_id = $2", countplayers, GetUserID(conn, chat_id))
	return err
}
func UpdateTimeSession(conn *sql.DB, chat_id int64, timesession string) error {
	_, err := conn.Exec("UPDATE lists_config lc SET timesession=$1 FROM lists_users lu WHERE lc.id = lu.list_id AND lu.user_id = $2", timesession, GetUserID(conn, chat_id))
	return err
}

func UpdateStateFilter(conn *sql.DB, chat_id int64, status string) error {
	_, err := conn.Exec("UPDATE lists_config lc SET switch=$1 FROM lists_users lu WHERE lc.id = lu.list_id AND lu.user_id = $2", status, GetUserID(conn, chat_id))
	return err
}

func GetConfig(conn *sql.DB, chat_id int64) (map[string]string, error) {
	row, err := conn.Query("SELECT lc.id, lc.price, lc.age, lc.countplayers, lc.timesession, lc.switch FROM lists_config lc INNER JOIN lists_users lu ON lc.id = lu.list_id WHERE lu.user_id = $1;", GetUserID(conn, chat_id))
	if err != nil {
		slog.Warn(err.Error())
		return map[string]string{}, err
	}
	defer row.Close()

	result := map[string]string{}
	if row.Next() {
		var id, price, age, countplayers, timesession, status string
		err := row.Scan(&id, &price, &age, &countplayers, &timesession, &status)
		if err != nil {
			return map[string]string{}, err
		}
		result["id"] = id
		result["price"] = price
		result["age"] = age
		result["countplayers"] = countplayers
		result["timesession"] = timesession

		if status == "1" {
			result["switch"] = "Активен"
		} else {
			result["switch"] = "Выключен"
		}
	}
	return result, err
}
