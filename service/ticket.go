package service

import (
	"context"
	"encoding/hex"
	"github.com/satori/go.uuid"
	"github.com/wq1019/cloud_disk/model"
	"time"
)

type ticketService struct {
	ts        model.TicketStore
	ticketTTL time.Duration
}

func (tSvc *ticketService) TicketTTL(ctx context.Context) time.Duration {
	return tSvc.ticketTTL
}

func (tSvc *ticketService) TicketIsValid(ticketId string) (isValid bool, userId int64, err error) {
	ticket, err := tSvc.ts.TicketLoad(ticketId)
	if err != nil {
		if tSvc.ts.TicketIsNotExistErr(err) {
			return false, 0, nil
		} else {
			return false, 0, err
		}
	}
	return time.Now().UTC().Before(ticket.ExpiredAt), ticket.UserId, nil
}

func (tSvc *ticketService) TicketGen(userId int64) (*model.Ticket, error) {
	u4 := uuid.NewV4()
	now := time.Now().UTC()
	ticket := &model.Ticket{
		Id:        hex.EncodeToString(u4.Bytes()),
		UserId:    userId,
		ExpiredAt: now.Add(tSvc.ticketTTL),
		CreatedAt: now,
	}
	err := tSvc.ts.TicketCreate(ticket)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}

func (tSvc *ticketService) TicketDestroy(ticketId string) error {
	return tSvc.ts.TicketDelete(ticketId)
}

func NewTicketService(ts model.TicketStore, ticketTTL time.Duration) model.TicketService {
	return &ticketService{ts: ts, ticketTTL: ticketTTL}
}

func TicketDestroy(ctx context.Context, ticketId string) error {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.TicketDestroy(ticketId)
	}
	return ServiceError
}

func TicketIsValid(ctx context.Context, ticketId string) (isValid bool, userId int64, err error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.TicketIsValid(ticketId)
	}
	return false, 0, ServiceError
}
