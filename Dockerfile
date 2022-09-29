FROM golang:1.16.2-alpine AS build
ENV CGO_ENABLED=0

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG MODE="docker"

RUN set -ex; \
    if [[ ${MODE} == "dev" ]]; then mv .env.example .env; \
    elif [[ ${MODE} == "docker" ]]; then mv .env.docker .env ; \
    else mv .env.testing .env; fi; \
    mkdir /app; \
    cp .env /app/.env

RUN go build -o /app/xkcd

FROM docker.io/improwised/golang-base
COPY --from=build /app/ /app/
ENTRYPOINT ["/app/xkcd"]