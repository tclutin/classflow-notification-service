FROM golang:1.23.3-alpine3.20

WORKDIR /app

COPY . .

RUN go mod download && go mod tidy

RUN go build -o ./ cmd/main.go

EXPOSE 8091

CMD ["./main"]