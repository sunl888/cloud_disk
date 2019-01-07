package service

import (
	"context"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/model"
	"github.com/wq1019/cloud_disk/pkg/hasher"
)

type userService struct {
	model.UserStore
	cs   model.CertificateStore
	tSvc model.TicketService
	h    hasher.Hasher
}

func (uSvc *userService) UserLogin(account, password string) (ticket *model.Ticket, err error) {
	c, err := uSvc.cs.CertificateLoadByAccount(account)
	if err != nil {
		if uSvc.cs.CertificateIsNotExistErr(err) { //账号不存在
			err = errors.ErrAccountNotFound()
		}
		return nil, err
	}
	user, err := uSvc.UserStore.UserLoad(c.UserId)
	if err != nil {
		return nil, err
	}
	if uSvc.h.Check(password, user.Password) {
		// 登录成功
		return uSvc.tSvc.TicketGen(user.Id)
	}

	return nil, errors.ErrPassword()
}

func (uSvc *userService) UserRegister(account string, certificateType model.CertificateType, password string) (userId int64, err error) {
	if exist, err := uSvc.cs.CertificateExist(account); err != nil {
		return 0, err
	} else if exist {
		return 0, errors.ErrAccountAlreadyExisted()
	}
	user := &model.User{
		Name:     account,
		Password: uSvc.h.Make(password),
		PwPlain:  password,
		UserInfo: &model.UserInfo{
			Nickname: account,
			Profile:  "这货很懒,什么都没有说哦",
			IsBan:    false,
			GroupId:  1,
		},
	}
	if err := uSvc.UserStore.UserCreate(user); err != nil {
		return 0, err
	}
	certificate := &model.Certificate{UserId: user.Id, Account: account, Type: certificateType}
	if err := uSvc.cs.CertificateCreate(certificate); err != nil {
		return 0, err
	}
	return user.Id, nil
}

func (uSvc *userService) UserUpdatePassword(userId int64, newPassword string) error {
	return uSvc.UserStore.UserUpdate(userId, map[string]interface{}{
		"password": uSvc.h.Make(newPassword),
		"pw_plain": newPassword,
	})
}

func NewUserService(us model.UserStore, cs model.CertificateStore, tSvc model.TicketService, h hasher.Hasher) model.UserService {
	return &userService{us, cs, tSvc, h}
}

func UserLoad(ctx context.Context, id int64) (*model.User, error) {
	return FromContext(ctx).UserLoad(id)
}

func UserLoadAndRelated(ctx context.Context, id int64) (*model.User, error) {
	return FromContext(ctx).UserLoadAndRelated(id)
}

func UserCreate(ctx context.Context, user *model.User) error {
	return FromContext(ctx).UserCreate(user)
}

func UserLogin(ctx context.Context, account, password string) (*model.Ticket, error) {
	return FromContext(ctx).UserLogin(account, password)
}

func UserRegister(ctx context.Context, account string, certificateType model.CertificateType, password string) (userId int64, err error) {
	return FromContext(ctx).UserRegister(account, certificateType, password)
}

func UserUpdatePassword(ctx context.Context, userId int64, newPassword string) error {
	return FromContext(ctx).UserUpdatePassword(userId, newPassword)
}

func UserListByUserIds(ctx context.Context, userIds []interface{}) ([]*model.User, error) {
	return FromContext(ctx).UserListByUserIds(userIds)
}
