# Go orders serivce (Go + PostgreSQL + Kafka + Redis)

## [Ссылка на демонстрацию](https://disk.yandex.ru/i/BHgx_g9N9tc_xQ)

## Заполнение данных

[Репозиторий с генератором](https://github.com/thavguard/orders-wb-generator)

## Setup .env

create `.env` file

run `docker compose --env-file=.env up -d`

## Linter

`golangci-lint run`

## Working with migrations

Для миграций используется пакет [golang-migrate](https://github.com/golang-migrate/migrate)

Пока что работа только из консоли

В .env создать переменную `DATABASE_URL`

### Create migrations

`make create-migration name=[name]`

### Up migrations

`make up-migration`

### Down migrations

`make down-migration`

## Описание проекта

Language - **go**\
DB - **postgres**\
Driver - **pgx** + **sqlx**\
Broker + DLQ - **kafka**\
HTTP - **gin**\
Cache - **redis LFU**\
Validation - **tags validator**\
Retry - **sethvargo/go-retry**\
