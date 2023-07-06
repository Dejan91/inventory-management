FROM golang:1.20.5-alpine

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

COPY . .
RUN go mod tidy

CMD ["air", "-c", ".air.toml"]