package service

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBService interface {
	Pool() *pgxpool.Pool
	Close()
}

type dbService struct {
	pool *pgxpool.Pool
}

func NewDBService(pool *pgxpool.Pool) (DBService, error) {
	return &dbService{pool: pool}, nil
}

func (s *dbService) Pool() *pgxpool.Pool {
	return s.pool
}

func (s *dbService) Close() {
	s.pool.Close()
}
