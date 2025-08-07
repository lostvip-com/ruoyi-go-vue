package service

import "common/models"

type I18nService struct {
}

func (s I18nService) GetValue(locale, key string) string {
	i18n := new(models.SysI18n)
	i18n, err := i18n.FindOne(locale, key)
	if err == nil {
		return i18n.LocaleName
	}
	return ""
}

var i18nService *I18nService

func GetI18nService() *I18nService {
	if i18nService == nil {
		i18nService = new(I18nService)
	}
	return i18nService
}
