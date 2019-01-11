package service

import (
	"context"
	"github.com/wq1019/cloud_disk/model"
)

type groupService struct {
	model.GroupStore
}

func (g *groupService) GroupCreate(group *model.Group) (err error) {
	isExist, err := g.GroupStore.GroupExist(group.Name)
	if err != nil {
		return err
	}
	if isExist {
		return model.ErrGroupAlreadyExist
	}
	return g.GroupStore.GroupCreate(group)
}

func GroupCreate(ctx context.Context, group *model.Group) (err error) {
	return FromContext(ctx).GroupCreate(group)
}

func GroupDelete(ctx context.Context, id int64) (err error) {
	return FromContext(ctx).GroupDelete(id)
}

func GroupExist(ctx context.Context, name string) (isExist bool, err error) {
	return FromContext(ctx).GroupExist(name)
}

func GroupUpdate(ctx context.Context, id int64, data map[string]interface{}) (err error) {
	return FromContext(ctx).GroupUpdate(id, data)
}

func GroupList(ctx context.Context, offset, limit int64) (groups []*model.WrapGroupList, count int64, err error) {
	return FromContext(ctx).GroupList(offset, limit)
}

func NewGroupService(gs model.GroupStore) model.GroupService {
	return &groupService{gs}
}
