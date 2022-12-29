# syntax=docker/dockerfile:1

## Build
FROM golang:1.16-buster AS build

WORKDIR /src

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /simplerick

## Deploy
FROM gcr.io/distroless/base-debian10

COPY --from=build /simplerick /simplerick

EXPOSE 3000

ENTRYPOINT ["/simplerick"]