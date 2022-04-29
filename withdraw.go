package gophermart

import "time"

type Withdraw struct {
	Id           int
	UserID       int
	Order        string
	Sum          float32
	Processed_at time.Time
}
