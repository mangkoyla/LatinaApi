FROM node:lts as web_builder

FROM golang:latest as web_app

WORKDIR /usr/src/app

COPY . .
COPY --from=web_builder /usr/src/web/docs/.vitepress/dist/ /usr/src/app/public/

# Drop replace
RUN go mod edit -dropreplace="github.com/mangkoyla/LatinaBot"
RUN go mod edit -dropreplace="github.com/mangkoyla/LatinaSub-go"

RUN go get -v github.com/mangkoyla/LatinaBot@main
RUN go get -v github.com/mangkoyla/LatinaSub-go@main
RUN go mod download && go mod tidy && go mod verify
RUN go build -tags with_grpc -o ./latinaapi ./cmd/latinaapi/main.go

ENV GIN_MODE=release
ENV API_MODE=true
EXPOSE 8080

CMD ["./latinaapi"]
