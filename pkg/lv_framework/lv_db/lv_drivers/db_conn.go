package lv_drivers

import "gorm.io/gorm"

type DbConn struct {
	driver Driver
}

func NewDbConn(cfg *DbConfig) *DbConn {
	conn := &DbConn{}
	switch cfg.DriverType {
	case "mysql":
		conn.driver = &MySQLDriver{config: cfg}
	case "sqlite":
		conn.driver = &SQLiteDriver{config: cfg}
	}
	return conn
}

// 代理所有方法到具体驱动
func (c *DbConn) Setup() (*gorm.DB, error) {
	return c.driver.Setup()
}

func (c *DbConn) GetGormDB() *gorm.DB {
	return c.driver.GetGormDB()
}

func (c *DbConn) TestConnInstance() {
	c.driver.TestConnInstance()
}
