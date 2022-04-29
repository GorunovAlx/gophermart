package services

import (
	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/withdraw"
	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/withdraw/memory"
)

type WithdrawConfiguration func(ws *WithdrawService) error

type WithdrawService struct {
	withdrawals withdraw.WithdrawRepository
}

func NewWithdrawService(cfgs ...WithdrawConfiguration) (*WithdrawService, error) {
	ws := &WithdrawService{}
	for _, cfg := range cfgs {
		err := cfg(ws)
		if err != nil {
			return nil, err
		}
	}
	return ws, nil
}

func WithWithdrawRepository(wr withdraw.WithdrawRepository) WithdrawConfiguration {
	return func(ws *WithdrawService) error {
		ws.withdrawals = wr
		return nil
	}
}

func WithMemoryWithdrawRepository() WithdrawConfiguration {
	wr := memory.New()
	return WithWithdrawRepository(wr)
}

func (ws *WithdrawService) Register(order string, sum float32, userID int) {
	w := withdraw.NewWithdraw(order, sum, userID)
	ws.withdrawals.Add(w)
}

func (ws *WithdrawService) GetWithdrawals(userID int) []withdraw.Withdraw {
	return ws.withdrawals.GetWithdrawals(userID)
}
