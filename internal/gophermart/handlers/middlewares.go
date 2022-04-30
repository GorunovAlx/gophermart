package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"

	"time"

	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/user"
	loyaltyService "github.com/GorunovAlx/gophermart/internal/gophermart/services/loyalty"
	userService "github.com/GorunovAlx/gophermart/internal/gophermart/services/user"
	"github.com/golang-jwt/jwt/v4"
	"github.com/urfave/negroni"
)

type contextKey int

const (
	registerPath                    = "/api/user/register"
	loginPath                       = "/api/user/login"
	contextToken         contextKey = iota
	cookieDuration                  = 5 * time.Minute
	refreshTimeForCookie            = 60 * time.Second
)

func AuthMiddleware(us *userService.UserService) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		if r.RequestURI == registerPath || r.RequestURI == loginPath {
			next.ServeHTTP(w, r)
			return
		}
		c, err := r.Cookie("token")
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

		if time.Until(claims.ExpiresAt.Time) < refreshTimeForCookie {
			expirationTime := time.Now().Add(cookieDuration)
			claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tknStr, err = token.SignedString(jwtKey)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = us.SetAuthToken(claims.Username, tknStr)
			if err != nil {
				if errors.Is(err, user.ErrUserNotFound) {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("login not found"))
					return
				}
			}

			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   tknStr,
				Path:    "/",
				Expires: expirationTime,
			})
		}

		newCtx := context.WithValue(r.Context(), contextToken, tknStr)
		next.ServeHTTP(w, r.WithContext(newCtx))
	}
}

func UpdateOrdersMiddleware(ls *loyaltyService.LoyaltySystem) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if r.RequestURI == registerPath || r.RequestURI == loginPath {
			next.ServeHTTP(w, r)
			return
		}

		userID, err := getUserID(ls, r)
		log.Printf("UpdateOrdersMiddleware. user: %v", userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		orders, err := ls.OrderService.GetOrdersNotProcessed(userID)
		log.Printf("UpdateOrdersMiddleware. len(orders): %v", len(orders))
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
				log.Printf("UpdateOrdersMiddleware. Update:order: %v", order)
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
