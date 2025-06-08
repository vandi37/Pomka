package test_grpc

import (
	"promos/internal/models/promos"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var promo = &promos.CreatePromoIn{
	Name:     uuid.New().String(),
	Value:    10,
	Currency: 1,
	Creator:  "test",
	ExpAt:    timestamppb.New(time.Now().Add(time.Hour * 6)),
}
