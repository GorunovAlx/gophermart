package services

import (
	"errors"
	"time"

	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/order"
	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/user"
	accrualService "github.com/GorunovAlx/gophermart/internal/gophermart/services/accrual"
	orderService "github.com/GorunovAlx/gophermart/internal/gophermart/services/order"
	userService "github.com/GorunovAlx/gophermart/internal/gophermart/services/user"
	withdrawService "github.com/GorunovAlx/gophermart/internal/gophermart/services/withdraw"
)

// LoyaltySystemConfiguration is an alias that takes a pointer and modifies the LoyaltySystem
type LoyaltySystemConfiguration func(ls *LoyaltySystem) error

type LoyaltySystem struct {
	// userservice is used to handle users
	UserService *userService.UserService
	// orderservice is used to handle orders
	OrderService    *orderService.OrderService
	WithdrawService *withdrawService.WithdrawService
	AccrualService  *accrualService.AccrualService
}

const (
	shortWait       = 500 * time.Nanosecond
	longWait        = 3 * 500 * time.Nanosecond
	processedStatus = "PROCESSED"
	invalidStatus   = "INVALID"
)

// NewLoyaltySystem takes a variable amount of LoyaltySystemConfigurations and builds a LoyaltySystem
func NewLoyaltySystem(cfgs ...LoyaltySystemConfiguration) (*LoyaltySystem, error) {
	// Create the LoyaltySystem
	ls := &LoyaltySystem{}
	// Apply all Configurations passed in
	for _, cfg := range cfgs {
		// Pass the service into the configuration function
		err := cfg(ls)
		if err != nil {
			return nil, err
		}
	}
	return ls, nil
}

//
func WithUserService(us *userService.UserService) LoyaltySystemConfiguration {
	return func(ls *LoyaltySystem) error {
		ls.UserService = us
		return nil
	}
}

// WithOrderService applies a given OrderService to the LoyaltySystem
func WithOrderService(os *orderService.OrderService) LoyaltySystemConfiguration {
	// return a function that matches the LoyaltySystemConfiguration signature
	return func(ls *LoyaltySystem) error {
		ls.OrderService = os
		return nil
	}
}

func WithWithdrawService(ws *withdrawService.WithdrawService) LoyaltySystemConfiguration {
	// return a function that matches the LoyaltySystemConfiguration signature
	return func(ls *LoyaltySystem) error {
		ls.WithdrawService = ws
		return nil
	}
}

func WithAccrualService(as *accrualService.AccrualService) LoyaltySystemConfiguration {
	return func(ls *LoyaltySystem) error {
		ls.AccrualService = as
		return nil
	}
}

func (ls *LoyaltySystem) UpdateOrder(job *order.OrderJob) error {
	done := make(chan bool)
	duration := shortWait
	for {
		acOrder, err := ls.AccrualService.GetAccrualOrder(job.Number)
		if err != nil {
			if errors.Is(err, accrualService.ErrDataProcessing) {
				return accrualService.ErrDataProcessing
			}
			if errors.Is(err, accrualService.ErrDataRetrievalError) || errors.Is(err, accrualService.ErrStatusNotOk) {
				duration = longWait
			}
		}

		acOrder.OrderID = job.ID
		if acOrder.Status != job.Status {
			err = ls.OrderService.UpdateOrder(acOrder)
			if errors.Is(err, order.ErrUpdateOrderNotExists) {
				return order.ErrUpdateOrderNotExists
			}
		}
		if acOrder.Status == processedStatus || acOrder.Status == invalidStatus {
			if acOrder.Accrual > 0 {
				err = ls.UserService.ChangeBalance(job.UserID, acOrder.Accrual)
				if errors.Is(err, user.ErrUserNotFound) {
					return user.ErrUserNotFound
				}
			}
			done <- true
		}

		select {
		case <-done:
			return nil
		default:
			time.Sleep(duration)
		}
	}
}

func (ls *LoyaltySystem) Update(number string, orderID int, userID int) error {
	acOrder, err := ls.AccrualService.GetAccrualOrder(number)
	if err != nil {
		if errors.Is(err, accrualService.ErrDataProcessing) {
			return accrualService.ErrDataProcessing
		}
	}

	acOrder.OrderID = orderID
	err = ls.OrderService.UpdateOrder(acOrder)
	if errors.Is(err, order.ErrUpdateOrderNotExists) {
		return order.ErrUpdateOrderNotExists
	}
	if acOrder.Status == processedStatus || acOrder.Status == invalidStatus {
		if acOrder.Accrual > 0 {
			err = ls.UserService.ChangeBalance(userID, acOrder.Accrual)
			if errors.Is(err, user.ErrUserNotFound) {
				return user.ErrUserNotFound
			}
		}
	}

	return nil
}

func (ls *LoyaltySystem) RegisterWithdraw(order string, sum float32, userID int) error {
	u, err := ls.UserService.GetUser(userID)
	if err != nil {
		return err
	}
	if u.GetCurrentBalance() < sum {
		return user.ErrNotEnoughFunds
	}

	ls.WithdrawService.Register(order, sum, userID)
	err = ls.UserService.TakeOutSum(userID, sum)
	if err != nil {
		return err
	}

	return nil
}
