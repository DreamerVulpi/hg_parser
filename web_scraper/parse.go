package web_scraper

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

const (
	base       = "https://hobbygames.ru"
	catalog    = "/catalog/search"
	boardgames = "/nastolnie"
)

func WriteJSON(data []Product) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		slog.Info("Unable to create JSON file")
		return
	}
	os.WriteFile("productsHG.json", file, 0644)
}

func getProduct(e *colly.HTMLElement, filter map[string]string) (Product, error) {
	var element Product
	element.Price = e.ChildText("span.price")
	if isAvailiable(element.Price) {
		element.Name = e.ChildAttr("div.name-desc a", "title")
		element.Img = e.ChildAttr("div.image a picture img", "src")
		element.Link = e.ChildAttr("div.image a", "href")
		element.CountPlayers = e.ChildText("div.params__item.players span")
		element.TimeSession = e.ChildText("div.params__item.time span")
		element.AgePlayers = e.ChildText("div.age__number")
		if biggerPrice(element.Price, filter["price"]) && agePlayers(element.AgePlayers, filter["ageMin"]) && countPlayers(element.CountPlayers, filter["countPlayer"]) && timeSession(element.TimeSession, filter["timeSession"]) {
			return element, nil
		}
	}
	return Product{}, fmt.Errorf("Product out of stock")
}

func biggerPrice(priceStr, filterPrice string) bool {
	if !strings.Contains(priceStr, "бон.") {
		filter, _ := strconv.Atoi(filterPrice)
		temp := strings.Replace(priceStr, "₽", "", -1)
		price, _ := strconv.Atoi(strings.Replace(temp, " ", "", -1))
		slog.Info(strings.Replace(temp, " ", "", -1))
		return filter <= price
	}
	return false
}

func isAvailiable(price string) bool {
	return price != ""
}

func agePlayers(ageStr, filterAge string) bool {
	filter, _ := strconv.Atoi(filterAge)
	age, _ := strconv.Atoi(strings.Replace(ageStr, "+", "", -1))
	return filter <= age // TODO: Child filter
}

func countPlayers(countPlayersStr, filterCountPlayers string) bool {
	filter, _ := strconv.Atoi(filterCountPlayers)
	if strings.Contains(countPlayersStr, "-") {
		sliceCountPlayers := strings.Split(countPlayersStr, "-")
		countPlayersMin, _ := strconv.Atoi(sliceCountPlayers[0])
		countPlayersMax, _ := strconv.Atoi(sliceCountPlayers[1])
		return countPlayersMin <= filter && filter <= countPlayersMax || filter == 0
	}
	countPlayers, _ := strconv.Atoi(strings.Replace(countPlayersStr, "+", "", -1))
	return filter <= countPlayers || countPlayers == 0
}

func timeSession(timeSessionStr, filterTimeSession string) bool {
	filter, _ := strconv.Atoi(filterTimeSession)
	if strings.Contains(timeSessionStr, "-") {
		sliceTimeSession := strings.Split(timeSessionStr, "-")
		timeSessionMin, _ := strconv.Atoi(sliceTimeSession[0])
		timeSessionMax, _ := strconv.Atoi(sliceTimeSession[1])
		return timeSessionMax >= filter && filter >= timeSessionMin || filter == 0
	}
	timeSession, _ := strconv.Atoi(strings.Replace(timeSessionStr, "+", "", -1))
	return filter <= timeSession || filter == 0
}

func ParseProducts(collector *colly.Collector, filter map[string]string, search string) []Product {
	keyword := fmt.Sprintf("keyword=%s", search)
	link := base + catalog + "?" + keyword
	sliceProducts := make([]Product, 0)
	sliceTest := make([]Product, 0)

	collector.OnHTML("div.product-item__content", func(e *colly.HTMLElement) {
		product, err := getProduct(e, filter)
		if err != nil {
			return
		}
		sliceTest = append(sliceTest, product)
		//sliceProducts = append(sliceProducts, element)
		// sliceTest = append(sliceTest, Filter(&element, map[string]string{"AgeMin": "12"}))

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
