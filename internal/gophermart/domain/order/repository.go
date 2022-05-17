package order

import (
	"errors"

	accrual "github.com/GorunovAlx/gophermart/internal/gophermart/services/accrual"
)

var (
	ErrOrderAlreadyRegisteredByUser      = errors.New("the order number has already been uploaded by this user")
	ErrOrderAlreadyRegisteredByOtherUser = errors.New("the order number has already been uploaded by another user")
	ErrUpdateOrderNotExists              = errors.New("failed to update the order in the repository: order does not exist")
)

type OrderRepository interface {
	Add(Order) (int, error)
	GetOrders(userID int) ([]Order, error)
	GetOrdersNotProcessed(userID int) ([]Order, error)
	GetOrderUserIDByNumber(orderNumber string) int
	Update(order accrual.AccrualOrder) error
}
