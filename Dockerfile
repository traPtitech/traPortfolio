# syntax=docker/dockerfile:1

##
## Build stage
##
FROM golang:1.22.4-alpine AS build

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

ENTRYPOINT /traPortfolio
