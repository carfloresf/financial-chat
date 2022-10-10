APP_NAME = financial-chat

rabbitmq:
	docker run -d --hostname my-rabbit --name rabbitmq -p 15672:15672 -p 5672:5672 rabbitmq:3.11.0-management-alpine

run-server:
	go run cmd/server/*.go

run-bot:
	go run cmd/bot/*.go

test:
	./scripts/unit-test.sh

lint:
	./scripts/lint.sh

migrate-create:
	./bin/migrate create -ext sql -dir db/migrations/ $(MIGRATION_NAME)

migrate-up:
	./bin/migrate -path db/migrations/ -database db/file.db up

.PHONY: build run test lint
