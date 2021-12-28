# Build stage
FROM golang:1.17-buster AS build

WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY src/ /build/src
COPY main.go ./

RUN go build -o /api

# Deploy stage
FROM gcr.io/distroless/base-debian10

WORKDIR /app

COPY --from=build /api /app/api
COPY config.env /app/config.env

EXPOSE 8000

USER nonroot:nonroot

ENTRYPOINT [ "/app/api" ]