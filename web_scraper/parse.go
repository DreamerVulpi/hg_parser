package web_scraper

import (
	"fmt"
	"log/slog"

	"github.com/gocolly/colly"
)

// TODO: ограничить количество выдаваемых результатов до нескольких на выбор пользователя
// TODO: при этом найти компромисс между производительностью и конечным результатом

func ParseProducts(collector *colly.Collector, search string) []Product {
	keyword := fmt.Sprintf("keyword=%s", search)
	link := base + catalog + "?" + keyword
	sliceProducts := make([]Product, 0)

	collector.OnHTML("div.product-item__content", func(e *colly.HTMLElement) {
		var element Product
		element.Price = e.ChildText("span.price")
		if element.Price != "" {
			element.Name = e.ChildAttr("div.name-desc a", "title")
			element.Img = e.ChildAttr("div.image a picture img", "src")
			element.Link = e.ChildAttr("div.image a", "href")
			sliceProducts = append(sliceProducts, element)
		}
	})

	collector.OnHTML("ul.pagination li", func(e *colly.HTMLElement) {
		rst := e.ChildAttr("a.next", "href")
		if rst != "" {
			link = base + catalog + rst + "&" + keyword
			err := collector.Visit(link)
			if err != nil {
				panic(err)
			}
		}
	})

	collector.OnRequest(func(r *colly.Request) {
		slog.Info("Посещение: " + link)
	})
	collector.OnResponse(func(r *colly.Response) {
		slog.Info("Страница посещена:" + link)
	})

	collector.OnScraped(func(r *colly.Response) {
		slog.Info("Запарсено!")
	})

	collector.Visit(link)
	return sliceProducts
}
