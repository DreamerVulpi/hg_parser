package main

import (
	"hg_parser/web_scraper"
	"log/slog"
)

func main() {
	slog.Info("Start program..")
	web_scraper.Parse("")
}
