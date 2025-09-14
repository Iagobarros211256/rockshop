package main

import (
	"log"
	"os"

	"github.com/Iagobarros211256/rockshop/internal/store"
)

func main() {
	dataFile := os.Getenv("DB_FILE")
	if datafile == "" {
		datafile = "data/store.json"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	st, err := store.NewJSONStore(datafile)
	if err != nil {
		log.Fatalf("faailed to init store: %v", err)
	}
}
