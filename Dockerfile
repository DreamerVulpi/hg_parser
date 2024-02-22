FROM golang:latest
WORKDIR /app
COPY . .
RUN go get github.com/redis/go-redis/v9
RUN go get github.com/spf13/viper
RUN go get github.com/lib/pq
RUN go get github.com/gocolly/colly
RUN go get github.com/go-telegram-bot-api/telegram-bot-api/v5
CMD ["go", "run", "main.go"]