package services

import (
	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/withdraw"
	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/withdraw/memory"
	withdrawDB "github.com/GorunovAlx/gophermart/internal/gophermart/domain/withdraw/postgres"
	"github.com/jackc/pgx/v4"
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

func WithPostgresWithdrawRepository(pool *pgx.Conn) WithdrawConfiguration {
	return func(ws *WithdrawService) error {
		pwr := withdrawDB.NewPostgresRepository(pool)
		ws.withdrawals = pwr
		return nil
	}
}

func (ws *WithdrawService) Register(order string, sum float32, userID int) (int, error) {
	w := withdraw.NewWithdraw(order, sum, userID)
	id, err := ws.withdrawals.Add(w)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (ws *WithdrawService) GetWithdrawals(userID int) []withdraw.Withdraw {
	return ws.withdrawals.GetWithdrawals(userID)
}
