package accrual

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/GorunovAlx/gophermart/internal/gophermart/entity"
)

var (
	ErrDataRetrievalError = errors.New("an error occurred while receiving data")
	ErrStatusNotOk        = errors.New("status isn't ok")
	ErrDataProcessing     = errors.New("error occurred while processing the data")
)

type AccrualOrder struct {
	Number  string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual,omitempty"`
}

type AccrualService struct {
	Address string
	Os      entity.OrderRepository
}

func NewAccrualService(address string, o entity.OrderRepository) *AccrualService {
	return &AccrualService{
		Address: address,
		Os:      o,
	}
}

func (as *AccrualService) GetAccrualOrder(number string) (AccrualOrder, error) {
	parsedURL, err := url.Parse(fmt.Sprintf("%s/api/orders/%s", as.Address, number))
	if err != nil {
		return AccrualOrder{}, err
	}
	res, err := http.Get(parsedURL.String())
	if err != nil {
		return AccrualOrder{}, ErrDataRetrievalError
	}
	if res.StatusCode != http.StatusOK {
		return AccrualOrder{}, nil
	}

	var order AccrualOrder
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&order); err != nil {
		return AccrualOrder{}, ErrDataProcessing
	}

	return order, nil
}

func (as *AccrualService) UpdateAccrualOrder(number string) error {
	order, err := as.GetAccrualOrder(number)
	if err != nil {
		return err
	}

	err = as.Os.Update(order.Status, order.Accrual, order.Number)
	if err != nil {
		return err
	}

	return nil
}
