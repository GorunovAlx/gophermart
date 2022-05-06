module github.com/GorunovAlx/gophermart

go 1.18

replace github.com/GorunovAlx/gophermart/internal/gophermart/application/config => ../internal/gophermart/application/config

require (
	github.com/caarlos0/env/v6 v6.9.1
	github.com/golang-migrate/migrate/v4 v4.15.2
	github.com/gorilla/mux v1.8.0
	github.com/jackc/pgx/v4 v4.16.0
	github.com/rs/zerolog v1.26.1
	github.com/theplant/luhn v0.0.0-20170224032821-81a1a381387a
	github.com/urfave/negroni v1.0.0
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4
)

require (
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.12.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.11.0 // indirect
	github.com/lib/pq v1.10.2 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/text v0.3.7 // indirect
)
