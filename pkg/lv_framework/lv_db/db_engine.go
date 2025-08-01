package lv_db

import (
	"fmt"
	_ "github.com/glebarez/go-sqlite"
	"github.com/glebarez/sqlite"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lostvip-com/lv_framework/lv_global"
	"github.com/lostvip-com/lv_framework/lv_log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"strings"
	"sync"
)

type dbEngine struct {
	gormMap map[string]*gorm.DB
}

var (
	instance *dbEngine
	once     sync.Once
)

func init() {
	fmt.Println("-----init orm-------")
}

// GetInstance 初始化数据操作 driver为数据库类型
func GetInstance() *dbEngine {
	once.Do(func() {
		instance = new(dbEngine)
		instance.gormMap = make(map[string]*gorm.DB)
	})
	return instance
}

// GetDB 获取操作实例 如果传入slave 并且成功配置了slave 返回slave orm引擎 否则返回master orm引擎
func (db *dbEngine) GetDB(dbName string) *gorm.DB {
	gdb := db.gormMap[dbName]
	if gdb == nil {
		var config = lv_global.Config()
		driverName := config.GetDriver(dbName)
		url := config.GetValueStr(fmt.Sprintf("application.datasource.%s.url", dbName))
		gdb = createGormDB(driverName, url)
		db.gormMap[dbName] = gdb
	}
	return gdb
}
func GetOrmDefault() *gorm.DB {
	var config = lv_global.Config()
	dbName := config.GetDBNameDefault()
	if dbName == "" { //先取配置文件
		dbName = lv_global.Config().GetValueStr("application.datasource.default")
		config.SetDBNameDefault(dbName)
		config.SetDBDriverDefault(config.GetDriver(dbName))
	}
	if dbName == "" {
		panic("default database not found!")
	}
	defaultDb := GetInstance().gormMap[dbName]
	if defaultDb == nil {
		defaultDb = GetInstance().GetDB(dbName)
	}
	return defaultDb
}

func createGormDB(driverName, url string) *gorm.DB {
	if !strings.Contains(url, "?") {
		url = url + "?"
	}
	params := "parseTime=true"
	if !strings.Contains(url, params) { //自动解析时间类型到time.Time!!
		if strings.HasSuffix(url, "?") {
			url = url + params
		} else {
			url = url + "&" + params
		}
	}
	params = "charset=utf8mb4"
	if !strings.Contains(url, params) {
		if !strings.Contains(url, params) {
			if strings.HasSuffix(url, "?") {
				url = url + params
			} else {
				url = url + "&" + params
			}
		}
	}
	var dialector gorm.Dialector
	if "mysql" == driverName {
		dialector = mysql.Open(url)
	} else if "sqlite" == driverName {
		dialector = sqlite.Open(url)
	} else {
		panic("不支持的数据库类型：" + driverName)
	}
	showSql := lv_global.Config().GetBool("application.datasource.show-sql")
	config := &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}} //表名用单数
	if showSql {
		config.Logger = logger.Default.LogMode(logger.Info)
	}

	gormDB, err := gorm.Open(dialector, config)
	if err != nil {
		panic("  gorm init fail！" + err.Error())
	}
	if lv_global.IsDebug {
		gormDB = gormDB.Debug() //会开启sql打印
	}
	sqlDB, err := gormDB.DB() //dr
	if err != nil {
		panic("连接数据库失败")
	}
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(50)
	lv_log.Info("###########  gorm init success！ #################")
	return gormDB
}
