package services

import (
	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/order"
	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/order/memory"
	accrual "github.com/GorunovAlx/gophermart/internal/gophermart/services/accrual"
)

// OrderConfiguration is an alias for a function that will take in a pointer to an OrderService and modify it
type OrderConfiguration func(os *OrderService) error

// OrderService is a implementation of the OrderService
type OrderService struct {
	orders order.OrderRepository
}

// NewOrderService takes a variable amount of OrderConfiguration functions and returns a new OrderService
// Each OrderConfiguration will be called in the Order they are passed in
func NewOrderService(cfgs ...OrderConfiguration) (*OrderService, error) {
	// Create the Orderservice
	os := &OrderService{}
	// Apply all Configurations passed in
	for _, cfg := range cfgs {
		// Pass the service into the configuration function
		err := cfg(os)
		if err != nil {
			return nil, err
		}
	}
	return os, nil
}

// WithOrderRepository applies a given Order repository to the OrderService
func WithOrderRepository(or order.OrderRepository) OrderConfiguration {
	// return a function that matches the OrderConfiguration alias,
	// You need to return this so that the parent function can take in all the needed parameters
	return func(os *OrderService) error {
		os.orders = or
		return nil
	}
}

// WithMemoryOrderRepository applies a memory Order repository to the OrderService
func WithMemoryOrderRepository() OrderConfiguration {
	// Create the memory repo, if we needed parameters, such as connection strings they could be inputted here
	ur := memory.New()
	return WithOrderRepository(ur)
}

func (os *OrderService) RegisterOrder(orderNumber string, userID int) (int, error) {
	uID := os.orders.GetOrderUserIDByNumber(orderNumber)
	if uID != -1 {
		if uID == userID {
			return -1, order.ErrOrderAlreadyRegisteredByUser
		}
		if uID != userID {
			return -1, order.ErrOrderAlreadyRegisteredByOtherUser
		}
	}

	o := order.NewOrder(orderNumber, userID)
	err := os.orders.Add(o)
	if err != nil {
		return -1, err
	}

	return o.GetID(), nil
}

func (os *OrderService) UpdateOrder(order accrual.AccrualOrder) error {
	err := os.orders.Update(order)
	if err != nil {
		return err
	}
	return nil
}

func (os *OrderService) GetOrdersByUserID(userID int) ([]order.Order, error) {
	res, err := os.orders.GetOrders(userID)
	if err != nil {
		return []order.Order{}, err
	}

	return res, nil
}
