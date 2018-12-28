package model

import (
	"context"
	"errors"
	"time"
)

// 登录凭证
type Ticket struct {
	Id        string `gorm:"type:CHAR(32);PRIMARY_KEY;NOT NULL"`
	UserId    int64  `gorm:"type:BIGINT;index"`
	ExpiredAt time.Time
	CreatedAt time.Time
}

type TicketStore interface {
	TicketLoad(id string) (*Ticket, error)
	TicketCreate(ticket *Ticket) error
	TicketDelete(id string) error
	TicketIsNotExistErr(err error) bool
}

type TicketService interface {
	TicketIsValid(ticketId string) (isValid bool, userId int64, err error)
	// 生成 ticket
	TicketGen(userId int64) (*Ticket, error)
	TicketTTL(ctx context.Context) time.Duration
	TicketDestroy(ticketId string) error
}

var (
	ErrTicketNotExist = errors.New("ticket not exist")
	ErrTicketExisted  = errors.New("ticket 已经存在")
)

func TicketIsNotExistErr(err error) bool {
	return err == ErrTicketNotExist
}
