package web_scraper

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/gocolly/colly"
)

type Product struct {
	Name  string `json:"name"`
	Img   string `json:"img"`
	Price string `json:"price"`
	Link  string `json:"link"`
}

func Parse(search string) {
	c := colly.NewCollector(
		// TODO: params
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"),
	)

	listProducts := make([]Product, 0)

	c.OnRequest(func(r *colly.Request) {
		slog.Info("Visiting: ", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		slog.Warn(err.Error())
	})

	c.OnResponse(func(r *colly.Response) {
		slog.Info("Page visited: ", r.Request.URL)
	})

	c.OnHTML("div.product-item__content", func(e *colly.HTMLElement) {
		var element Product
		element.Link = e.ChildAttr("div.image a", "href")
		element.Price = e.ChildText("span.price")
		element.Name = e.ChildAttr("div.name-desc a", "title")
		element.Img = e.ChildAttr("div.image a picture img", "src")
		listProducts = append(listProducts, element)
	})

	c.OnScraped(func(r *colly.Response) {
		slog.Info(r.Request.URL.String(), " scraped!")
	})

	c.Visit("https://hobbygames.ru/nastolnie")
	writeJSON(listProducts)
}

func writeJSON(data []Product) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		slog.Info("Unable to create JSON file")
		return
	}
	os.WriteFile("productsHG.json", file, 0644)
}
