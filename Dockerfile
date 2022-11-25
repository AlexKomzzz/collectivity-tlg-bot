FROM golang:latest AS build

# RUN go version
WORKDIR /github.com/AlexKomzzz/collectivity-tlg-bot/

COPY go.* ./

RUN go mod download

COPY . /github.com/AlexKomzzz/collectivity-tlg-bot/

RUN CGO_ENABLED=0 go build -o ./.bin/bot ./cmd/bot/main.go

FROM scratch

WORKDIR /

COPY --from=build /github.com/AlexKomzzz/collectivity-tlg-bot/.bin/bot .
COPY --from=build /github.com/AlexKomzzz/collectivity-tlg-bot/configs configs/
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/  

EXPOSE 9090

CMD ["./bot"]