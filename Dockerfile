FROM node:lts as web_builder

WORKDIR /usr/src/web

RUN git clone https://github.com/LalatinaHub/LatinaDocs .
RUN npm install
RUN npm run build

FROM golang:latest as web_app

WORKDIR /usr/src/app

COPY . .
COPY --from=web_builder /usr/src/web/docs/.vitepress/dist/ /usr/src/app/public/

# RUN apt-get update -y
# RUN curl -fsSL https://deb.nodesource.com/setup_18.x -o ./node_source.sh 
# RUN bash ./node_source.sh && apt-get install nodejs build-essential -y
# RUN git clone https://github.com/LalatinaHub/LatinaSub
# RUN cd LatinaSub && npm install
# RUN mkdir LatinaSub/bin
# RUN curl "https://github.com/tindy2013/subconverter/releases/download/v0.7.2/subconverter_linux64.tar.gz" -L -o ./LatinaSub/bin/subconverter.tar.gz
# RUN go install -v -tags with_shadowsocksr,with_grpc github.com/sagernet/sing-box/cmd/sing-box@latest
# RUN cp $GOPATH/bin/sing-box ./LatinaSub/bin/
# RUN tar -xf ./LatinaSub/bin/subconverter.tar.gz -C ./LatinaSub/bin
# RUN rm -rf ./LatinaSub/bin/subconverter.tar.gz
# RUN chmod +x ./LatinaSub/bin/*

RUN go mod download && go mod tidy && go mod verify
RUN go build -o ./latinaapi ./cmd/latinaapi/main.go

ENV GIN_MODE=release
ENV API_MODE=true
EXPOSE 8080

CMD ["./latinaapi"]
