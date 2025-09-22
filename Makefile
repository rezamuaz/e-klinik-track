.PHONY: up
up:
	docker compose up -d

.PHONY: start
start: up wait-db migrate-up
	@echo "App stack is up and migrations are applied."

migration_up: 
	migrate -path migrations/sql/ -database "postgresql://destiny:qn8prVZ6Cr75@localhost:5433/pvsave?sslmode=disable" -verbose up
migration_drop:
	migrate -path migrations/sql/ -database "postgresql://destiny:qn8prVZ6Cr75@localhost:5433/pvsave?sslmode=disable" -verbose drop