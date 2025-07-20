package lv_drivers

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

// Driver 数据库驱动接口
type Driver interface {
	Setup() (*gorm.DB, error)
	GetGormDB() *gorm.DB
	TestConnInstance()
}

// DbConfig 通用数据库配置
type DbConfig struct {
	DriverType string // 驱动类型 mysql/sqlite
	User       string
	Password   string
	Host       string
	Port       string
	DbName     string
	Options    map[string]string
}

// 新增辅助函数处理连接选项
func buildOptions(options map[string]string) string {
	var opts []string
	for k, v := range options {
		opts = append(opts, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(opts, "&")
}
