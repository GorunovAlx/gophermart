package entity

import (
	"context"
	"errors"

	"github.com/GorunovAlx/gophermart/internal/gophermart/database"
)

type User struct {
	ID        int
	Login     string
	Password  string
	AuthToken string
}

type UserStorage struct {
	S database.Storage
}

type UserRepository interface {
	Get(id int) (User, error)
	GetIDByToken(token string) int
	SetAuthToken(login, token string) error
	GetUserByLogin(login string) (User, error)
	Add(login, password string) error
	GetBalance(userID int) (float32, error)
}

var (
	ErrUserNotFound    = errors.New("the user was not found in the repository")
	ErrFailedToAddUser = errors.New("failed to add the user to the repository")
	ErrNotEnoughFunds  = errors.New("not enough funds on the user's balance")
)

func (us UserStorage) Get(id int) (User, error) {
	var u User
	err := us.S.PGpool.QueryRow(
		context.Background(),
		"select id, login, password, authtoken, current, withdrawn from users where id=$1",
		id,
	).Scan(
		&u.ID,
		&u.Login,
		&u.Password,
		&u.AuthToken,
	)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func (us UserStorage) GetIDByToken(token string) int {
	var userID int
	err := us.S.PGpool.QueryRow(
		context.Background(),
		"select id from users where authtoken=$1",
		token,
	).Scan(&userID)
	if err != nil {
		return -1
	}

	return userID
}

func (us UserStorage) SetAuthToken(login, token string) error {
	sqlStmt := `
	update users set authtoken = $1 where login = $2;`

	_, err := us.S.PGpool.Exec(
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

func (us UserStorage) GetUserByLogin(login string) (User, error) {
	var u User
	err := us.S.PGpool.QueryRow(
		context.Background(),
		"select id, login, password, authtoken from users where login=$1",
		login,
	).Scan(
		&u.ID,
		&u.Login,
		&u.Password,
		&u.AuthToken,
	)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func (us UserStorage) Add(login, password string) error {
	insertStatement := `
	INSERT INTO users (login, password)
	VALUES ($1, $2);`

	commandTag, err := us.S.PGpool.Exec(
		context.Background(),
		insertStatement,
		login,
		password,
	)

	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("unsuccess add")
	}

	return nil
}

func (us UserStorage) GetBalance(userID int) (float32, error) {
	var current float32
	var withdrawn float32
	err := us.S.PGpool.QueryRow(
		context.Background(),
		"select coalesce(sum(orders.accrual), 0) from orders where orders.user_id = $1 and orders.status = 'PROCESSED';",
		userID,
	).Scan(&current)
	if err != nil {
		return 0, err
	}

	err = us.S.PGpool.QueryRow(
		context.Background(),
		"select coalesce(sum(withdrawals.sum), 0) from withdrawals where withdrawals.user_id = $1;",
		userID,
	).Scan(&withdrawn)
	if err != nil {
		return 0, err
	}

	return current - withdrawn, nil
}
