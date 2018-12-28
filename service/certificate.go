package service

import (
	"github.com/wq1019/cloud_disk/model"
)

type certificateService struct {
	model.CertificateStore
}

func NewCertificateService(cs model.CertificateStore) model.CertificateService {
	return &certificateService{cs}
}
