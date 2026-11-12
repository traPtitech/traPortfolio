# syntax=docker/dockerfile:1

##
## Build stage
##
FROM golang:1.25.3-alpine@sha256:aee43c3ccbf24fdffb7295693b6e33b21e01baec1b2a55acc351fde345e9ec34 AS build

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
FROM alpine:3@sha256:f59abe656b692d901e184be3445e084786602108599d2a51d229c0f5a48c71bf AS deploy

WORKDIR /

COPY --from=build /traPortfolio /traPortfolio

ENV TPF_DB_PORT="1323"
EXPOSE ${TPF_DB_PORT}

ENTRYPOINT ["/traPortfolio"]
