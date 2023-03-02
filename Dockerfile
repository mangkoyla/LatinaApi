FROM node:lts as web_builder

WORKDIR /usr/src/web

RUN git clone https://github.com/LalatinaHub/LatinaDocs .
RUN npm install
RUN npm run build

FROM golang:latest as web_app

WORKDIR /usr/src/app

COPY . .
COPY --from=web_builder /usr/src/web/docs/.vitepress/dist/ /usr/src/app/public/

ENV GIN_MODE=release
ENV API_MODE=true
EXPOSE 8080

CMD ["bash", "start.sh"]
