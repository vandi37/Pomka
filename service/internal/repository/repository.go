package repository

import (
	"github.com/sirupsen/logrus"
)

type Repository struct {
	logger *logrus.Logger
}

func NewRepository(logger *logrus.Logger) *Repository {
	return &Repository{logger: logger}
}
