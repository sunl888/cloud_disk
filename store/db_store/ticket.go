package db_store

import (
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/model"
)

type dbTicket struct {
	db *gorm.DB
}

func (dt *dbTicket) TicketIsNotExistErr(err error) bool {
	return model.TicketIsNotExistErr(err)
}

func (dt *dbTicket) TicketLoad(id string) (ticket *model.Ticket, err error) {
	if id == "" {
		return nil, model.ErrTicketNotExist
	}
	ticket = &model.Ticket{}
	err = dt.db.Where(model.Ticket{Id: id}).First(ticket).Error
	if gorm.IsRecordNotFoundError(err) {
		err = model.ErrTicketNotExist
	}
	return
}

func (dt *dbTicket) TicketCreate(ticket *model.Ticket) error {
	return dt.db.Create(ticket).Error
}

func (dt *dbTicket) TicketDelete(id string) error {
	return dt.db.Delete(model.Ticket{Id: id}).Error
}

func NewDBTicket(db *gorm.DB) model.TicketStore {
	return &dbTicket{db: db}
}
