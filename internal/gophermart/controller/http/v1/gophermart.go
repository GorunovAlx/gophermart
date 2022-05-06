package v1

import (
	//"context"
	"encoding/json"
	"errors"

	"io/ioutil"

	//"log"
	"net/http"

	//"runtime"

	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/theplant/luhn"
	"golang.org/x/crypto/bcrypt"

	//"golang.org/x/sync/errgroup"

	"github.com/GorunovAlx/gophermart/internal/gophermart/entity"
	"github.com/jackc/pgx/v4"
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

	user, err := h.Users.GetUserByLogin(u.Login)
	if user.ID != 0 {
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil && err != pgx.ErrNoRows {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.Users.Add(u.Login, string(hashedPassword))
	if err != nil {
		if errors.Is(err, entity.ErrFailedToAddUser) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	userIDToken, err := GenerateUserIDToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	expiration := time.Now().Add(cookieDuration)
	cookie := http.Cookie{
		Name:    "token",
		Value:   userIDToken,
		Path:    "/",
		Expires: expiration,
	}

	err = h.Users.SetAuthToken(u.Login, userIDToken)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
	}

	http.SetCookie(w, &cookie)

	//newCtx := context.WithValue(r.Context(), contextLogin, u.Login)
	//h.setToken(w, r.WithContext(newCtx))
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var u UserRequest
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.Users.GetUserByLogin(u.Login)
	if err == pgx.ErrNoRows {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("wrong credentials"))
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("wrong credentials"))
		return
	}

	userIDToken, err := GenerateUserIDToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	expiration := time.Now().Add(cookieDuration)
	cookie := http.Cookie{
		Name:    "token",
		Value:   userIDToken,
		Path:    "/",
		Expires: expiration,
	}

	err = h.Users.SetAuthToken(u.Login, userIDToken)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
	}

	http.SetCookie(w, &cookie)

	//newCtx := context.WithValue(r.Context(), contextLogin, u.Login)
	//h.setToken(w, r.WithContext(newCtx))
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) setToken(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value(contextLogin).(string)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.Users.SetAuthToken(login, tokenString)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Path:    "/",
		Expires: expirationTime,
	})
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

	userID := h.GetUserID(r)
	order, err := h.Orders.GetOrderByNumber(string(b))
	if err != nil && err != pgx.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if order.ID != 0 && order.UserID == userID {
		w.WriteHeader(http.StatusOK)
		return
	}
	if order.ID != 0 && order.UserID != userID {
		w.WriteHeader(http.StatusConflict)
		return
	}

	orderID, err := h.Orders.Add(userID, 0, statusNewValue, string(b))
	if orderID == -1 || err != nil {
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
	userID := h.GetUserID(r)
	res, err := h.Orders.GetOrders(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(res) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var orders []OrderResponse
	for _, order := range res {
		or := OrderResponse{
			Number:     order.Number,
			Status:     order.Status,
			Accrual:    order.Accrual,
			UploadedAt: order.UploadedAt.Format(time.RFC3339),
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
	userID := h.GetUserID(r)
	current, err := h.Users.GetBalance(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	withdrawn, err := h.Withdrawals.GetUserSumWithdrawn(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	balance := BalanceResponse{
		Current:   current,
		Withdrawn: withdrawn,
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
	var withdraw WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&withdraw); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderNumber, err := strconv.Atoi(withdraw.OrderNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !luhn.Valid(orderNumber) {
		http.Error(w, "incorrect order number format", http.StatusUnprocessableEntity)
		return
	}

	userID := h.GetUserID(r)
	current, err := h.Users.GetBalance(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if current < withdraw.Sum {
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}

	err = h.Withdrawals.Add(userID, withdraw.OrderNumber, withdraw.Sum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getWithdrawals(w http.ResponseWriter, r *http.Request) {
	userID := h.GetUserID(r)
	res, err := h.Withdrawals.GetWithdrawals(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(res) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var withdrawals []WithdrawResponse
	for _, withdraw := range res {
		wr := WithdrawResponse{
			Order:       withdraw.Order,
			Sum:         withdraw.Sum,
			ProcessedAt: withdraw.ProcessedAt.Format(time.RFC3339),
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

func (h *Handler) GetUserID(r *http.Request) int {
	contextToken := r.Context().Value(contextUserID)
	if contextToken == nil {
		return -1
	}
	token := contextToken.(string)
	return h.Users.GetIDByToken(token)
}
