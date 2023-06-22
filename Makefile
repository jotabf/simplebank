DB_NAME=simple_bank
DB_USER=root
DB_PASSWORD=secret
DB_URL="postgresql://$(DB_USER):$(DB_PASSWORD)@localhost:5432/$(DB_NAME)?sslmode=disable"

startdb:
	docker start postgres

postgres:
	docker run --name postgres --network simplebank-net -p 5432:5432 -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres:alpine

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

db_docs:
	dbdocs build doc/db.dbml 

db_schema:
	dbml2sql  --postgres -o doc/schema.sql doc/db.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen --build_flags=--mod=mod -package=mockdb -destination=db/mock/store.go github.com/jotabf/simplebank/db/sqlc Store

image: 
	docker build -t simplebank:latest .

container:
	docker run --name simplebank --network simplebank-net -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:secret@postgres:5432/simple_bank?sslmode=disable" simplebank:latest

.PHONY: startdb postgres createdb dropdb migrateup migratedown migratecreate db_docs db_schema sqlc test server mock image container