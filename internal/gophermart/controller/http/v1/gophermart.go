package v1

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"

	//"log"
	"net/http"

	//"runtime"
	"sort"
	"strconv"
	"time"

	"github.com/theplant/luhn"
	"golang.org/x/crypto/bcrypt"

	//"golang.org/x/sync/errgroup"

	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/order"
	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/user"
	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/withdraw"
	loyaltyService "github.com/GorunovAlx/gophermart/internal/gophermart/services/loyalty"
)

var (
	statusNewValue = "NEW"
)

const (
	contextLogin contextKey = iota
)

func (h *Handler) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var u UserRequest
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = h.Services.Users.AddUser(u.Login, string(hashedPassword))
	if err != nil {
		if errors.Is(err, user.ErrFailedToAddUser) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	newCtx := context.WithValue(r.Context(), contextLogin, u.Login)
	h.setToken(w, r.WithContext(newCtx))
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var u UserRequest
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userValue := h.Services.Users.GetUserByLogin(u.Login)
	if (user.User{}) == userValue {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("wrong credentials"))
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userValue.GetPassword()), []byte(u.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("wrong credentials"))
		return
	}

	newCtx := context.WithValue(r.Context(), contextLogin, u.Login)
	h.setToken(w, r.WithContext(newCtx))
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) registerOrderHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderNumber, err := strconv.Atoi(string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !luhn.Valid(orderNumber) {
		http.Error(w, "incorrect order number format", http.StatusUnprocessableEntity)
		return
	}

	userID, err := getUserID(h.Services.Loyalty, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.Services.Orders.RegisterOrder(string(b), userID)
	if err != nil {
		if errors.Is(err, order.ErrOrderAlreadyRegisteredByUser) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(err.Error()))
			return
		}
		if errors.Is(err, order.ErrOrderAlreadyRegisteredByOtherUser) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(err.Error()))
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	/*
		jobCh := make(chan *order.OrderJob)
		g, _ := errgroup.WithContext(context.Background())

		for i := 0; i < runtime.NumCPU(); i++ {
			g.Go(func() error {
				for job := range jobCh {
					if err = h.Services.Loyalty.UpdateOrder(job); err != nil {
						return err
					}
				}
				return nil
			})
		}

		job := &order.OrderJob{
			Number: string(b),
			ID:     oID,
			UserID: userID,
			Status: statusNewValue,
		}
		jobCh <- job

		go func() {
			if err := g.Wait(); err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}()
	*/

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(h.Services.Loyalty, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := h.Services.Orders.GetOrdersByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(res) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	sort.Sort(order.ByUploadedAt(res))

	var orders []OrderResponse
	for _, order := range res {
		or := OrderResponse{
			Number:     order.GetNumber(),
			Status:     order.GetStatus(),
			Accrual:    order.GetAccrual(),
			UploadedAt: order.GetUploadedAt().Format(time.RFC3339),
		}
		orders = append(orders, or)
	}

	resp, err := json.MarshalIndent(orders, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *Handler) getCurrentBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(h.Services.Loyalty, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u, err := h.Services.Users.GetUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	balance := BalanceResponse{
		Current:   u.GetCurrentBalance(),
		Withdrawn: u.GetWithdrawnBalance(),
	}
	resp, err := json.MarshalIndent(balance, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *Handler) registerWithdraw(w http.ResponseWriter, r *http.Request) {
	var b BalanceRequest
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderNumber, err := strconv.Atoi(string(b.OrderNumber))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !luhn.Valid(orderNumber) {
		http.Error(w, "incorrect order number format", http.StatusUnprocessableEntity)
		return
	}

	userID, err := getUserID(h.Services.Loyalty, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.Services.Loyalty.RegisterWithdraw(b.OrderNumber, b.Sum, userID)
	if err != nil {
		if errors.Is(err, user.ErrNotEnoughFunds) {
			w.WriteHeader(http.StatusPaymentRequired)
			w.Write([]byte(err.Error()))
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getWithdrawals(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(h.Services.Loyalty, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := h.Services.Loyalty.WithdrawService.GetWithdrawals(userID)
	if len(res) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	sort.Sort(withdraw.ByUploadedAt(res))

	var withdrawals []WithdrawResponse
	for _, withdraw := range res {
		wr := WithdrawResponse{
			Order:       withdraw.GetOrder(),
			Sum:         withdraw.GetSum(),
			ProcessedAt: withdraw.GetProcessedAt().Format(time.RFC3339),
		}
		withdrawals = append(withdrawals, wr)
	}

	resp, err := json.MarshalIndent(withdrawals, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func getUserID(ls *loyaltyService.LoyaltySystem, r *http.Request) (int, error) {
	token := r.Context().Value(contextToken).(string)
	return ls.UserService.GetUserIDByToken(token)
}
