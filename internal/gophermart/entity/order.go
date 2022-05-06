package entity

import (
	"context"
	"errors"
	"time"

	"github.com/GorunovAlx/gophermart/internal/gophermart/accrual"
	"github.com/GorunovAlx/gophermart/internal/gophermart/database"
)

var (
	ErrOrderAlreadyRegisteredByUser      = errors.New("the order number has already been uploaded by this user")
	ErrOrderAlreadyRegisteredByOtherUser = errors.New("the order number has already been uploaded by another user")
	ErrUpdateOrderNotExists              = errors.New("failed to update the order in the repository: order does not exist")
)

type Order struct {
	ID         int
	UserID     int
	Number     string
	Accrual    float32
	Status     string
	UploadedAt time.Time
}

type OrderStorage struct {
	S database.Storage
}

type OrderRepository interface {
	Add(userID int, accrual float32, status, number string) (int, error)
	GetOrders(userID int) ([]Order, error)
	GetOrdersNotProcessed(userID int) ([]Order, error)
	GetOrderUserIDByNumber(orderNumber string) int
	Update(order accrual.AccrualOrder, userID int) error
}

func (os OrderStorage) Add(userID int, accrual float32, status, number string) (int, error) {
	insertStatement := `
	INSERT INTO orders (user_id, number, status, accrual, uploaded_at)
	VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	var orderID int
	err := os.S.PGpool.QueryRow(
		context.Background(),
		insertStatement,
		userID,
		number,
		status,
		accrual,
		time.Now().Format(time.RFC3339),
	).Scan(&orderID)
	if err != nil {
		return -1, err
	}

	return orderID, nil
}

func (os OrderStorage) GetOrders(userID int) ([]Order, error) {
	var result []Order

	selectStatement := "select id, user_id, number, status, accrual, uploaded_at from orders where user_id=$1"
	rows, err := os.S.PGpool.Query(context.Background(), selectStatement, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o Order
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

		result = append(result, o)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}

func (os OrderStorage) GetOrdersNotProcessed(userID int) ([]Order, error) {
	var result []Order

	selectStatement := `
	select id, user_id, number, status, accrual, uploaded_at 
	from orders where user_id=$1 and status not in ('INVALID', 'PROCESSED')`

	rows, err := os.S.PGpool.Query(context.Background(), selectStatement, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o Order
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

		result = append(result, o)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}

func (os OrderStorage) GetOrderUserIDByNumber(orderNumber string) int {
	var userID int
	err := os.S.PGpool.QueryRow(
		context.Background(),
		"select user_id from orders where number=$1",
		orderNumber,
	).Scan(&userID)
	if err != nil {
		return -1
	}

	return userID
}

func (os OrderStorage) Update(order accrual.AccrualOrder, userID int) error {
	sqlStmt := `
	update orders set status = $1, accrual = $2 where user_id = $3;`

	_, err := os.S.PGpool.Exec(
		context.Background(),
		sqlStmt,
		order.Status,
		order.Accrual,
		userID,
	)

	if err != nil {
		return err
	}

	return nil
}
