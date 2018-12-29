package service

import "github.com/wq1019/cloud_disk/model"

type shareService struct {
	model.ShareStore
}

func NewShareService(ss model.ShareStore) model.ShareService {
	return &shareService{ss}
}
