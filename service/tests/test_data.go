package test_grpc

import (
	"promos/internal/models/promos"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var promo = &promos.CreatePromoIn{
	Name:     uuid.New().String(),
	Value:    0,
	Currency: 0,
	Creator:  0,
	ExpAt:    timestamppb.New(time.Now().Add(time.Hour * 6)),
}
