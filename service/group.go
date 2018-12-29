package service

import "github.com/wq1019/cloud_disk/model"

type groupService struct {
	model.GroupStore
}

func NewGroupService(gs model.GroupStore) model.GroupService {
	return &groupService{gs}
}
