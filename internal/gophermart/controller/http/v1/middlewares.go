package v1

import (
	"context"
	"net/http"

	"time"

	"github.com/GorunovAlx/gophermart/internal/gophermart/entity"
	"github.com/golang-jwt/jwt/v4"
	"github.com/urfave/negroni"
)

type (
	Claims struct {
		Username string `json:"username"`
		jwt.RegisteredClaims
	}

	contextKey int
)

var (
	jwtKey         = []byte(secretKey)
	expirationTime = time.Now().Add(60 * time.Minute)
)

const (
	registerPath                    = "/api/user/register"
	loginPath                       = "/api/user/login"
	contextUserID        contextKey = iota
	cookieDuration                  = 65 * time.Minute
	refreshTimeForCookie            = 60 * time.Second
	secretKey                       = "my_secret_key"
	tokenString                     = "token"
	fullPath                        = "/"
	loginNotFound                   = "login not found"
)

func AuthMiddleware(us entity.UserRepository) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if r.RequestURI == registerPath || r.RequestURI == loginPath {
			next.ServeHTTP(w, r)
			return
		}

		userIDToken := getCookieByName("token", r)

		if len(userIDToken) != 0 {
			isAuthentic, err := AuthUserIDToken(userIDToken)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if isAuthentic {
				id := us.GetIDByToken(userIDToken)
				if id == -1 {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				ctx := r.Context()
				ctx = context.WithValue(ctx, contextUserID, userIDToken)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		next.ServeHTTP(w, r)
	}
}

func UpdateOrdersMiddleware(h *Handler) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if r.RequestURI == registerPath || r.RequestURI == loginPath {
			next.ServeHTTP(w, r)
			return
		}

		userID := h.GetUserID(r)
		if userID == -1 {
			http.Error(w, entity.ErrUserNotFound.Error(), http.StatusInternalServerError)
			return
		}

		orders, err := h.Orders.GetOrdersNotProcessed(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(orders) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		go func() {
			for _, order := range orders {
				accrualOrder, err := h.Accruals.GetAccrualOrder(order.Number)
				err = h.Orders.Update(accrualOrder, userID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}()

		next.ServeHTTP(w, r)
	}
}

func getCookieByName(cName string, r *http.Request) string {
	receivedCookie := r.Cookies()
	var value string
	if len(receivedCookie) != 0 {
		for _, cookie := range receivedCookie {
			if cookie.Name == cName {
				value = cookie.Value
			}
		}
	}

	return value
}
