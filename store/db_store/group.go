package db_store

import (
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/model"
)

type dbGroup struct {
	db *gorm.DB
}

func (g *dbGroup) GroupCreate(group *model.Group) (err error) {
	err = g.db.Create(&group).Error
	return
}

func (g *dbGroup) GroupExist(name string) (isExist bool, err error) {
	var count int8
	err = g.db.Model(model.Group{}).Where("name = ?", name).Limit(1).Count(&count).Error
	isExist = count > 0
	return
}

func (g *dbGroup) GroupDelete(id int64) (err error) {
	group := model.Group{}
	err = g.db.Where("id = ?", id).First(&group).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = errors.RecordNotFound("用户组不存在")
		}
		return err
	}
	userCount := g.db.Model(&group).Association("Users").Count()
	if userCount > 0 {
		err = errors.GroupNotAllowBeDelete("该组不允许删除, 因为组里面有用户")
		return
	}
	err = g.db.Delete(&group).Error
	return err
}

func (g *dbGroup) GroupUpdate(id int64, data map[string]interface{}) (err error) {
	if id <= 0 {
		return model.ErrGroupNotExist
	}
	return g.db.Model(model.Group{Id: id}).Select("name", "max_storage", "allow_share").Updates(data).Error
}

func (g *dbGroup) GroupList(offset, limit int64) (groups []*model.WrapGroupList, count int64, err error) {
	groups = make([]*model.WrapGroupList, 0, 10)
	err = g.db.Table("`groups` g").
		Select("g.*,count(u.id) as user_count").
		Joins("LEFT JOIN `users` u ON u.group_id = g.id").
		Group("g.id").
		Offset(offset).
		Limit(limit).
		Scan(&groups).
		Count(&count).
		Error
	return
}

func NewDBGroup(db *gorm.DB) model.GroupStore {
	return &dbGroup{db}
}
