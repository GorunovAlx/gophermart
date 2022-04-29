module github.com/GorunovAlx/gophermart

go 1.18

replace github.com/GorunovAlx/gophermart/internal/gophermart/entity => ../internal/gophermart/entity

require (
	github.com/caarlos0/env/v6 v6.9.1
	github.com/golang-jwt/jwt/v4 v4.4.1
	github.com/gorilla/mux v1.8.0
	github.com/theplant/luhn v0.0.0-20170224032821-81a1a381387a
	github.com/urfave/negroni v1.0.0
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)
