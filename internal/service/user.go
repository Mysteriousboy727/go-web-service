package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"go-industry-server/internal/model"
	"go-industry-server/internal/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, req model.CreateUserRequest) (*model.User, error)
	GetUser(ctx context.Context, id string) (*model.User, error)
	ListUsers(ctx context.Context) ([]*model.User, error)
	UpdateUser(ctx context.Context, id string, req model.UpdateUserRequest) (*model.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type userService struct {
	repo   repository.UserRepository
	logger *slog.Logger
}

func NewUserService(repo repository.UserRepository, logger *slog.Logger) UserService {
	return &userService{repo: repo, logger: logger}
}

func (s *userService) CreateUser(ctx context.Context, req model.CreateUserRequest) (*model.User, error) {
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("name, email and password are required")
	}

	_, err := s.repo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, fmt.Errorf("email already in use")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("checking email uniqueness: %w", err)
	}

	now := time.Now().UTC()
	user := &model.User{
		ID:        newID(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  req.Password,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}
	return user, nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting user %s: %w", id, err)
	}
	return user, nil
}

func (s *userService) ListUsers(ctx context.Context) ([]*model.User, error) {
	return s.repo.List(ctx)
}

func (s *userService) UpdateUser(ctx context.Context, id string, req model.UpdateUserRequest) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user %s not found: %w", id, err)
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	user.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("updating user %s: %w", id, err)
	}
	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}
