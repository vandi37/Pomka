package repository

import "utils/hasher"

type Repository struct {
	h *hasher.Hasher
}

func NewRepository(hasher *hasher.Hasher) *Repository {
	return &Repository{h: hasher}
}
