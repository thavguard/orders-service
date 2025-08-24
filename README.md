# Go orders serivce (Go + PostgreSQL + Kafka)

## Setup .env

create `.env` file

run `docker compose --env-file=.env up -d`

## Working with migrations

Для миграций используется пакет [golang-migrate](https://github.com/golang-migrate/migrate)

Пока что работа только из консоли

В .env создать переменную `DATABASE_URL`

после создания / обновления .env выполнить команду `export $(grep -v '^#' .env | xargs)`

### Create migrations

`migrate create -ext sql -dir src/db/migrations -seq [migration name]`

### Up migrations

`migrate -database "$DATABASE_URL" -source file://src/db/migrations up`

### Down migrations

`migrate -database "$DATABASE_URL" -source file://src/db/migrations down`

## Описание проекта

Language - **go**\
DB - **postgres**\
Driver - **pgx** + **sqlx**\
Broker - **kafka**\
HTTP - **gin**\
Cache - **redis LFU**\
