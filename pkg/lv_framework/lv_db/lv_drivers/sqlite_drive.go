package lv_drivers

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"github.com/lostvip-com/lv_framework/lv_global"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"path/filepath"
)

type SQLiteDriver struct {
	config *DbConfig
	gormDB *gorm.DB
}

func (d *SQLiteDriver) GetGormDB() *gorm.DB {
	return d.gormDB
}

func (d *SQLiteDriver) TestConnInstance() {
	//TODO implement me
	panic("implement me")
}

// params := "parseTime=true"
// if !strings.Contains(url, params) { //自动解析时间类型到time.Time!!
// if strings.HasSuffix(url, "?") {
// url = url + params
// } else {
// url = url + "&" + params
// }
// }
// params = "charset=utf8mb4"
func (d *SQLiteDriver) Setup() (*gorm.DB, error) {
	dsn := fmt.Sprintf("file:%s?%s", d.config.DbName, buildOptions(d.config.Options))
	fn := filepath.Join(dsn, "db")
	dialector := sqlite.Open(fn)
	gormcfg := &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}}
	if lv_global.IsDebug { // 开启SQL日志记录，并设置为Info级别
		gormcfg.Logger = logger.Default.LogMode(logger.Info)
	} else {
		gormcfg.Logger = logger.Default.LogMode(logger.Error)
	}
	gormDB, err := gorm.Open(dialector, gormcfg)
	if err != nil {
		panic("连接数据库失败" + err.Error())
	}
	sqlDB, err := gormDB.DB() //dr
	if err != nil {
		panic("连接数据库失败")
	}
	sqlDB.SetMaxIdleConns(0)
	sqlDB.SetMaxOpenConns(100)
	// ... 相同连接池配置 ...
	return gormDB, err
}
