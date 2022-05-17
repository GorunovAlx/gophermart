package user

import (
	"errors"

	"github.com/GorunovAlx/gophermart/internal/gophermart/entity"
)

var (
	// ErrInvalidPerson is returned when the person is not valid in the NewUser factory
	ErrInvalidPerson = errors.New("a user has to have an valid data")
)

type User struct {
	person  *entity.Person
	balance *entity.Balance
}

func NewUser(login, pass string) (User, error) {
	if login == "" || pass == "" {
		return User{}, ErrInvalidPerson
	}

	person := &entity.Person{
		Login:    login,
		Password: pass,
	}

	b := &entity.Balance{
		Current:   0,
		Withdrawn: 0,
	}

	return User{
		person:  person,
		balance: b,
	}, nil
}

func (u *User) GetID() int {
	return u.person.ID
}

func (u *User) SetID(id int) {
	if u.person == nil {
		u.person = &entity.Person{}
	}
	u.person.ID = id
}

func (u *User) GetLogin() string {
	return u.person.Login
}

func (u *User) SetLogin(login string) {
	if u.person == nil {
		u.person = &entity.Person{}
	}
	u.person.Login = login
}

func (u *User) GetPassword() string {
	return u.person.Password
}

func (u *User) SetPassword(pass string) {
	if u.person == nil {
		u.person = &entity.Person{}
	}
	u.person.Password = pass
}

func (u *User) GetToken() string {
	return u.person.AuthToken
}

func (u *User) SetAuthToken(token string) {
	if u.person == nil {
		u.person = &entity.Person{}
	}
	u.person.AuthToken = token
}

func (u *User) GetCurrentBalance() float32 {
	return u.balance.Current
}

func (u *User) SetCurrentBalance(current float32) {
	if u.person == nil {
		u.person = &entity.Person{}
	}
	u.balance.Current = current
}

func (u *User) SetWithdrawnBalance(withdrawn float32) {
	if u.person == nil {
		u.person = &entity.Person{}
	}
	u.balance.Withdrawn = withdrawn
}

func (u *User) ChangeCurrentBalance(current float32) {
	if u.person == nil {
		u.person = &entity.Person{}
	}
	u.balance.Current += current
}

func (u *User) GetWithdrawnBalance() float32 {
	return u.balance.Withdrawn
}

func (u *User) ChangeWithdrawnBalance(withdrawn float32) {
	if u.person == nil {
		u.person = &entity.Person{}
	}
	u.balance.Withdrawn += withdrawn
}
