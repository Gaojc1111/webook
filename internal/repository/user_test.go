package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	cache_mocksvc "webook/internal/repository/cache/mock"
	"webook/internal/repository/dao"
	dao_mocksvc "webook/internal/repository/dao/mock"
)

func TestGormUserDAO_FindByID(t *testing.T) {
	testCases := []struct {
		name string

		mock     func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)
		ctx      context.Context
		id       int64
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "查询成功，缓存未命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				ud := dao_mocksvc.NewMockUserDAO(ctrl)
				uc := cache_mocksvc.NewMockUserCache(ctrl)
				userID := int64(1)
				uc.EXPECT().Get(gomock.Any(), userID).
					Return(domain.User{}, cache.ErrKeyNotExist)
				ud.EXPECT().FindByID(gomock.Any(), userID).
					Return(dao.User{
						ID: 1,
						Email: sql.NullString{
							String: "666@qq.com",
							Valid:  true,
						},
						Password:   "QQqq11!!",
						CreateTime: 1715593591685,
						UpdateTime: 1715593591685,
					}, nil)
				uc.EXPECT().Set(gomock.Any(), domain.User{
					ID:        1,
					Email:     "666@qq.com",
					Password:  "QQqq11!!",
					CreatedAt: 1715593591685,
					UpdatedAt: 1715593591685,
				}).Return(nil)
				return ud, uc
			},
			id: 1,
			wantUser: domain.User{
				ID:        1,
				Email:     "666@qq.com",
				Password:  "QQqq11!!",
				CreatedAt: 1715593591685,
				UpdatedAt: 1715593591685,
			},
			wantErr: nil,
		},
		{
			name: "缓存命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				ud := dao_mocksvc.NewMockUserDAO(ctrl)
				uc := cache_mocksvc.NewMockUserCache(ctrl)
				userID := int64(1)
				uc.EXPECT().Get(gomock.Any(), userID).
					Return(domain.User{
						ID:        1,
						Email:     "666@qq.com",
						Password:  "QQqq11!!",
						CreatedAt: 1715593591685,
						UpdatedAt: 1715593591685,
					}, nil)
				return ud, uc
			},
			id: 1,
			wantUser: domain.User{
				ID:        1,
				Email:     "666@qq.com",
				Password:  "QQqq11!!",
				CreatedAt: 1715593591685,
				UpdatedAt: 1715593591685,
			},
			wantErr: nil,
		},
		{
			name: "未查询到用户",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				ud := dao_mocksvc.NewMockUserDAO(ctrl)
				uc := cache_mocksvc.NewMockUserCache(ctrl)
				userID := int64(121)
				uc.EXPECT().Get(gomock.Any(), userID).
					Return(domain.User{}, cache.ErrKeyNotExist)
				ud.EXPECT().FindByID(gomock.Any(), userID).
					Return(dao.User{}, ErrUserNotFound)
				return ud, uc
			},
			id:       121,
			wantUser: domain.User{},
			wantErr:  ErrUserNotFound,
		},
		{
			name: "redis缓存失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				ud := dao_mocksvc.NewMockUserDAO(ctrl)
				uc := cache_mocksvc.NewMockUserCache(ctrl)
				userID := int64(1)
				uc.EXPECT().Get(gomock.Any(), userID).
					Return(domain.User{}, cache.ErrKeyNotExist)
				ud.EXPECT().FindByID(gomock.Any(), userID).
					Return(dao.User{
						ID: 1,
						Email: sql.NullString{
							String: "666@qq.com",
							Valid:  true,
						},
						Password:   "QQqq11!!",
						CreateTime: 1715593591685,
						UpdateTime: 1715593591685,
					}, nil)
				uc.EXPECT().Set(gomock.Any(), domain.User{
					ID:        1,
					Email:     "666@qq.com",
					Password:  "QQqq11!!",
					CreatedAt: 1715593591685,
					UpdatedAt: 1715593591685,
				}).Return(errors.New("redis err"))
				return ud, uc
			},
			id: 1,
			wantUser: domain.User{
				ID:        1,
				Email:     "666@qq.com",
				Password:  "QQqq11!!",
				CreatedAt: 1715593591685,
				UpdatedAt: 1715593591685,
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ud, uc := tc.mock(ctrl)
			ur := NewCachedUserRepository(ud, uc)
			user, err := ur.FindByID(tc.ctx, tc.id)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}
}
