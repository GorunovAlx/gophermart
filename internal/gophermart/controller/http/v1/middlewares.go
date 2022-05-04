package v1

import (
	"context"
	"errors"
	"net/http"

	"time"

	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/user"
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
	jwtKey         = []byte(_secretKey)
	expirationTime = time.Now().Add(10 * time.Minute)
)

const (
	_registerPath                    = "/api/user/register"
	_loginPath                       = "/api/user/login"
	_contextToken         contextKey = iota
	_cookieDuration                  = 5 * time.Minute
	_refreshTimeForCookie            = 60 * time.Second
	_secretKey                       = "my_secret_key"
	_tokenString                     = "token"
	_fullPath                        = "/"
	_loginNotFound                   = "login not found"
)

func AuthMiddleware(us *userService.UserService) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		if r.RequestURI == _registerPath || r.RequestURI == _loginPath {
			next.ServeHTTP(w, r)
			return
		}
		c, err := r.Cookie(_tokenString)
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

		if time.Until(claims.ExpiresAt.Time) < _refreshTimeForCookie {
			expirationTime := time.Now().Add(_cookieDuration)
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
					w.Write([]byte(_loginNotFound))
					return
				}
			}

			http.SetCookie(w, &http.Cookie{
				Name:    _tokenString,
				Value:   tknStr,
				Path:    _fullPath,
				Expires: expirationTime,
			})
		}

		newCtx := context.WithValue(r.Context(), _contextToken, tknStr)
		next.ServeHTTP(w, r.WithContext(newCtx))
	}
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
