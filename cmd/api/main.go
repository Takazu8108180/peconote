package main

import (
	"log"

	"github.com/peconote/peconote/internal/infrastructure/db"
	"github.com/peconote/peconote/internal/infrastructure/router"
)

func main() {
	gormDB, err := db.NewDB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	sqlxDB, err := db.NewSqlxDB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	r := router.NewRouter(gormDB, sqlxDB)
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
