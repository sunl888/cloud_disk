package service

import (
	"context"
	"github.com/gin-gonic/gin"
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
	user := &model.User{Password: uSvc.h.Make(password), PwPlain: password}
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
	if service, ok := ctx.Value("service").(Service); ok {
		return service.UserLoad(id)
	}
	return nil, ServiceError
}

func UserCreate(ctx context.Context, user *model.User) error {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.UserCreate(user)
	}
	return ServiceError
}

func UserLogin(ctx *gin.Context, account, password string) (*model.Ticket, error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.UserLogin(account, password)
	}
	return nil, ServiceError
}

func UserRegister(ctx *gin.Context, account string, certificateType model.CertificateType, password string) (userId int64, err error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.UserRegister(account, certificateType, password)
	}
	return 0, ServiceError
}

func UserUpdatePassword(ctx context.Context, userId int64, newPassword string) error {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.UserUpdatePassword(userId, newPassword)
	}
	return ServiceError
}

func UserListByUserIds(ctx context.Context, userIds []interface{}) ([]*model.User, error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.UserListByUserIds(userIds)
	}
	return nil, ServiceError
}
