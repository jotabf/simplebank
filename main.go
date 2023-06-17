package main

import (
	"database/sql"
	"log"

	"github.com/jotabf/simplebank/api"
	db "github.com/jotabf/simplebank/db/sqlc"
	"github.com/jotabf/simplebank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	store := db.NewStore(connDB)
	server, err := api.NewServer(*config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddr)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
