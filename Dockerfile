# syntax=docker/dockerfile:1

## Build
FROM golang:1.16-buster AS build

WORKDIR /src

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /simplerick

## Deploy
FROM gcr.io/distroless/static-debian12

COPY --from=build /simplerick /simplerick

EXPOSE 3000

ENTRYPOINT ["/simplerick"]