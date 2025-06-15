package service

import (
	"warns/pkg/models/warns"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceWarns struct {
	repo RepositoryWarns
	db   *pgxpool.Pool
	warns.UnimplementedWarnsServer
}

type RepositoryWarns interface {
}

func NewServiceWarns(repo RepositoryWarns, db *pgxpool.Pool) *ServiceWarns {
	return &ServiceWarns{repo: repo, db: db}
}
