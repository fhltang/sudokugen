# suntax=docker/dockerfile:1

## Build
FROM golang:1.18-bullseye as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd/ ./cmd
COPY internal/ ./internal

RUN go build -o /web cmd/web/web.go

RUN find . -print


## Deploy
from gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /web /web

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT [ "/web" ]
