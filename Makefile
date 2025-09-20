ifneq (,$(wildcard .env))
    include .env
    export
endif


MIGRATION_SOURCE=file://src/db/migrations

no-target:
	@echo "❌ Нужно указать цель!
	@exit 1

create-migration:
	migrate create -ext sql -dir src/db/migrations -seq $(name)

up-migration:
	migrate -database "$(MIGRATE_URL)" -source $(MIGRATION_SOURCE) up

down-migration:
	migrate -database "$(MIGRATE_URL)" -source $(MIGRATION_SOURCE) down