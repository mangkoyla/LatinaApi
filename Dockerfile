FROM node:latest as web_builder

WORKDIR /usr/src/web

RUN git clone https://github.com/LalatinaHub/LatinaDocs .
RUN npm install
RUN npm run build

FROM golang:latest

WORKDIR /usr/src/app

COPY . .
COPY --from=web_builder /usr/src/web/docs/.vitepress/dist/ /usr/src/app/public/

RUN go mod download && go mod tidy && go mod verify
RUN go build -o ./latinaapi ./cmd/latinaapi/main.go

ENV GIN_MODE=release
EXPOSE 8080

CMD ["./latinaapi"]
