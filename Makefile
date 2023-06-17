DB_NAME=simple_bank
DB_USER=root
DB_PASSWORD=secret
DB_URL="postgresql://$(DB_USER):$(DB_PASSWORD)@localhost:5432/$(DB_NAME)?sslmode=disable"

startdb:
	docker start postgres

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres:alpine

createdb:
	docker exec -it postgres createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

dropdb:
	docker exec -it postgres dropdb $(DB_NAME)

migrateup:
	migrate -path db/migration/ -database $(DB_URL) -verbose up $v

migratedown:
	migrate -path db/migration/ -database $(DB_URL) -verbose down $v

migratecreate:
	migrate create -ext sql -dir db/migration/ -seq $(NAME)

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen --build_flags=--mod=mod -package=mockdb -destination=db/mock/store.go github.com/jotabf/simplebank/db/sqlc Store

.PHONY: startdb postgres createdb dropdb migrateup migratedown migratecreate sqlc test server mock