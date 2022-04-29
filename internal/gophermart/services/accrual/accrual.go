package accrual

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var (
	ErrDataRetrievalError = errors.New("an error occurred while receiving data")
	ErrStatusNotOk        = errors.New("status isn't ok")
	ErrDataProcessing     = errors.New("error occurred while processing the data")
)

var statuses = []string{"NEW", "PROCESSING", "INVALID", "PROCESSED"}

type AccrualOrder struct {
	OrderID int
	Number  string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual,omitempty"`
}

type AccrualService struct {
	Address string
}

func NewAccrualService(address string) *AccrualService {
	return &AccrualService{
		Address: address,
	}
}

func random() int {
	rand.Seed(time.Now().UnixNano())
	min := 0
	max := 3
	return rand.Intn(max-min+1) + min
}

func (as *AccrualService) GetAccrualOrder(number string) (AccrualOrder, error) {
	url := fmt.Sprintf("http://%v/api/orders/%v", as.Address, number)
	res, err := http.Get(url)
	if err != nil {
		return AccrualOrder{}, ErrDataRetrievalError
	}
	if res.StatusCode != http.StatusOK {
		return AccrualOrder{}, ErrStatusNotOk
	}

	var order AccrualOrder
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&order); err != nil {
		return AccrualOrder{}, ErrDataProcessing
	}
	r := random()
	order.Status = statuses[r]

	return order, nil
}
