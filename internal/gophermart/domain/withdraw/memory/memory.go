package memory

import (
	"sync"

	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/withdraw"
)

type MemoryWithdrawRepository struct {
	withdrawals map[int]withdraw.Withdraw
	sync.Mutex
}

func New() *MemoryWithdrawRepository {
	return &MemoryWithdrawRepository{
		withdrawals: make(map[int]withdraw.Withdraw),
	}
}

func (mr *MemoryWithdrawRepository) Add(w withdraw.Withdraw) {
	if mr.withdrawals == nil {
		mr.Lock()
		mr.withdrawals = make(map[int]withdraw.Withdraw)
		mr.Unlock()
	}
	withdrawID := getNextID(mr)
	w.SetID(withdrawID)

	mr.Lock()
	mr.withdrawals[withdrawID] = w
	mr.Unlock()
}

func (mr *MemoryWithdrawRepository) GetWithdrawals(userID int) []withdraw.Withdraw {
	var res []withdraw.Withdraw
	for _, withdraw := range mr.withdrawals {
		if withdraw.GetUserID() == userID {
			res = append(res, withdraw)
		}
	}

	return res
}

func getNextID(mr *MemoryWithdrawRepository) int {
	var idCount = 0
	for range mr.withdrawals {
		idCount++
	}

	return idCount + 1
}
