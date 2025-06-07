FROM node:24-alpine3.21 as node-builder

WORKDIR /app

COPY ./web .

RUN npm install
RUN npm run build

FROM golang:1.24.3-alpine3.21 AS go-builder

ARG APP_VERSION=0.0.1

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /app

COPY --from=node-builder /app ./web

RUN go mod download
RUN go build -mod=mod -a -installsuffix cgo -ldflags="-s -w -X 'github.com/alnovi/sso/config.Version=$APP_VERSION'" -o ./app ./cmd/server/main.go

FROM alpine:3.21

ENV APP_ENVIRONMENT=production
ENV HTTP_PORT=8080

VOLUME /app/certs

WORKDIR /app

COPY --from=go-builder /app/app .

EXPOSE 8080

CMD ["./app"]
