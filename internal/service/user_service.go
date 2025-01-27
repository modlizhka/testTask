package service

import (
	"context"

	"user-service/internal/model"

	"github.com/shopspring/decimal"
)

//go:generate mockgen -source=user_service.go -destination=mocks/mock.go
type storage interface {
	ReplenishmentBalance(ctx context.Context, id int, amount decimal.Decimal) error
	Payment(ctx context.Context, senderID int, recipientID int, amount decimal.Decimal) error
	RecentOperations(ctx context.Context, id int) ([]model.Operation, error)
}

type UserService struct {
	dataBaseStorage storage
}

func NewUserService(dataBaseStorage storage) *UserService {
	return &UserService{dataBaseStorage: dataBaseStorage}
}

func (s *UserService) ReplenishmentBalance(ctx context.Context, id int, amount decimal.Decimal) error {
	err := s.dataBaseStorage.ReplenishmentBalance(ctx, id, amount)
	return err
}

func (s *UserService) Payment(ctx context.Context, senderID int, recipientID int, amount decimal.Decimal) error {
	err := s.dataBaseStorage.Payment(ctx, senderID, recipientID, amount)
	return err
}
func (s *UserService) RecentOperations(ctx context.Context, id int) ([]model.Operation, error) {
	res, err := s.dataBaseStorage.RecentOperations(ctx, id)
	return res, err
}
