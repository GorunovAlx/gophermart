package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	//"runtime"
	"sort"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/theplant/luhn"
	"github.com/urfave/negroni"
	"golang.org/x/crypto/bcrypt"

	//"golang.org/x/sync/errgroup"

	"github.com/GorunovAlx/gophermart/internal/gophermart/config"
	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/order"
	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/user"
	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/withdraw"
	accrualService "github.com/GorunovAlx/gophermart/internal/gophermart/services/accrual"
	loyaltyService "github.com/GorunovAlx/gophermart/internal/gophermart/services/loyalty"
	orderService "github.com/GorunovAlx/gophermart/internal/gophermart/services/order"
	userService "github.com/GorunovAlx/gophermart/internal/gophermart/services/user"
	withdrawService "github.com/GorunovAlx/gophermart/internal/gophermart/services/withdraw"
)

type Handler struct {
	Negroni       *negroni.Negroni
	Router        *mux.Router
	LoyaltySystem *loyaltyService.LoyaltySystem
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var (
	jwtKey         = []byte("my_secret_key")
	expirationTime = time.Now().Add(10 * time.Minute)
	statusNewValue = "NEW"
)

const (
	contextLogin contextKey = iota
	//cookieDuration                  = 5 * time.Minute
	//refreshTimeForCookie            = 60 * time.Second
)

func NewHandler() *Handler {
	us, err := userService.NewUserService(
		userService.WithMemoryUserRepository(),
	)
	if err != nil {
		log.Printf("NewHandler, NewUserService: %v", err.Error())
	}

	os, err := orderService.NewOrderService(
		orderService.WithMemoryOrderRepository(),
	)
	if err != nil {
		log.Printf("NewHandler, NewOrderService: %v", err.Error())
	}
	ws, err := withdrawService.NewWithdrawService(
		withdrawService.WithMemoryWithdrawRepository(),
	)
	if err != nil {
		log.Printf("NewHandler, NewWithdrawService: %v", err.Error())
	}

	as := accrualService.NewAccrualService(config.Cfg.AccrualAddress)

	ls, err := loyaltyService.NewLoyaltySystem(
		loyaltyService.WithUserService(us),
		loyaltyService.WithOrderService(os),
		loyaltyService.WithAccrualService(as),
		loyaltyService.WithWithdrawService(ws),
	)
	if err != nil {
		log.Printf("NewHandler, NewLoyaltySystem: %v", err.Error())
	}

	r := mux.NewRouter()

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseFunc(AuthMiddleware(us))
	n.UseFunc(UpdateOrdersMiddleware(ls))
	n.UseHandler(r)

	h := &Handler{
		Negroni:       n,
		Router:        r,
		LoyaltySystem: ls,
	}

	return h
}

func Initialize() *Handler {
	h := NewHandler()
	h.initializeRoutes()

	return h
}

func (h *Handler) initializeRoutes() {
	s := h.Router.PathPrefix("/api/user").Subrouter()
	s.HandleFunc("/register", h.registerUserHandler).Methods("POST")
	s.HandleFunc("/login", h.loginUserHandler).Methods("POST")
	s.HandleFunc("/orders", h.registerOrderHandler).Methods("POST")
	s.HandleFunc("/orders", h.getOrdersHandler).Methods("GET")
	s.HandleFunc("/balance", h.getCurrentBalance).Methods("GET")
	s.HandleFunc("/balance/withdraw", h.registerWithdraw).Methods("POST")
	s.HandleFunc("/balance/withdrawals", h.getWithdrawals).Methods("GET")
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

	err = h.LoyaltySystem.UserService.SetAuthToken(login, tokenString)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
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
	err = h.LoyaltySystem.UserService.AddUser(u.Login, string(hashedPassword))
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

	userValue := h.LoyaltySystem.UserService.GetUserByLogin(u.Login)
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

	userID, err := getUserID(h.LoyaltySystem, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.LoyaltySystem.OrderService.RegisterOrder(string(b), userID)
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

		for i := 0; i < 3; i++ {
			g.Go(func() error {
				for job := range jobCh {
					if err = h.LoyaltySystem.UpdateOrder(job); err != nil {
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

		//go func() {
		if err := g.Wait(); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//}()


			if err = h.LoyaltySystem.Update(string(b), oID, userID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
	*/
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(h.LoyaltySystem, r)
	log.Printf("user: %v", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := h.LoyaltySystem.OrderService.GetOrdersByUserID(userID)
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
	userID, err := getUserID(h.LoyaltySystem, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u, err := h.LoyaltySystem.UserService.GetUser(userID)
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

	userID, err := getUserID(h.LoyaltySystem, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.LoyaltySystem.RegisterWithdraw(b.OrderNumber, b.Sum, userID)
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
	userID, err := getUserID(h.LoyaltySystem, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := h.LoyaltySystem.WithdrawService.GetWithdrawals(userID)
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
