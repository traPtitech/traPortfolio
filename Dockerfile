# syntax=docker/dockerfile:1

##
## Build stage
##
FROM golang:1.23.2-alpine AS build

WORKDIR /app

RUN --mount=type=cache,target=/go/pkg/mod/,sharing=locked \
  --mount=type=bind,source=go.sum,target=go.sum \
  --mount=type=bind,source=go.mod,target=go.mod \
  go mod download

RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,target=. \
  go build -o /traPortfolio

##
## Deployment stage
##
FROM alpine:3 AS deploy

WORKDIR /

COPY --from=build /traPortfolio /traPortfolio

ENV TPF_DB_PORT="1323"
EXPOSE ${TPF_DB_PORT}

ENTRYPOINT ["/traPortfolio"]
