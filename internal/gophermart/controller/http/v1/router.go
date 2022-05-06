package v1

import (
	"github.com/GorunovAlx/gophermart/internal/gophermart/accrual"
	"github.com/GorunovAlx/gophermart/internal/gophermart/entity"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type Handler struct {
	Negroni     *negroni.Negroni
	Router      *mux.Router
	Users       entity.UserRepository
	Orders      entity.OrderRepository
	Withdrawals entity.WithdrawRepository
	Accruals    *accrual.AccrualService
}

func NewHandler(u entity.UserRepository, o entity.OrderRepository, w entity.WithdrawRepository, a *accrual.AccrualService) *Handler {
	r := mux.NewRouter()
	n := negroni.New()

	h := &Handler{
		Negroni:     n,
		Router:      r,
		Users:       u,
		Orders:      o,
		Withdrawals: w,
		Accruals:    a,
	}

	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseFunc(AuthMiddleware(u))
	n.UseFunc(UpdateOrdersMiddleware(h))
	n.UseHandler(r)

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
