package service_test

import (
	"context"
	"testing"

	"user-service/internal/model"
	"user-service/internal/service"
	mock_service "user-service/internal/service/mocks"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestUserService_ReplenishmentBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock_service.NewMockstorage(ctrl)
	userService := service.NewUserService(mockStorage)

	ctx := context.Background()
	userID := 1
	amount := decimal.NewFromFloat(100.0)

	mockStorage.EXPECT().ReplenishmentBalance(ctx, userID, amount).Return(nil)
	err := userService.ReplenishmentBalance(ctx, userID, amount)
	assert.NoError(t, err)

	mockStorage.EXPECT().ReplenishmentBalance(ctx, userID, amount).Return(assert.AnError)
	err = userService.ReplenishmentBalance(ctx, userID, amount)
	assert.Error(t, err)
}

func TestUserService_Payment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock_service.NewMockstorage(ctrl)
	userService := service.NewUserService(mockStorage)

	ctx := context.Background()
	senderID := 1
	recipientID := 2
	amount := decimal.NewFromFloat(50.0)

	mockStorage.EXPECT().Payment(ctx, senderID, recipientID, amount).Return(nil)
	err := userService.Payment(ctx, senderID, recipientID, amount)
	assert.NoError(t, err)

	mockStorage.EXPECT().Payment(ctx, senderID, recipientID, amount).Return(assert.AnError)
	err = userService.Payment(ctx, senderID, recipientID, amount)
	assert.Error(t, err)
}

func TestUserService_RecentOperations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock_service.NewMockstorage(ctrl)
	userService := service.NewUserService(mockStorage)

	ctx := context.Background()
	userID := 1
	operations := []model.Operation{
		{ID: 1, Amount: decimal.NewFromFloat(100.0)},
		{ID: 2, Amount: decimal.NewFromFloat(50.0)},
	}

	mockStorage.EXPECT().RecentOperations(ctx, userID).Return(operations, nil)
	res, err := userService.RecentOperations(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, operations, res)

	mockStorage.EXPECT().RecentOperations(ctx, userID).Return(nil, assert.AnError)
	res, err = userService.RecentOperations(ctx, userID)
	assert.Error(t, err)
	assert.Nil(t, res)
}
