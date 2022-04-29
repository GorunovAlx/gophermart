package gophermart

import "time"

type Withdraw struct {
	ID          int
	UserID      int
	Order       string
	Sum         float32
	ProcessedAt time.Time
}
