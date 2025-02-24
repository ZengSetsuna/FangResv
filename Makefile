postgres:
	docker run --name pgfr -p 5433:5432 -e POSTGRES_USER=shu -e POSTGRES_PASSWORD=shu -d postgres:12-alpine

createdb:
	docker exec -it pgfr createdb --username=shu --owner=shu fr

dropdb:
	docker exec -it pgfr dropdb --username=shu --owner=shu fr

migrateup:
	migrate -path db/migration -database "postgresql://shu:shu@localhost:5433/fr?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://shu:shu@localhost:5433/fr?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://shu:shu@localhost:5433/fr?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://shu:shu@localhost:5433/fr?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb sqlc test migrateup migratedown server migratedown1 migrateup1