package db

import (
	"crypto/sha256"
	"database/sql"
)

func HashPassword(password string) ([]byte, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(password))
	if err != nil {
		return []byte{}, err
	}
	rslt := hash.Sum(nil)
	return rslt, nil
}

func GetUser(conn *sql.DB, username, password string) (int, error) {
	password_hash, err := HashPassword(password)
	if err != nil {
		return 0, err
	}
	row, err := conn.Query("SELECT id FROM users WHERE username=$1 AND password_hash=$2", username, password_hash)
	result := 0
	if err != nil {
		return 0, err
	} else {
		row.Scan(&result)
		return result, err
	}
}

func RegisterAccount(conn *sql.DB, login, password string) error {
	password_hash, err := HashPassword(password)
	if err != nil {
		return err
	}
	_, err = conn.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", login, password_hash)
	// TODO:  Генерация дефолтной индивидуальной конфигурации для поиска товаров при регистрации
	// ("INSERT INTO lists_configuration (keyword, ageMin, countPlayer, timeSession) VALUES ($1, $2,)")
	return err
}
