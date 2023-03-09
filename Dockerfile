FROM golang:1.20.2-alpine AS build
WORKDIR /go/src/github.com/traPtitech/traPortfolio
COPY ./go.* ./
RUN go mod download
COPY . .
RUN go build -o /traPortfolio .

FROM alpine:3.17.2
WORKDIR /app

COPY --from=build /traPortfolio ./

ENTRYPOINT ./traPortfolio -c /opt/portfolio/config.yaml
