package service

import "github.com/wq1019/cloud_disk/model"

type userInfoService struct {
	model.UserInfoStore
}

func NewUserInfoService(ss model.ShareStore) model.UserInfoService {
	return &userInfoService{ss}
}
