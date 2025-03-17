package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/AmoabaKelvin/temp-mail/api/internal/db"
	"github.com/AmoabaKelvin/temp-mail/api/internal/store"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .envrc file: %v", err)
	}

	mo_var := os.Getenv("DB_MAX_OPEN_CONNS")
	mo, err := strconv.Atoi(mo_var)
	if err != nil {
		log.Fatalf("Error converting DB_MAX_OPEN_CONNS to int: %v", err)
	}

	mi_var := os.Getenv("DB_MAX_IDLE_CONNS")
	mi, err := strconv.Atoi(mi_var)
	if err != nil {
		log.Fatalf("Error converting DB_MAX_IDLE_CONNS to int: %v", err)
	}

	mit_var := os.Getenv("DB_MAX_IDLE_TIME")
	if mit_var == "" {
		log.Fatalf("Error converting DB_MAX_IDLE_TIME to int: %v", err)
	}

	cfg := config{
		addr: os.Getenv("ADDR"),
		db: dbConfig{
			addr:         os.Getenv("DB_ADDR"),
			maxOpenConns: mo,
			maxIdleConns: mi,
			MaxIdleTime:  mit_var,
		},
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.MaxIdleTime)

	if err != nil {
		log.Fatalf("db.New: %v", err)
	}

	store := store.NewPostgresStorage(db)

	app := application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
