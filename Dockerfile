FROM golang:1.20.4-alpine3.18

# RUN apk add --no-cache docker-cli

WORKDIR /app

COPY . .

RUN go build -o docknest-server

ENTRYPOINT ["/app/docknest-server"]
