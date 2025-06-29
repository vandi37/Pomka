package service

import (
	"conn"
	"protobuf/checks"
	"protobuf/market"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceMarket struct {
	RepositoryMarket
	db *pgxpool.Pool
	conn.UserService
	checks.UnimplementedChecksServer
	market.UnimplementedMarketServer
}

type RepositoryMarket interface {
}

func NewServiceMarket(repo RepositoryMarket, db *pgxpool.Pool, users conn.UserService) *ServiceMarket {
	return &ServiceMarket{RepositoryMarket: repo, db: db, UserService: users}
}
