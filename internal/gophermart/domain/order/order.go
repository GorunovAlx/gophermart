package order

import (
	"time"

	"github.com/GorunovAlx/gophermart/internal/gophermart/entity"
)

type Order struct {
	order *entity.Order
}

type OrderJob struct {
	Number string
	ID     int
	UserID int
	Status string
}

type ByUploadedAt []Order

func NewOrder(number string, userID int) Order {
	o := &entity.Order{
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

func (o *Order) SetID(id int) {
	if o.order == nil {
		o.order = &entity.Order{}
	}
	o.order.ID = id
}

func (o *Order) GetUserID() int {
	return o.order.UserID
}

func (o *Order) SetUserID(userID int) {
	if o.order == nil {
		o.order = &entity.Order{}
	}
	o.order.UserID = userID
}

func (o *Order) GetNumber() string {
	return o.order.Number
}

func (o *Order) SetNumber(number string) {
	if o.order == nil {
		o.order = &entity.Order{}
	}
	o.order.Number = number
}

func (o *Order) GetStatus() string {
	return o.order.Status
}

func (o *Order) SetStatus(status string) {
	if o.order == nil {
		o.order = &entity.Order{}
	}
	o.order.Status = status
}

func (o *Order) GetAccrual() float32 {
	return o.order.Accrual
}

func (o *Order) SetAccrual(accrual float32) {
	if o.order == nil {
		o.order = &entity.Order{}
	}
	o.order.Accrual = accrual
}

func (o *Order) GetUploadedAt() time.Time {
	return o.order.UploadedAt
}

func (o *Order) SetUploadedAt(time time.Time) {
	if o.order == nil {
		o.order = &entity.Order{}
	}
	o.order.UploadedAt = time
}
