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
	if user.IsBan == true {
		return nil, errors.UserIsBanned()
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
		Nickname: account,
		Profile:  "这货很懒,什么都没有说哦",
		GroupId:  1,
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

func (uSvc *userService) UserUpdateUsedStorage(userId int64, newUsedStorage uint64) error {
	return uSvc.UserStore.UserUpdate(userId, map[string]interface{}{
		"used_storage": newUsedStorage,
	})
}

func (uSvc *userService) UserUpdateBanStatus(userId int64, newBanStatus bool) error {
	return uSvc.UserStore.UserUpdate(userId, map[string]interface{}{
		"is_ban": newBanStatus,
	})
}

func NewUserService(us model.UserStore, cs model.CertificateStore, tSvc model.TicketService, h hasher.Hasher) model.UserService {
	return &userService{us, cs, tSvc, h}
}

func UserLoad(ctx context.Context, id int64) (*model.User, error) {
	return FromContext(ctx).UserLoad(id)
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

func UserUpdateUsedStorage(ctx context.Context, userId int64, newUsedStorage uint64) error {
	return FromContext(ctx).UserUpdateUsedStorage(userId, newUsedStorage)
}

func UserUpdate(ctx context.Context, userId int64, data map[string]interface{}) error {
	return FromContext(ctx).UserUpdate(userId, data)
}

func UserUpdateBanStatus(ctx context.Context, userId int64, newBanStatus bool) error {
	return FromContext(ctx).UserUpdateBanStatus(userId, newBanStatus)
}
