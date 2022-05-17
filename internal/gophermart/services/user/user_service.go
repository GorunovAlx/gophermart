package services

import (
	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/user"
	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/user/memory"
	userDB "github.com/GorunovAlx/gophermart/internal/gophermart/domain/user/postgres"
	"github.com/jackc/pgx/v4"
)

// UserConfiguration is an alias for a function that will take in a pointer to an UserService and modify it
type UserConfiguration func(us *UserService) error

// UserService is a implementation of the UserService
type UserService struct {
	users user.UserRepository
}

// NewUserService takes a variable amount of UserConfiguration functions and returns a new UserService
// Each UserConfiguration will be called in the user they are passed in
func NewUserService(cfgs ...UserConfiguration) (*UserService, error) {
	// Create the userservice
	us := &UserService{}
	// Apply all Configurations passed in
	for _, cfg := range cfgs {
		// Pass the service into the configuration function
		err := cfg(us)
		if err != nil {
			return nil, err
		}
	}
	return us, nil
}

// WithUserRepository applies a given user repository to the UserService
func WithUserRepository(ur user.UserRepository) UserConfiguration {
	// return a function that matches the UserConfiguration alias,
	// You need to return this so that the parent function can take in all the needed parameters
	return func(us *UserService) error {
		us.users = ur
		return nil
	}
}

// WithMemoryUserRepository applies a memory user repository to the UserService
func WithMemoryUserRepository() UserConfiguration {
	// Create the memory repo, if we needed parameters, such as connection strings they could be inputted here
	ur := memory.New()
	return WithUserRepository(ur)
}

func WithPostgresUserRepository(pool *pgx.Conn) UserConfiguration {
	return func(us *UserService) error {
		pur := userDB.NewPostgresRepository(pool)
		us.users = pur
		return nil
	}
}

func (us *UserService) GetUser(userID int) (*user.User, error) {
	u, err := us.users.Get(userID)
	if err != nil {
		return &user.User{}, err
	}
	return u, nil
}

func (us *UserService) AddUser(login, pass string) error {
	u, err := user.NewUser(login, pass)
	if err != nil {
		return err
	}
	if err = us.users.Add(u); err != nil {
		return err
	}
	return nil
}

func (us *UserService) GetUserByLogin(login string) user.User {
	return us.users.GetUserByLogin(login)
}

func (us *UserService) SetAuthToken(login, token string) error {
	if err := us.users.SetAuthToken(login, token); err != nil {
		return err
	}
	return nil
}

func (us *UserService) GetUserIDByToken(token string) (int, error) {
	id := us.users.GetIDByToken(token)
	if id == -1 {
		return -1, user.ErrUserNotFound
	}
	return id, nil
}

func (us *UserService) ChangeBalance(userID int, accrual float32) error {
	err := us.users.ChangeCurrentBalance(userID, accrual)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) TakeOutSum(userID int, sum float32) error {
	err := us.users.ChangeCurrentBalance(userID, -sum)
	if err != nil {
		return err
	}
	err = us.users.ChangeWithdrawnBalance(userID, sum)
	if err != nil {
		return err
	}

	return nil
}
