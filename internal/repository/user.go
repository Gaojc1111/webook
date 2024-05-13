package repository

import (
	"context"
	"database/sql"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
)

var (
	ErrUserDuplicated = dao.ErrUserDuplicated
	ErrUserNotFound   = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: c,
	}
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		ID:       u.ID,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Password: u.Password,
	}
}

func (r *UserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		ID: u.ID,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Password: u.Password,
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.toEntity(u))
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID:       user.ID,
		Email:    user.Email.String,
		Password: user.Password,
	}, err
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (domain.User, error) {
	// 1.先查缓存
	user, err := r.cache.Get(ctx, id)
	if err != nil {
		return user, err
	}
	// 2.再查DB
	u, err := r.dao.FindByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	user = domain.User{
		ID:       u.ID,
		Email:    u.Email.String,
		Password: u.Password,
	}
	go func() {
		err = r.cache.Set(ctx, user)
		// 缓存崩了怎么办
		if err != nil {
			//日志，监控
		}
	}()
	return user, err
}

func (repo *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}
