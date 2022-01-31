FROM golang:1.16-alpine as build

WORKDIR /pow

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . /pow

RUN go build -o "./bin/app" ./cmd/app

FROM alpine:3.15

WORKDIR /

COPY --from=build /pow/bin/app /app
COPY ./pkg/quote/quotes.csv /pkg/quote/quotes.csv

ENTRYPOINT ["/app"]
