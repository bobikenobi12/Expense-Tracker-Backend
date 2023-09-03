dev:
	air

swagger:
	swag init --dir ./,./handlers

create-migration name:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	migrate -path migrations -database "postgres://$(PSQL_USER):$(PSQL_PASSWORD)@localhost:5432/$(PSQL_DATABASE)?sslmode=disable" -verbose up

migrate-down:
	migrate -path migrations -database "postgres://$(PSQL_USER):$(PSQL_PASSWORD)@localhost:5432/$(PSQL_DATABASE)?sslmode=disable" -verbose down

migrate-force:
	migrate -path migrations -database "postgres://$(PSQL_USER):$(PSQL_PASSWORD)@localhost:5432/$(PSQL_DATABASE)?sslmode=disable" -verbose force $(version)

migrate-version:
	migrate -path migrations -database "postgres://$(PSQL_USER):$(PSQL_PASSWORD)@localhost:5432/$(PSQL_DATABASE)?sslmode=disable" -verbose version

migrate-drop:
	migrate -path migrations -database "postgres://$(PSQL_USER):$(PSQL_PASSWORD)@localhost:5432/$(PSQL_DATABASE)?sslmode=disable" -verbose drop

migrate-reset:
	migrate -path migrations -database "postgres://$(PSQL_USER):$(PSQL_PASSWORD)@localhost:5432/$(PSQL_DATABASE)?sslmode=disable" -verbose reset



