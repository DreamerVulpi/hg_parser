package web_scraper

import (
	"log/slog"

	"github.com/gocolly/colly"
)

type Product struct {
	Name         string `json:"name"`
	Img          string `json:"img"`
	Price        string `json:"price"`
	CountPlayers string `json:"countPlayers"`
	TimeSession  string `json:"timeSession"`
	AgePlayers   string `json:"agePlayers"`
	Link         string `json:"link"`
}

func Init() *colly.Collector {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"),
	)

	c.OnError(func(_ *colly.Response, err error) {
		slog.Warn(err.Error())
	})
	return c
}
