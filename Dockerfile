FROM golang:latest
WORKDIR /app
COPY . .
RUN go get github.com/redis/go-redis/v9
RUN go get github.com/spf13/viper
RUN go get github.com/lib/pq
RUN go get github.com/gocolly/colly
CMD ["go", "run", "main.go"]