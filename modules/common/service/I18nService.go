package service

import (
	"common/models"
	"github.com/gin-gonic/gin"
)

type I18nService struct {
}

func (s *I18nService) GetLang(c *gin.Context) string {
	lang := c.GetHeader("Accept-Language")
	if lang == "" {
		return "zh"
	}
	return lang[0:2]
}
func (s *I18nService) GetLocale(c *gin.Context) string {
	lang := c.GetHeader("Accept-Language")
	if lang == "" {
		return "zh-CN"
	}
	return lang
}
func (s *I18nService) GetI18nText(locale, key string) string {
	i18n := new(models.SysI18n)
	i18n, err := i18n.FindOne(locale, key)
	if err != nil || i18n.LocaleName == "" {
		return key
	}
	return i18n.LocaleName
}

var i18nService *I18nService

func GetI18nServiceInstance() *I18nService {
	if i18nService == nil {
		i18nService = new(I18nService)
	}
	return i18nService
}
