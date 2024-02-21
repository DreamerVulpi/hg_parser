package db

import (
	"context"
	"database/sql"
	"log/slog"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func ConnectToPostgres(driver, params string) (*sql.DB, error) {
	conn, err := sql.Open(driver, params)
	if err != nil {
		return nil, err
	}
	return conn, err
}

func ConnectToRedis(params map[string]string) (*redis.Client, error) {
	dbname, err := strconv.Atoi(params["dbname"])
	conn := redis.NewClient(&redis.Options{
		Addr:     "redis:" + params["port"],
		Password: params["password"],
		DB:       dbname,
	})
	ctx := context.Background()
	pong := conn.Ping(ctx)
	slog.Info(pong.Result())
	return conn, err
}
