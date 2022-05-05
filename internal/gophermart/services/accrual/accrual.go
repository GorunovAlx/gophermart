package accrual

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrDataRetrievalError = errors.New("an error occurred while receiving data")
	ErrStatusNotOk        = errors.New("status isn't ok")
	ErrDataProcessing     = errors.New("error occurred while processing the data")
)

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

func (as *AccrualService) GetAccrualOrder(number string) (AccrualOrder, error) {
	res, err := http.Get(fmt.Sprintf("%s/api/orders/%s", as.Address, number))
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

	return order, nil
}
