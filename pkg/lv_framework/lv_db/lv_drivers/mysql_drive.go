package lv_drivers

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type MySQLDriver struct {
	config *DbConfig
	gormDB *gorm.DB
}

func (d *MySQLDriver) GetGormDB() *gorm.DB {
	return d.gormDB
}

func (d *MySQLDriver) TestConnInstance() {
	//TODO implement me
	panic("implement me")
}

// 实现Driver接口
func (d *MySQLDriver) Setup() (*gorm.DB, error) {
	url := "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=5000ms"
	url = fmt.Sprintf(url, d.config.User, d.config.Password, d.config.Host, d.config.Port, d.config.DbName)

	// 添加选项处理
	if len(d.config.Options) > 0 {
		url += "&" + buildOptions(d.config.Options)
	}

	db, err := gorm.Open(mysql.Open(url), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("连接数据库失败" + err.Error())
	}
	sqlDB, err := d.gormDB.DB() //dr
	if err != nil {
		panic("连接数据库失败")
	}
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(50)
	//e.GormSqlite.LogMode(true) // ====================打印sql
	return db, err
}

//
//func (e *DbConn) CloseConn() {
//	sqlDB, err := e.gormDB.DB() //dr
//	if err == nil {
//		sqlDB.Close()
//	}
//	myDbConn = nil
//	Logger.Info(" close success !")
//	return
//}
