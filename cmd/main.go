package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/go-redis/redis/v9"
	_ "github.com/lib/pq"


	"github.com/mirasildev/note_project/api"
	"github.com/mirasildev/note_project/config"
	"github.com/mirasildev/note_project/storage"
)

func main() {
	cfg := config.Load(".")
	fmt.Println(cfg)
	psqlUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)

	psqlConn, err := sqlx.Connect("postgres", psqlUrl)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
	})

	strg := storage.NewStoragePg(psqlConn)
	inMemory := storage.NewInMemoryStorage(rdb)

	apiServer := api.New(&api.RouterOptions{
		Cfg: &cfg,
		Storage: strg,
		InMemory: inMemory,
	})

	err = apiServer.Run(cfg.HttpPort)
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

	log.Print("Server stopped")

}
