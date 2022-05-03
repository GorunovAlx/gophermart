package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/GorunovAlx/gophermart"
	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/user"
)

type PostgresUserRepository struct {
	*pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		db,
	}
}

func (db *PostgresUserRepository) Get(id int) (*user.User, error) {
	conn, err := db.Acquire(context.Background())
	if err != nil {
		return &user.User{}, err
	}
	defer conn.Release()

	var person gophermart.Person
	var balance gophermart.Balance
	err = conn.QueryRow(
		context.Background(),
		"select id, login, password, authtoken, current, withdrawn from users where id=$1",
		id,
	).Scan(
		&person.ID,
		&person.Login,
		&person.Password,
		&person.AuthToken,
		&balance.Current,
		&balance.Withdrawn,
	)
	if err != nil {
		return &user.User{}, err
	}
	u, err := user.NewUser(person.Login, person.Password)
	if err != nil {
		return &user.User{}, err
	}
	u.SetID(person.ID)
	u.SetAuthToken(person.AuthToken)
	u.SetCurrentBalance(balance.Current)
	u.SetWithdrawnBalance(balance.Withdrawn)

	return &u, nil
}

func (db *PostgresUserRepository) ChangeCurrentBalance(userID int, accrual float32) error {
	conn, err := db.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	sqlStmt := `
	update users set current = current + $1 where id = $2;`

	_, err = conn.Exec(
		context.Background(),
		sqlStmt,
		accrual,
		userID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (db *PostgresUserRepository) ChangeWithdrawnBalance(userID int, withdraw float32) error {
	conn, err := db.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	sqlStmt := `
	update users set withdrawn = withdrawn + $1 where id = $2;`

	_, err = conn.Exec(
		context.Background(),
		sqlStmt,
		withdraw,
		userID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (db *PostgresUserRepository) GetIDByToken(token string) int {
	conn, err := db.Acquire(context.Background())
	if err != nil {
		return -1
	}
	defer conn.Release()

	var userID int
	err = conn.QueryRow(
		context.Background(),
		"select id from users where authtoken=$1",
		token,
	).Scan(&userID)
	if err != nil {
		return -1
	}

	return userID
}

func (db *PostgresUserRepository) SetAuthToken(login, token string) error {
	conn, err := db.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	sqlStmt := `
	update users set authtoken = $1 where login = $2;`

	_, err = conn.Exec(
		context.Background(),
		sqlStmt,
		token,
		login,
	)

	if err != nil {
		return err
	}

	return nil
}

func (db *PostgresUserRepository) GetUserByLogin(login string) user.User {
	conn, err := db.Acquire(context.Background())
	if err != nil {
		return user.User{}
	}
	defer conn.Release()

	var person gophermart.Person
	var balance gophermart.Balance
	err = conn.QueryRow(
		context.Background(),
		"select id, login, password, authtoken, current, withdrawn from users where login=$1",
		login,
	).Scan(
		&person.ID,
		&person.Login,
		&person.Password,
		&person.AuthToken,
		&balance.Current,
		&balance.Withdrawn,
	)
	if err != nil {
		return user.User{}
	}
	u, err := user.NewUser(person.Login, person.Password)
	if err != nil {
		return user.User{}
	}
	u.SetID(person.ID)
	u.SetAuthToken(person.AuthToken)
	u.SetCurrentBalance(balance.Current)
	u.SetWithdrawnBalance(balance.Withdrawn)

	return u
}

func (db *PostgresUserRepository) Add(u user.User) error {
	conn, err := db.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	insertStatement := `
	INSERT INTO users (login, password, authtoken, current, withdrawn)
	VALUES ($1, $2, $3, $4, $5);`

	commandTag, err := conn.Exec(
		context.Background(),
		insertStatement,
		u.GetLogin(),
		u.GetPassword(),
		u.GetToken(),
		u.GetCurrentBalance(),
		u.GetWithdrawnBalance(),
	)

	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("Unsuccess add")
	}

	return nil
}
