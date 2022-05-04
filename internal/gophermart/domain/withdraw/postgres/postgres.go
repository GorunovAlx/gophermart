package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/withdraw"
	"github.com/GorunovAlx/gophermart/internal/gophermart/entity"
)

type PostgresWithdrawRepository struct {
	*pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresWithdrawRepository {
	return &PostgresWithdrawRepository{
		db,
	}
}

func (db *PostgresWithdrawRepository) Add(w withdraw.Withdraw) (int, error) {
	conn, err := db.Acquire(context.Background())
	if err != nil {
		return -1, err
	}
	defer conn.Release()

	insertStatement := `
	INSERT INTO withdrawals (user_id, "order", sum, processed_at)
	VALUES ($1, $2, $3, $4) RETURNING id;`

	var withdrawID int
	err = conn.QueryRow(
		context.Background(),
		insertStatement,
		w.GetUserID(),
		w.GetOrder(),
		w.GetSum(),
		time.Now().Format(time.RFC3339),
	).Scan(&withdrawID)
	if err != nil {
		return -1, err
	}

	return withdrawID, nil
}

func (db *PostgresWithdrawRepository) GetWithdrawals(userID int) []withdraw.Withdraw {
	conn, err := db.Acquire(context.Background())
	if err != nil {
		return []withdraw.Withdraw{}
	}
	defer conn.Release()

	var result []withdraw.Withdraw

	selectStatement := `select id, user_id, "order", sum, processed_at from withdrawals where user_id=$1`
	rows, err := conn.Query(context.Background(), selectStatement, userID)
	if err != nil {
		return []withdraw.Withdraw{}
	}
	defer rows.Close()

	for rows.Next() {
		var w entity.Withdraw
		err = rows.Scan(
			&w.ID,
			&w.UserID,
			&w.Order,
			&w.Sum,
			&w.ProcessedAt,
		)
		if err != nil {
			return []withdraw.Withdraw{}
		}

		aggregateWithdraw := withdraw.NewWithdraw(w.Order, w.Sum, w.UserID)
		aggregateWithdraw.SetID(w.ID)
		aggregateWithdraw.SetProcessedAt(w.ProcessedAt)

		result = append(result, aggregateWithdraw)
	}

	if rows.Err() != nil {
		return []withdraw.Withdraw{}
	}

	return result
}
