package v1

import (
	"fmt"

	"github.com/GorunovAlx/gophermart/config"
	"github.com/GorunovAlx/gophermart/internal/gophermart/database"
	accrualService "github.com/GorunovAlx/gophermart/internal/gophermart/services/accrual"
	loyaltyService "github.com/GorunovAlx/gophermart/internal/gophermart/services/loyalty"
	orderService "github.com/GorunovAlx/gophermart/internal/gophermart/services/order"
	userService "github.com/GorunovAlx/gophermart/internal/gophermart/services/user"
	withdrawService "github.com/GorunovAlx/gophermart/internal/gophermart/services/withdraw"
)

type ServiceShelf struct {
	Accruals    *accrualService.AccrualService
	Users       *userService.UserService
	Orders      *orderService.OrderService
	Withdrawals *withdrawService.WithdrawService
	Loyalty     *loyaltyService.LoyaltySystem
}

func NewServiceShelf(cfg *config.Config, st *database.Storage) (*ServiceShelf, error) {
	if st == nil {
		return NewShelfWithMemoryStorage(cfg)
	}

	return NewShelfWithPostgresStorage(cfg, st)
}

func NewShelfWithPostgresStorage(cfg *config.Config, st *database.Storage) (*ServiceShelf, error) {
	us, err := userService.NewUserService(
		userService.WithPostgresUserRepository(st.PGpool),
	)
	if err != nil {
		return nil, fmt.Errorf("NewShelfWithPostgresStorage, NewUserService: %v", err.Error())
	}

	os, err := orderService.NewOrderService(
		orderService.WithPostgresOrderRepository(st.PGpool),
	)
	if err != nil {
		return nil, fmt.Errorf("NewShelfWithPostgresStorage, NewOrderService: %v", err.Error())
	}

	ws, err := withdrawService.NewWithdrawService(
		withdrawService.WithPostgresWithdrawRepository(st.PGpool),
	)
	if err != nil {
		return nil, fmt.Errorf("NewShelfWithPostgresStorage, NewWithdrawService: %v", err.Error())
	}

	as := accrualService.NewAccrualService(cfg.AccrualAddress)

	ls, err := loyaltyService.NewLoyaltySystem(
		loyaltyService.WithUserService(us),
		loyaltyService.WithOrderService(os),
		loyaltyService.WithAccrualService(as),
		loyaltyService.WithWithdrawService(ws),
	)
	if err != nil {
		return nil, fmt.Errorf("NewShelfWithPostgresStorage, NewLoyaltySystem: %v", err.Error())
	}

	return &ServiceShelf{
		Accruals:    as,
		Users:       us,
		Orders:      os,
		Withdrawals: ws,
		Loyalty:     ls,
	}, nil
}

func NewShelfWithMemoryStorage(cfg *config.Config) (*ServiceShelf, error) {
	us, err := userService.NewUserService(
		userService.WithMemoryUserRepository(),
	)
	if err != nil {
		return nil, fmt.Errorf("NewShelfWithMemoryStorage, NewUserService: %v", err.Error())
	}

	os, err := orderService.NewOrderService(
		orderService.WithMemoryOrderRepository(),
	)
	if err != nil {
		return nil, fmt.Errorf("NewShelfWithMemoryStorage, NewOrderService: %v", err.Error())
	}

	ws, err := withdrawService.NewWithdrawService(
		withdrawService.WithMemoryWithdrawRepository(),
	)
	if err != nil {
		return nil, fmt.Errorf("NewShelfWithMemoryStorage, NewWithdrawService: %v", err.Error())
	}

	as := accrualService.NewAccrualService(cfg.AccrualAddress)

	ls, err := loyaltyService.NewLoyaltySystem(
		loyaltyService.WithUserService(us),
		loyaltyService.WithOrderService(os),
		loyaltyService.WithAccrualService(as),
		loyaltyService.WithWithdrawService(ws),
	)
	if err != nil {
		return nil, fmt.Errorf("NewShelfWithMemoryStorage, NewLoyaltySystem: %v", err.Error())
	}

	return &ServiceShelf{
		Accruals:    as,
		Users:       us,
		Orders:      os,
		Withdrawals: ws,
		Loyalty:     ls,
	}, nil
}
