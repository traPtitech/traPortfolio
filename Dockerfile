# syntax=docker/dockerfile:1

##
## Build stage
##
FROM golang:1.22.2-alpine AS build

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -o /traPortfolio .

##
## Deployment stage
##
FROM alpine:3 AS deploy

WORKDIR /

COPY --from=build /traPortfolio /traPortfolio
COPY dev/config.yaml /opt/traPortfolio/config.yaml

ENTRYPOINT /traPortfolio -c /opt/traPortfolio/config.yaml --db-host mysql
