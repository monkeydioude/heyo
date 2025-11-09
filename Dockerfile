FROM golang:1.25.0-alpine AS builder

RUN apk add --no-cache bash

WORKDIR /app
COPY . .
RUN env

RUN go build -C cmd/explorer -o /explorer
RUN go build -C cmd/heyo -o /heyo

# RUN PART
FROM alpine:latest

RUN apk add --no-cache bash postgresql-client

WORKDIR /app

COPY --from=builder /explorer .
COPY --from=builder /heyo .

CMD ["./heyo"]
