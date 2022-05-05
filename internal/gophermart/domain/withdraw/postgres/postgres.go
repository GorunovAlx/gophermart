package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/withdraw"
	"github.com/GorunovAlx/gophermart/internal/gophermart/entity"
)

type PostgresWithdrawRepository struct {
	*pgx.Conn
}

func NewPostgresRepository(db *pgx.Conn) *PostgresWithdrawRepository {
	return &PostgresWithdrawRepository{
		db,
	}
}

func (db *PostgresWithdrawRepository) Add(w withdraw.Withdraw) (int, error) {
	insertStatement := `
	INSERT INTO withdrawals (user_id, "order", sum, processed_at)
	VALUES ($1, $2, $3, $4) RETURNING id;`

	var withdrawID int
	err := db.Conn.QueryRow(
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
	var result []withdraw.Withdraw

	selectStatement := `select id, user_id, "order", sum, processed_at from withdrawals where user_id=$1`
	rows, err := db.Conn.Query(context.Background(), selectStatement, userID)
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
