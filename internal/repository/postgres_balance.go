package repository

import (
	"context"
	"errors"
	"strconv"

	"user-service/internal/model"
	"user-service/pkg/postgres"

	"github.com/shopspring/decimal"
)

var NotFoundErr = errors.New("entry not found")
var InsufficientFundsErr = errors.New("insufficient funds")

type DataBaseStorage struct {
	pool *postgres.Pool
}

func NewDataBaseStorage(pool *postgres.Pool) *DataBaseStorage {
	return &DataBaseStorage{pool: pool}
}
func (s *DataBaseStorage) ReplenishmentBalance(ctx context.Context, id int, amount decimal.Decimal) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "UPDATE users SET balance = balance + $1 WHERE id = $2", amount, id)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}
	_, err = tx.Exec(ctx, "INSERT INTO operations (type, recipient, sender, amount) VALUES ($1, $2, $3, $4)", "replenishment", id, 0, amount)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}
	err = tx.Commit(context.Background())
	return err
}

func (s *DataBaseStorage) Payment(ctx context.Context, senderID int, recipientID int, amount decimal.Decimal) error {
	balanceCheck, err := s.checkBalance(ctx, senderID, amount)
	if err != nil {
		return err
	}
	if !balanceCheck {
		return InsufficientFundsErr
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE users SET balance = balance - $1 WHERE id = $2", amount, senderID)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE users SET balance = balance + $1 WHERE id = $2", amount, recipientID)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	_, err = tx.Exec(ctx, "INSERT INTO operations (type, recipient, sender, amount) VALUES ($1, $2, $3, $4)", "payment", recipientID, senderID, amount)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}
	err = tx.Commit(context.Background())
	return err

}

func (s *DataBaseStorage) RecentOperations(ctx context.Context, id int) ([]model.Operation, error) {
	query := `SELECT id, type, recipient, sender, amount, created_at FROM operations WHERE recipient = $1 OR sender = $1 ORDER BY created_at DESC LIMIT 10`

	rows, err := s.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var operations []model.Operation

	for rows.Next() {
		var operation model.Operation
		err := rows.Scan(&operation.ID, &operation.Type, &operation.Recipient, &operation.Sender, &operation.Amount, &operation.CreatedAt)
		if err != nil {
			return operations, err
		}
		operations = append(operations, operation)
	}

	return operations, rows.Err()
}
func (s *DataBaseStorage) checkBalance(ctx context.Context, userID int, amount decimal.Decimal) (bool, error) {
	var balance decimal.Decimal

	query := `SELECT balance FROM users WHERE id = $1`

	err := s.pool.QueryRow(ctx, query, strconv.Itoa(userID)).Scan(&balance)
	if err != nil {
		return false, err
	}
	return balance.GreaterThanOrEqual(amount), nil
}
