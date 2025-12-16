.PHONY: migrate-up migrate-down migrate-down-all migrate-version migrate-force migrate-create

MIGRATOR_PATH = cmd/migrator/migrator.go

migrate-up:
	go run $(MIGRATOR_PATH) -command=up

migrate-down:
	go run $(MIGRATOR_PATH) -command=down

migrate-down-all:
	go run $(MIGRATOR_PATH) -command=down-all

migrate-version:
	go run $(MIGRATOR_PATH) -command=version

migrate-force:
	go run $(MIGRATOR_PATH) -command=force -force-version=$(V)

gen-mocks:
	mockery

gen-docs:
	swag init -g cmd/service/service.go -o docs