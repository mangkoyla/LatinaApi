FROM golang:latest

WORKDIR /usr/src/app

COPY . .

RUN go mod download && go mod tidy && go mod verify
RUN go build -o ./latinaapi ./cmd/latinaapi/main.go

ENV GIN_MODE=release
EXPOSE 8080

CMD ["./latinaapi"]
