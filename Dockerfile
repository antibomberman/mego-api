FROM golang:1.22.5
WORKDIR /app
COPY . .

RUN go mod tidy

RUN go build -o api cmd/api/main.go
CMD ["./api"]