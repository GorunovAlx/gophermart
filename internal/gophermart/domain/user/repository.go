package user

import (
	"errors"
)

var (
	ErrUserNotFound    = errors.New("the user was not found in the repository")
	ErrFailedToAddUser = errors.New("failed to add the user to the repository")
	ErrNotEnoughFunds  = errors.New("not enough funds on the user's balance")
)

type UserRepository interface {
	Get(int) (*User, error)
	GetIDByToken(token string) int
	SetAuthToken(login, token string) error
	GetUserByLogin(login string) User
	Add(User) error
	ChangeCurrentBalance(userID int, current float32) error
	ChangeWithdrawnBalance(userID int, withdraw float32) error
}
