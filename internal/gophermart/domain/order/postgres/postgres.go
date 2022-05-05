package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/order"
	"github.com/GorunovAlx/gophermart/internal/gophermart/entity"
	accrual "github.com/GorunovAlx/gophermart/internal/gophermart/services/accrual"
)

type PostgresOrderRepository struct {
	*pgx.Conn
}

func NewPostgresRepository(db *pgx.Conn) *PostgresOrderRepository {
	return &PostgresOrderRepository{
		db,
	}
}

func (db *PostgresOrderRepository) Add(o order.Order) (int, error) {
	insertStatement := `
	INSERT INTO orders (user_id, number, status, accrual, uploaded_at)
	VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	var orderID int
	err := db.Conn.QueryRow(
		context.Background(),
		insertStatement,
		o.GetUserID(),
		o.GetNumber(),
		o.GetStatus(),
		o.GetAccrual(),
		time.Now().Format(time.RFC3339),
	).Scan(&orderID)
	if err != nil {
		return -1, err
	}

	return orderID, nil
}

func (db *PostgresOrderRepository) GetOrders(userID int) ([]order.Order, error) {
	var result []order.Order

	selectStatement := "select id, user_id, number, status, accrual, uploaded_at from orders where user_id=$1"
	rows, err := db.Conn.Query(context.Background(), selectStatement, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o entity.Order
		err = rows.Scan(
			&o.ID,
			&o.UserID,
			&o.Number,
			&o.Status,
			&o.Accrual,
			&o.UploadedAt,
		)
		if err != nil {
			return nil, err
		}
		aggregateOrder := order.NewOrder(o.Number, o.UserID)
		aggregateOrder.SetID(o.ID)
		aggregateOrder.SetStatus(o.Status)
		aggregateOrder.SetAccrual(o.Accrual)
		aggregateOrder.SetUploadedAt(o.UploadedAt)

		result = append(result, aggregateOrder)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}

func (db *PostgresOrderRepository) GetOrdersNotProcessed(userID int) ([]order.Order, error) {
	var result []order.Order

	selectStatement := `
	select id, user_id, number, status, accrual, uploaded_at 
	from orders where user_id=$1 and status not in ('INVALID', 'PROCESSED')`

	rows, err := db.Conn.Query(context.Background(), selectStatement, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o entity.Order
		err = rows.Scan(
			&o.ID,
			&o.UserID,
			&o.Number,
			&o.Status,
			&o.Accrual,
			&o.UploadedAt,
		)
		if err != nil {
			return nil, err
		}
		aggregateOrder := order.NewOrder(o.Number, o.UserID)
		aggregateOrder.SetID(o.ID)
		aggregateOrder.SetStatus(o.Status)
		aggregateOrder.SetAccrual(o.Accrual)
		aggregateOrder.SetUploadedAt(o.UploadedAt)

		result = append(result, aggregateOrder)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}

func (db *PostgresOrderRepository) GetOrderUserIDByNumber(orderNumber string) int {
	var userID int
	err := db.Conn.QueryRow(
		context.Background(),
		"select user_id from orders where number=$1",
		orderNumber,
	).Scan(&userID)
	if err != nil {
		return -1
	}

	return userID
}

func (db *PostgresOrderRepository) Update(order accrual.AccrualOrder) error {
	sqlStmt := `
	update orders set status = $1, accrual = $2 where id = $3;`

	_, err := db.Conn.Exec(
		context.Background(),
		sqlStmt,
		order.Status,
		order.Accrual,
		order.OrderID,
	)

	if err != nil {
		return err
	}

	return nil
}
