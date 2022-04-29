package order

import (
	"time"

	"github.com/GorunovAlx/gophermart"
)

type Order struct {
	order *gophermart.Order
}

type OrderJob struct {
	Number string
	ID     int
	UserID int
	Status string
}

type ByUploadedAt []Order

func NewOrder(number string, userID int) Order {
	o := &gophermart.Order{
		Number: number,
		UserID: userID,
		Status: "NEW",
	}

	return Order{
		order: o,
	}
}

func (b ByUploadedAt) Len() int { return len(b) }

func (b ByUploadedAt) Less(i, j int) bool {
	return b[i].GetUploadedAt().Unix() < b[j].GetUploadedAt().Unix()
}

func (b ByUploadedAt) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

func (o *Order) GetID() int {
	return o.order.ID
}

func (u *Order) SetID(id int) {
	if u.order == nil {
		u.order = &gophermart.Order{}
	}
	u.order.ID = id
}

func (o *Order) GetUserID() int {
	return o.order.UserID
}

func (u *Order) SetUserID(userID int) {
	if u.order == nil {
		u.order = &gophermart.Order{}
	}
	u.order.UserID = userID
}

func (o *Order) GetNumber() string {
	return o.order.Number
}

func (u *Order) SetNumber(number string) {
	if u.order == nil {
		u.order = &gophermart.Order{}
	}
	u.order.Number = number
}

func (o *Order) GetStatus() string {
	return o.order.Status
}

func (u *Order) SetStatus(status string) {
	if u.order == nil {
		u.order = &gophermart.Order{}
	}
	u.order.Status = status
}

func (o *Order) GetAccrual() float32 {
	return o.order.Accrual
}

func (u *Order) SetAccrual(accrual float32) {
	if u.order == nil {
		u.order = &gophermart.Order{}
	}
	u.order.Accrual = accrual
}

func (o *Order) GetUploadedAt() time.Time {
	return o.order.Uploaded_at
}

func (u *Order) SetUploadedAt(time time.Time) {
	if u.order == nil {
		u.order = &gophermart.Order{}
	}
	u.order.Uploaded_at = time
}
