package v1

import (
	"context"
	"errors"
	"net/http"

	"time"

	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/user"
	loyaltyService "github.com/GorunovAlx/gophermart/internal/gophermart/services/loyalty"
	userService "github.com/GorunovAlx/gophermart/internal/gophermart/services/user"
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
	contextToken         contextKey = iota
	cookieDuration                  = 5 * time.Minute
	refreshTimeForCookie            = 60 * time.Second
	secretKey                       = "my_secret_key"
	tokenString                     = "token"
	fullPath                        = "/"
	loginNotFound                   = "login not found"
)

func AuthMiddleware(us *userService.UserService) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		if r.RequestURI == registerPath || r.RequestURI == loginPath {
			next.ServeHTTP(w, r)
			return
		}
		c, err := r.Cookie(tokenString)
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tknStr := c.Value
		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		newCtx := context.WithValue(r.Context(), contextToken, tknStr)
		next.ServeHTTP(w, r.WithContext(newCtx))
	}
}

func (h *Handler) setToken(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value(contextLogin).(string)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: login,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.Services.Users.SetAuthToken(login, tokenString)
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

func UpdateOrdersMiddleware(ls *loyaltyService.LoyaltySystem) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if r.RequestURI == registerPath || r.RequestURI == loginPath {
			next.ServeHTTP(w, r)
			return
		}

		userID, err := getUserID(ls, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		orders, err := ls.OrderService.GetOrdersNotProcessed(userID)
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
				err = ls.Update(order.GetNumber(), order.GetID(), userID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}()

		next.ServeHTTP(w, r)
	}
}
