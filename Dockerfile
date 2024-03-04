FROM golang:latest
WORKDIR /hg_parser/app
COPY . .

RUN go get github.com/spf13/viper
RUN go get github.com/lib/pq
RUN go get github.com/gocolly/colly
RUN go get github.com/go-telegram-bot-api/telegram-bot-api/v5
RUN go get github.com/looplab/fsm
RUN go get github.com/golang-migrate/migrate/v4
CMD ["go", "run", "main.go"]
