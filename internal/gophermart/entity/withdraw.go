package entity

import (
	"context"
	"time"

	"github.com/GorunovAlx/gophermart/internal/gophermart/database"
)

type Withdraw struct {
	ID          int
	UserID      int
	Order       string
	Sum         float32
	ProcessedAt time.Time
}

type WithdrawRepository interface {
	Add(userID int, order string, sum float32) error
	GetWithdrawals(userID int) ([]Withdraw, error)
	GetUserSumWithdrawn(userID int) (float32, error)
}

type WithdrawStorage struct {
	S database.Storage
}

func (ws WithdrawStorage) Add(userID int, order string, sum float32) error {
	insertStatement := `
	INSERT INTO withdrawals (user_id, "order", sum)
	VALUES ($1, $2, $3);`

	_, err := ws.S.PGpool.Exec(
		context.Background(),
		insertStatement,
		userID,
		order,
		sum,
	)
	if err != nil {
		return err
	}

	return nil
}

func (ws WithdrawStorage) GetWithdrawals(userID int) ([]Withdraw, error) {
	var result []Withdraw

	selectStatement := `select id, user_id, "order", sum, processed_at from withdrawals where user_id=$1`
	rows, err := ws.S.PGpool.Query(context.Background(), selectStatement, userID)
	if err != nil {
		return []Withdraw{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var w Withdraw
		err = rows.Scan(
			&w.ID,
			&w.UserID,
			&w.Order,
			&w.Sum,
			&w.ProcessedAt,
		)
		if err != nil {
			return []Withdraw{}, err
		}

		result = append(result, w)
	}

	if rows.Err() != nil {
		return []Withdraw{}, rows.Err()
	}

	return result, nil
}

func (ws WithdrawStorage) GetUserSumWithdrawn(userID int) (float32, error) {
	var withdrawn float32
	err := ws.S.PGpool.QueryRow(
		context.Background(),
		"select coalesce(sum(withdrawals.sum), 0) from withdrawals where withdrawals.user_id = $1;",
		userID,
	).Scan(&withdrawn)
	if err != nil {
		return 0, err
	}

	return withdrawn, nil
}
