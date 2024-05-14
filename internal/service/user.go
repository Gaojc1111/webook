package service

import (
	"context"
	"errors"
	"webook/internal/domain"
	"webook/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

// var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
var ErrUserDuplicated = repository.ErrUserDuplicated
var ErrInvalidUserOrPassword = errors.New("邮箱或密码错误")

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email, password string) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) SignUp(ctx context.Context, u domain.User) error {
	// 加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	// 存储
	return svc.repo.Create(ctx, u)
}

func (svc *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
	// 先查询有没有这个用户
	user, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		// todo
		return user, ErrInvalidUserOrPassword
	}
	if err != nil {
		return user, err
	}
	// 比较密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// todo 打印日志
		return user, ErrInvalidUserOrPassword
	}
	// 没问题
	return user, nil
}

func (svc *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
	user, err := svc.repo.FindByID(ctx, id)
	return user, err
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	// 先找一下，大部分用户是已经存在的用户
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		// err == nil, 找到User，直接返回
		// err != nil，系统错误，直接返回
		return u, err
	}
	// 用户没找到，注册
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	// 有两种可能，一种是 err 恰好是唯一索引冲突（phone）
	// 一种是 err != nil，系统错误
	if err != nil && err != ErrUserDuplicated {
		return domain.User{}, err
	}
	// 要么 err ==nil，要么ErrDuplicateUser，也代表用户存在
	// 主从延迟，理论上来讲，强制走主库
	return svc.repo.FindByPhone(ctx, phone)
}
