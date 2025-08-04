package main

import (
	"log"

	"github.com/peconote/peconote/internal/infrastructure/db"
	"github.com/peconote/peconote/internal/infrastructure/router"
)

func main() {
	d, err := db.NewDB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	r := router.NewRouter(d)
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
