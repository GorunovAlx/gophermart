package v1

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type Handler struct {
	Negroni  *negroni.Negroni
	Router   *mux.Router
	Services *ServiceShelf
}

func NewHandler(s *ServiceShelf) *Handler {
	r := mux.NewRouter()
	n := negroni.New()

	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseFunc(AuthMiddleware(s.Users))
	n.UseFunc(UpdateOrdersMiddleware(s.Loyalty))
	n.UseHandler(r)

	h := &Handler{
		Negroni:  n,
		Router:   r,
		Services: s,
	}

	return h
}

func Initialize(s *ServiceShelf) *Handler {
	h := NewHandler(s)
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
