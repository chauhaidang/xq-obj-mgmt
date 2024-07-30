package main

import (
	"log"
)

func main() {
	pgStorage, err := NewPGStorage(Envs)
	if err != nil {
		log.Fatal("Error while initialzing Postgres storage")
	}

	db, err := pgStorage.Init()
	if err != nil {
		log.Fatal("Error while initializing data tables")
	}

	store := NewStore(db)
	api := NewApiServer(Envs.Addr, store)
	api.Serve()
}
