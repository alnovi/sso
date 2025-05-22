FROM golang:1.24.3-alpine3.21 AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -mod=mod -a -installsuffix cgo -o ./app -ldflags="-s -w" ./cmd/server/main.go

FROM alpine:3.21

ARG APP_VERSION=0.0.1

ENV APP_VERSION=$APP_VERSION
ENV APP_ENVIRONMENT=production
ENV HTTP_PORT=8080

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]
