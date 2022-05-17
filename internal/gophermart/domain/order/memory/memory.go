package memory

import (
	"sync"
	"time"

	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/order"
	accrual "github.com/GorunovAlx/gophermart/internal/gophermart/services/accrual"
)

type MemoryOrderRepository struct {
	orders map[int]order.Order
	sync.Mutex
}

func New() *MemoryOrderRepository {
	return &MemoryOrderRepository{
		orders: make(map[int]order.Order),
	}
}

func (mr *MemoryOrderRepository) Add(o order.Order) (int, error) {
	if mr.orders == nil {
		mr.Lock()
		mr.orders = make(map[int]order.Order)
		mr.Unlock()
	}
	orderID := getNextID(mr)
	o.SetID(orderID)
	o.SetUploadedAt(time.Now())

	mr.Lock()
	mr.orders[orderID] = o
	mr.Unlock()

	return orderID, nil
}

func (mr *MemoryOrderRepository) Update(accrualOrder accrual.AccrualOrder) error {
	if _, ok := mr.orders[accrualOrder.OrderID]; !ok {
		return order.ErrUpdateOrderNotExists
	}

	mr.Lock()
	o := mr.orders[accrualOrder.OrderID]
	o.SetStatus(accrualOrder.Status)
	o.SetAccrual(accrualOrder.Accrual)
	mr.Unlock()

	return nil
}

func (mr *MemoryOrderRepository) GetOrderUserIDByNumber(orderNumber string) int {
	for _, order := range mr.orders {
		if order.GetNumber() == orderNumber {
			return order.GetUserID()
		}
	}

	return -1
}

func (mr *MemoryOrderRepository) GetOrders(userID int) ([]order.Order, error) {
	var res []order.Order
	for _, order := range mr.orders {
		if order.GetUserID() == userID {
			res = append(res, order)
		}
	}

	return res, nil
}

func (mr *MemoryOrderRepository) GetOrdersNotProcessed(userID int) ([]order.Order, error) {
	var res []order.Order
	for _, order := range mr.orders {
		status := order.GetStatus()
		if status != "PROCESSED" && status != "INVALID" {
			res = append(res, order)
		}
	}

	return res, nil
}

func getNextID(mr *MemoryOrderRepository) int {
	var idCount = 0
	for range mr.orders {
		idCount++
	}

	return idCount + 1
}
