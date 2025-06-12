package test_grpc

import (
	"promos/internal/models/promos"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var promo = &promos.CreatePromo{
	Name:     uuid.New().String(),
	Amount:   0,
	Currency: 0,
	Uses:     1,
	ExpAt:    timestamppb.New(time.Now().Add(time.Hour * 6)),
}
