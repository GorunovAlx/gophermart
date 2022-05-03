package memory

import (
	"fmt"
	"sync"

	"github.com/GorunovAlx/gophermart/internal/gophermart/domain/user"
)

type MemoryUserRepository struct {
	users map[int]user.User
	sync.Mutex
}

func New() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: make(map[int]user.User),
	}
}

func (mr *MemoryUserRepository) Get(id int) (*user.User, error) {
	if user, ok := mr.users[id]; ok {
		return &user, nil
	}

	return &user.User{}, user.ErrUserNotFound
}

func (mr *MemoryUserRepository) SetAuthToken(login, token string) error {
	for _, user := range mr.users {
		if user.GetLogin() == login {
			user.SetAuthToken(token)
			return nil
		}
	}

	return user.ErrUserNotFound
}

func (mr *MemoryUserRepository) ChangeCurrentBalance(userID int, current float32) error {
	for _, user := range mr.users {
		if user.GetID() == userID {
			user.ChangeCurrentBalance(current)
			return nil
		}
	}

	return user.ErrUserNotFound
}

func (mr *MemoryUserRepository) ChangeWithdrawnBalance(userID int, withdraw float32) error {
	for _, user := range mr.users {
		if user.GetID() == userID {
			user.ChangeWithdrawnBalance(withdraw)
			return nil
		}
	}

	return user.ErrUserNotFound
}

func (mr *MemoryUserRepository) GetUserByLogin(login string) user.User {
	for _, user := range mr.users {
		if user.GetLogin() == login {
			return user
		}
	}

	return user.User{}
}

func (mr *MemoryUserRepository) GetIDByToken(token string) int {
	for _, user := range mr.users {
		if user.GetToken() == token {
			return user.GetID()
		}
	}

	return -1
}

func (mr *MemoryUserRepository) Add(u user.User) error {
	if mr.users == nil {
		mr.Lock()
		mr.users = make(map[int]user.User)
		mr.Unlock()
	}
	userID := getNextID(mr)
	u.SetID(userID)

	for _, userValue := range mr.users {
		if userValue.GetLogin() == u.GetLogin() {
			return fmt.Errorf("login is already taken: %w", user.ErrFailedToAddUser)
		}
	}

	mr.Lock()
	mr.users[userID] = u
	mr.Unlock()
	return nil
}

func getNextID(mr *MemoryUserRepository) int {
	var idCount = 0
	for range mr.users {
		idCount++
	}

	return idCount + 1
}
