package main

import (
	"context"
	"log"

	"github.com/dtt4h/go-url-shortener/internal/config"
	//urlRepo "github.com/dtt4h/go-url-shortener/internal/repository/url"
	//"github.com/dtt4h/go-url-shortener/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.MustLoad()

	pool, err := pgxpool.New(context.Background(), cfg.DB.URL)
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()
	//dbService, err := service.NewDBService(pool)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//urlRepository := urlRepo.NewURLRepository(dbService.Pool())
	//urlService := service.NewURLService(urlRepository)
}
