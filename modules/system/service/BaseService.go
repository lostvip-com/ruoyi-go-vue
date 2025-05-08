package service

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
