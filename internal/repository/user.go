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

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByID(ctx context.Context, id int64) (domain.User, error)
	FindByWechat(ctx context.Context, openID string) (domain.User, error)
}

type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewCachedUserRepository(dao dao.UserDAO, c cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: c,
	}
}

func (repo *CachedUserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		ID:        u.ID,
		Email:     u.Email.String,
		Phone:     u.Phone.String,
		Password:  u.Password,
		CreatedAt: u.CreateTime,
		UpdatedAt: u.UpdateTime,
		WechatInfo: domain.WechatInfo{
			OpenID:  u.WechatOpenID.String,
			UnionID: u.WechatUnionID.String,
		},
	}
}

func (r *CachedUserRepository) toEntity(u domain.User) dao.User {
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
		WechatOpenID: sql.NullString{
			String: u.WechatInfo.OpenID,
			Valid:  u.WechatInfo.OpenID != "",
		},
		WechatUnionID: sql.NullString{
			String: u.WechatInfo.UnionID,
			Valid:  u.WechatInfo.UnionID != "",
		},
	}
}

func (r *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.toEntity(u))
}

func (r *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
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

func (r *CachedUserRepository) FindByID(ctx context.Context, id int64) (domain.User, error) {
	// 1.先查缓存
	user, err := r.cache.Get(ctx, id)
	if err == nil {
		// 缓存命中
		return user, nil
	}
	// 2.再查DB
	u, err := r.dao.FindByID(ctx, id)
	if err == ErrUserNotFound {
		return domain.User{}, ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, err
	}

	user = r.toDomain(u)
	//go func() {
	if err := r.cache.Set(ctx, user); err != nil {
		//日志，监控
	}
	return user, err
}

func (repo *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *CachedUserRepository) FindByWechat(ctx context.Context, openID string) (domain.User, error) {
	du, err := repo.dao.FindByWechat(ctx, openID)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(du), nil
}
