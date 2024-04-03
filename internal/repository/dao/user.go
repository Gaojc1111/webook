package dao

import (
	"context"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"

	"gorm.io/gorm"
)

var (
	ErrUserDuplicateEmail = errors.New("该邮箱已被注册")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

// User 直接对应数据库表
// entity/model
type User struct {
	ID         int64  `gorm:"primaryKey,autoIncrement"`
	Email      string `gorm:"unique"`
	Password   string
	CreateTime int64 `gorm:"column:createTime"`
	UpdateTime int64 `gorm:"column:updateTime"`
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.CreateTime = now
	u.UpdateTime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			// 邮箱冲突（唯一键）
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	user := User{}
	// 查找email = email 的第一条记录
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&user).Error

	return user, err
}
