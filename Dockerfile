FROM golang:1.20.2-alpine AS build
WORKDIR /go/src/github.com/traPtitech/traPortfolio
COPY ./go.* ./
RUN go mod download
COPY . .
RUN go build -o /traPortfolio .

FROM alpine:3.12.0
WORKDIR /app

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

COPY --from=build /traPortfolio ./

ENTRYPOINT ./traPortfolio
