#Отключаем преобразование путей
#У меня локально(на windows) нормально не создаются файлы миграций.
#На других ОС эти команды игнорируются
export MSYS_NO_PATHCONV=1
export MSYS2_ARG_CONV_EXCL=*

include .env
export

env-up:
	@docker compose up -d postgres

env-down:
	@docker compose down postgres

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутсвует необходимый параметр seq. Пример: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up: 
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует параметр action"; exit 1; \
	fi; \
	docker compose run --rm postgres-migrate \
		-path=/migrations/ \
		-database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres:5432/$(POSTGRES_DB)?sslmode=disable" \
		$(action)

kill-all:
	@docker compose down -v --remove-orphans
	@rm -rf ./out
	@echo "Complete."

logs:
	@docker compose logs -f

app-logs:
	@docker compose logs app

up:
	@docker compose up -d --build app

swagger-gen:
	@docker compose run --rm swagger \
		init \
		-g cmd/BoOfIs-app/main.go \
		-o docs \
		--parseInternal \
		--parseDependency

ps:
	@docker compose ps

