package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sodam/internal/handler"
	"sodam/internal/service"

	"github.com/hako/branca"
	_ "github.com/jackc/pgx/stdlib"
)

func main() {
	var (
		port        = env("PORT", "3000")
		origin      = env("ORIGIN", "http://localhost:"+port)
		databaseURL = env("DATABASE_URL", "postgresql://root@127.0.0.1:26257/sodam?sslmode=disable")
		brancaKey   = env("BRANACA_KEY", "supersecretkeyyoushouldnotcommit")
	)

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatalf("could not open db connetion: %v\n", err)
		return
	}

	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Fatalf("could not ping to db: %v\n", err)
		return
	}

	cdc := branca.NewBranca(brancaKey)
	cdc.SetTTL(uint32(service.TokenLifespan.Seconds()))
	s := service.New(db, cdc, origin)
	h := handler.New(s)
	log.Printf("accepting connetions on port %d\n", port)
	if err = http.ListenAndServe(":"+port, h); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}
}

func env(key, fallbackValue string) string {
	s := os.Getenv(key)
	if s == "" {
		return fallbackValue
	}
	return s
}
