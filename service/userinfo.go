package service

import (
	"context"
	"github.com/wq1019/cloud_disk/model"
)

type userInfoService struct {
	model.UserInfoStore
}

func CreateUserInfo(ctx context.Context, userInfo *model.UserInfo) (err error) {
	return FromContext(ctx).CreateUserInfo(userInfo)
}

func LoadUserInfo(ctx context.Context, uid int64) (*model.UserInfo, error) {
	return FromContext(ctx).LoadUserInfo(uid)
}

func NewUserInfoService(ss model.UserInfoStore) model.UserInfoService {
	return &userInfoService{ss}
}
