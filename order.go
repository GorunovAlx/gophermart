package gophermart

import "time"

type Order struct {
	ID         int
	UserID     int
	Number     string
	Accrual    float32
	Status     string
	UploadedAt time.Time
}
