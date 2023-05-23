package repository

import (
	"sync"

	"github.com/Totus-Floreo/asperitas-on-go/pkg/model"
	uuid "github.com/google/uuid"
)

type UserStorage struct {
	Storage map[string]*model.User
	mu      *sync.Mutex
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
		Storage: make(map[string]*model.User),
		mu:      new(sync.Mutex),
	}
}

func (s *UserStorage) GetUser(userID string) (*model.User, error) {
	for _, user := range s.Storage {
		if user.ID == userID {
			return user, nil
		}
	}

	return nil, model.ErrUserNotFound
}

func (s *UserStorage) AddUser(user *model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user.ID = uuid.New().String()
	s.Storage[user.Username] = user
	return nil
}
