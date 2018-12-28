package db_store

import (
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/model"
)

type dbCertificate struct {
	db *gorm.DB
}

func (c *dbCertificate) CertificateExist(account string) (bool, error) {
	var count uint8
	err := c.db.Model(&model.Certificate{}).Where(model.Certificate{Account: account}).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (c *dbCertificate) CertificateIsNotExistErr(err error) bool {
	return model.CertificateIsNotExistErr(err)
}

func (c *dbCertificate) CertificateLoadByAccount(account string) (certificate *model.Certificate, err error) {
	if account == "" {
		return nil, model.ErrCertificateNotExist
	}
	certificate = &model.Certificate{}
	err = c.db.Where(model.Certificate{Account: account}).First(&certificate).Error
	if gorm.IsRecordNotFoundError(err) {
		err = model.ErrCertificateNotExist
	}
	return
}

func (c *dbCertificate) CertificateCreate(certificate *model.Certificate) error {
	return c.db.Create(certificate).Error
}

func (c *dbCertificate) CertificateUpdate(oldAccount, newAccount string, certificateType model.CertificateType) error {
	return c.db.Model(&model.User{}).
		Where("account", oldAccount).
		Where("type", certificateType).
		UpdateColumn("account", newAccount).Error
}

func NewDBCertificate(db *gorm.DB) model.CertificateStore {
	return &dbCertificate{db: db}
}
