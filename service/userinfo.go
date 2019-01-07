package service

import (
	"context"
	"github.com/wq1019/cloud_disk/model"
)

type userInfoService struct {
	model.UserInfoStore
}

func UpdateUsedStorage(ctx context.Context, uid, usedStorage int64) error {
	return FromContext(ctx).UpdateUsedStorage(uid, usedStorage)
}

func NewUserInfoService(ss model.UserInfoStore) model.UserInfoService {
	return &userInfoService{ss}
}
