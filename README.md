# Single Sign-On (SSO)

[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/alnovi/sso)](https://go.dev/dl/)
[![GitHub License](https://img.shields.io/github/license/alnovi/sso)](https://github.com/alnovi/sso/blob/master/LICENSE.md)
[![Go Report Card](https://goreportcard.com/badge/github.com/alnovi/sso)](https://goreportcard.com/report/github.com/alnovi/sso)
![coverage](https://raw.githubusercontent.com/alnovi/sso/badges/.badges/master/coverage.svg)
![GitHub top language](https://img.shields.io/github/languages/top/alnovi/sso)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/alnovi/sso)
[![GitHub Release](https://img.shields.io/github/v/release/alnovi/sso)](https://github.com/alnovi/sso/releases)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/alnovi/sso/master.yml)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/alnovi/sso/deploy.yml?label=deploy)

**SSO**, расшифровывается как **Single Sign-On**, — это технология аутентификации, которая
позволяет пользователю войти в несколько связанных сервисов с единым логином и паролем.

## Генерация сертификатов

Для работы приложения требуются RSA сертификаты. При первом запуске приложения, если сертификатов нет,
они автоматически будут сгенерированы в папке certs. Вы можете сгенерировать пару открытого и закрытого ключей RSA и
сохранить их в папке certs (_private.pem_ и _public.pem_), которая монтируется в docker.

## Запуск в docker compose

Для работы приложения требуется СУБД postgres, подключить папку для сертификатов
и указать переменные окружения.

```yml
services:
  server:
    image: ghcr.io/alnovi/sso:latest
    restart: unless-stopped
    volumes:
      - ~/.sso/data/certs:/app/certs:rw
    ports:
      - "8080:8080"
    environment:
      APP_HOST: http://127.0.0.1:8080
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USERNAME: user
      DB_PASSWORD: secret
      MAIL_USERNAME: admin@example.com
      MAIL_PASSWORD: secret
      CLIENT_ADMIN_SECRET: secret
      USER_ADMIN_EMAIL: admin@example.com
      USER_ADMIN_PASSWORD: secret
    depends_on:
      - postgres
    networks:
      - backend

  postgres:
    image: postgres:16-alpine3.21
    restart: unless-stopped
    volumes:
      - ~/.sso/data/postgres:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    environment:
      POSTGRES_DB: sso
      POSTGRES_USER: user
      POSTGRES_PASSWORD: secret
    networks:
      - backend

networks:
  backend:
    driver: bridge
```
## Переменные окружения

Сервис SSO можно настраивать с использованием переменных окружения. Для обеспечения безопасности заполните следующие ключи:
`CLIENT_ADMIN_SECRET`, `USER_ADMIN_EMAIL`, `USER_ADMIN_PASSWORD`.

| Key                            | Require | Default           | Description                                    |
|:-------------------------------|:-------:|:------------------|:-----------------------------------------------|
| APP_HOST                       |   Да    |                   | Хост на котором работает сервис                |
| APP_SHUTDOWN                   |   Нет   | 10s               | Максимальное время остановки сервиса           |
| LOG_FORMAT                     |   Нет   | json              | Формат логов (text, json, pretty, discard)     |
| LOG_LEVEL                      |   Нет   | error             | Уровень логирования (debug, info, warn, error) |
| HTTP_HOST                      |   Нет   | 0.0.0.0           | Хост HTTP сервера                              |
| HTTP_PORT                      |   Нет   | 8080              | Порт HTTP сервера                              |
| DB_HOST                        |   Нет   | localhost         | Хост СУБД postgres                             |
| DB_PORT                        |   Нет   | 5432              | Порт СУБД postgres                             |
| DB_USERNAME                    |   Нет   | root              | Пользователь СУБД postgres                     |
| DB_PASSWORD                    |   Нет   | secret            | Пароль пользователя СУБД postgres              |
| DB_DATABASE                    |   Нет   | sso               | Название БД postgres                           |
| MAIL_HOST                      |   Нет   | smtp.gmail.com    | Хост почтового сервера                         |
| MAIL_PORT                      |   Нет   | 587               | Порт почтового сервера                         |
| MAIL_FROM                      |   Нет   | SSO               | Имя отправителя                                |
| MAIL_USERNAME                  |   Да    | sso@example.com   | Пользователь почтового сервера                 |
| MAIL_PASSWORD                  |   Да    | secret            | Пароль пользователя почтового сервера          |
| SCHEDULER_STOP_TIMEOUT         |   Нет   | 5s                | Максимальное время остановки планировщика      |
| SCHEDULER_DELETE_TOKEN_EXPIRED |   Нет   | 5m                | Интервал удаления не активных токенов          |
| SCHEDULER_DELETE_SESSION_EMPTY |   Нет   | 5m                | Интервал удаления не активных сессий           |
| CLIENT_ADMIN_ID                |   Нет   | sso-admin         | Client ID админки                              |
| CLIENT_ADMIN_NAME              |   Нет   | Пользователи      | Client name админки                            |
| CLIENT_ADMIN_SECRET            |   Да    | secret            | Client secret админки                          |
| USER_ADMIN_NAME                |   Нет   | Admin             | Имя admin пользователя                         |
| USER_ADMIN_EMAIL               |   Да    | admin@example.com | Логин admin пользователя                       |
| USER_ADMIN_PASSWORD            |   Да    | secret            | Пароль admin пользователя                      |
