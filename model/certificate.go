package model

import (
	"errors"
)

type CertificateType uint8

const (
	CertificateUserName CertificateType = iota
	CertificatePhoneNum
	CertificateEmail
)

type Certificate struct {
	Id      int64           `gorm:"type:BIGINT AUTO_INCREMENT;PRIMARY_KEY;NOT NULL"`
	UserId  int64           `gorm:"type:BIGINT;INDEX"`
	Account string          `gorm:"NOT NULL;UNIQUE"`
	Type    CertificateType `gorm:"type:TINYINT"`
}

type CertificateStore interface {
	CertificateExist(account string) (bool, error)
	CertificateLoadByAccount(account string) (*Certificate, error)
	CertificateIsNotExistErr(error) bool
	CertificateCreate(certificate *Certificate) error
	CertificateUpdate(oldAccount, newAccount string, certificateType CertificateType) error
}

var ErrCertificateNotExist = errors.New("certificate not exist")

func CertificateIsNotExistErr(err error) bool {
	return err == ErrCertificateNotExist
}

type CertificateService interface {
	CertificateStore
}
