package dao

import (
	"context"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicated = gorm.ErrDuplicatedKey
	ErrUserNotFound   = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(ctx context.Context, u User) error
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByID(ctx context.Context, id int64) (User, error)
	FindByWechat(ctx context.Context, openID string) (User, error)
}

type GormUserDAO struct {
	db *gorm.DB
}

// User 直接对应数据库表
// entity/model
type User struct {
	ID            int64          `gorm:"primaryKey,autoIncrement"`
	Email         sql.NullString `gorm:"unique"`
	Password      string         `gorm:"column:Password"`
	Phone         sql.NullString `gorm:"unique"`
	CreateTime    int64          `gorm:"column:createTime"`
	UpdateTime    int64          `gorm:"column:updateTime"`
	WechatOpenID  sql.NullString `gorm:"column:wechatOpenID"`
	WechatUnionID sql.NullString `gorm:"column:wechatUnionID"`
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GormUserDAO{db: db}
}

func (dao *GormUserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.CreateTime = now
	u.UpdateTime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			// 邮箱冲突（唯一键）
			return ErrUserDuplicated
		}
	}
	return err
}

func (dao *GormUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	user := User{}
	// 查找email = email 的第一条记录
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&user).Error

	return user, err
}

func (dao *GormUserDAO) FindByID(ctx context.Context, id int64) (User, error) {
	user := User{}
	// 查找email = email 的第一条记录
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return user, ErrUserNotFound
	}
	return user, err
}

func (dao *GormUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	return user, err
}

func (dao *GormUserDAO) FindByWechat(ctx context.Context, openID string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("wechatOpenID = ?", openID).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return user, ErrUserNotFound
	}
	return user, err
}
