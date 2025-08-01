package service

import (
	"github.com/lostvip-com/lv_framework/lv_db"
	"gorm.io/gorm"
)

type BaseService struct{}

var baseService *BaseService

func GetBaseService() *BaseService {
	if baseService == nil {
		baseService = &BaseService{}
	}
	return baseService
}

func (svc *BaseService) IsAdmin(userId int64) bool {
	if userId == 1 {
		return true
	} else {
		return false
	}
}

func (svc *BaseService) GetDb() *gorm.DB {
	return lv_db.GetOrmDefault()
}
