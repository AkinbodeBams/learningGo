include .envrc
MIGRATIONS_PATH = ./cmd/migrate/migrations

.PHONY:migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS)) 

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DBADDR) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DBADDR) down $(filter-out $@,$(MAKECMDGOALS))

