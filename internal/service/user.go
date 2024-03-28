package service

import (
	"Learn/LittleRedBook/internal/domain"
	"Learn/LittleRedBook/internal/repository"
	"context"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	// 加密

	// 存储
	return svc.repo.Create(ctx, u)
}
