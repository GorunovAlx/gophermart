package withdraw

import (
	"time"

	"github.com/GorunovAlx/gophermart"
)

type Withdraw struct {
	withdraw *gophermart.Withdraw
}

func NewWithdraw(order string, sum float32, userID int) Withdraw {
	w := &gophermart.Withdraw{
		UserID:       userID,
		Order:        order,
		Sum:          sum,
		Processed_at: time.Now(),
	}

	return Withdraw{
		withdraw: w,
	}
}

type ByUploadedAt []Withdraw

func (b ByUploadedAt) Len() int { return len(b) }

func (b ByUploadedAt) Less(i, j int) bool {
	return b[i].GetProcessedAt().Unix() < b[j].GetProcessedAt().Unix()
}

func (b ByUploadedAt) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

func (w *Withdraw) GetID() int {
	return w.withdraw.Id
}

func (w *Withdraw) SetID(id int) {
	w.withdraw.Id = id
}

func (w *Withdraw) GetUserID() int {
	return w.withdraw.UserID
}

func (w *Withdraw) GetOrder() string {
	return w.withdraw.Order
}

func (w *Withdraw) GetSum() float32 {
	return w.withdraw.Sum
}

func (w *Withdraw) GetProcessedAt() time.Time {
	return w.withdraw.Processed_at
}
