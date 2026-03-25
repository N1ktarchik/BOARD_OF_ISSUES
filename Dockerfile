FROM golang:1.25.1

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o /app/exe ./cmd/main.go

CMD ["/app/exe"]