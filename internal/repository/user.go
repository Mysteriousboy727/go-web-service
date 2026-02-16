package repository

import (
	"context"
	"errors"
	"sync"

	"go-industry-server/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	List(ctx context.Context) ([]*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
}

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*model.User
}

func NewInMemoryUserRepository() UserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*model.User),
	}
}

func (r *InMemoryUserRepository) Create(_ context.Context, user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return ErrAlreadyExists
	}
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) GetByID(_ context.Context, id string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}

func (r *InMemoryUserRepository) GetByEmail(_ context.Context, email string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, ErrNotFound
}

func (r *InMemoryUserRepository) List(_ context.Context) ([]*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*model.User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, u)
	}
	return users, nil
}

func (r *InMemoryUserRepository) Update(_ context.Context, user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[user.ID]; !ok {
		return ErrNotFound
	}
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[id]; !ok {
		return ErrNotFound
	}
	delete(r.users, id)
	return nil
}
